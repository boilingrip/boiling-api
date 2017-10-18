package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"

	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"os/signal"
	"sync"
	"syscall"

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

	d, err := db.New(cfg.Boiling.Database, cfg.Boiling.DatabaseUser, cfg.Boiling.DatabasePassword)
	if err != nil {
		log.Fatal(err)
	}

	a, err := api.New(d)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-quit

		log.Infoln("received SIGINT/SIGTERM, shutting down...")
		err := a.Stop()
		if err != nil {
			log.Warnln("unable to shut down cleanly: ", err)
		}

		err = d.Close()
		if err != nil {
			log.Warnln("unable to close DB cleanly: ", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = a.Run(iris.Addr(cfg.Boiling.ListenAddress))
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
