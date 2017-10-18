package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kataras/iris"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/boilingrip/boiling-api/api"
	"github.com/boilingrip/boiling-api/db"
)

type ConfigFile struct {
	Boiling Config `yaml:"boiling"`
}

type Config struct {
	Database         string `yaml:"database"`
	DatabaseUser     string `yaml:"database_user"`
	DatabasePassword string `yaml:"database_password"`

	ListenAddress string `yaml:"listen_addr"`

	CreateSQL string `yaml:"create_sql"`
	ResetHour int    `yaml:"reset_hour"`
}

func (c Config) validate() error {
	if len(c.Database) == 0 {
		return errors.New("database must be set")
	}
	if len(c.DatabaseUser) == 0 {
		return errors.New("database user must be set")
	}
	if len(c.DatabasePassword) == 0 {
		return errors.New("database password must be set")
	}
	if len(c.ListenAddress) == 0 {
		return errors.New("listen address must be set")
	}
	if len(c.CreateSQL) == 0 {
		return errors.New("create SQL must be set")
	}

	return nil
}

func parseConfig(path string) (*ConfigFile, error) {
	if path == "" {
		return nil, errors.New("no configPath path specified")
	}

	f, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfgFile ConfigFile
	err = yaml.Unmarshal(contents, &cfgFile)
	if err != nil {
		return nil, err
	}

	return &cfgFile, nil
}

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "boiling.yaml", "Path to the configuration file")
}

func main() {
	flag.Parse()
	cfg, err := parseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = cfg.Boiling.validate()
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("rebuilding DB...")
	err = cleanDB(cfg.Boiling)
	if err != nil {
		log.Fatalln("unable to reset DB: ", err)
	}

	d, err := db.New(cfg.Boiling.Database, cfg.Boiling.DatabaseUser, cfg.Boiling.DatabasePassword)
	if err != nil {
		log.Fatal(err)
	}

	a, err := api.New(d)
	if err != nil {
		log.Fatal(err)
	}

	closing := make(chan struct{})
	wg := sync.WaitGroup{}
	shutdown := make(chan struct{})

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-quit:
				log.Infoln("received SIGINT/SIGTERM, shutting down...")
				err := a.Stop()
				if err != nil {
					log.Warnln("unable to shut down cleanly: ", err)
				}

				err = d.Close()
				if err != nil {
					log.Warnln("unable to close DB cleanly: ", err)
				}

				close(closing)
			case <-closing:
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-t.C:
				if _, err := os.Stat("update_available"); !os.IsNotExist(err) {
					log.Infoln("update_available exists, attempting upgrade")

					err = os.Remove("update_available")
					if err != nil {
						log.Errorln("unable to delete update_available: ", err)
					}

					err = a.Stop()
					if err != nil {
						log.Warnln("unable to shut down cleanly: ", err)
					}

					err = d.Close()
					if err != nil {
						log.Warnln("unable to close DB cleanly: ", err)
					}

					os.Exit(0)
				}
			case <-closing:
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			now := time.Now()
			now = now.Add(24 * time.Hour)
			midnight := time.Date(now.Year(), now.Month(), now.Day(), cfg.Boiling.ResetHour, 0, 0, 0, time.Local)
			log.Infoln("using midnight: ", midnight)
			log.Infoln("it's now: ", time.Now())
			select {
			case <-time.After(midnight.Sub(time.Now())):
				log.Infoln("it's midnight - shutting down for DB rebuild")

				err = a.Stop()
				if err != nil {
					log.Warnln("unable to shut down cleanly: ", err)
					os.Exit(1)
				}
				<-shutdown

				err = d.Close()
				if err != nil {
					log.Warnln("unable to close DB cleanly: ", err)
					os.Exit(1)
				}

				os.Exit(0)
			case <-closing:
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = a.Run(iris.Addr(cfg.Boiling.ListenAddress))
		close(shutdown)
		if err != nil {
			close(closing)
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func cleanDB(cfg Config) error {
	dbConn, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s", cfg.DatabaseUser, cfg.DatabasePassword, cfg.Database))
	if err != nil {
		return err
	}
	defer dbConn.Close()

	createSQL := os.ExpandEnv(cfg.CreateSQL)
	log.Infoln("using SQL create script at: ", createSQL)
	file, err := ioutil.ReadFile(createSQL)
	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := dbConn.Exec(request)
		if err != nil {
			return err
		}
	}

	return nil
}
