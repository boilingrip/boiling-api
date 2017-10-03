package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kataras/iris"

	"github.com/boilingrip/boiling-api/db"
)

func TestMain(m *testing.M) {
	d, err := cleanDB()
	if err != nil {
		panic(err)
	}
	_, err = getDefaultAPIWithDB(d)
	if err != nil {
		panic(err)
	}
	ret := m.Run()
	err = stopDefaultAPI()
	if err != nil {
		panic(err)
	}
	os.Exit(ret)
}

func cleanDB() (db.BoilingDB, error) {
	d, err := db.New("boilingtest", "boilingtest", "boilingtest")
	if err != nil {
		return nil, err
	}

	inner, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s", "boilingtest", "boilingtest", "boilingtest"))
	if err != nil {
		return nil, err
	}
	defer inner.Close()

	file, err := ioutil.ReadFile("../db/create.sql")
	if err != nil {
		return nil, err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_, err := inner.Exec(request)
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

var defaultAPI *struct {
	api *API
	wg  *sync.WaitGroup
}

func getDefaultAPIWithDB(d db.BoilingDB) (*API, error) {
	if defaultAPI != nil {
		defaultAPI.api.db = d
		return defaultAPI.api, nil
	}
	toReturn := &struct {
		api *API
		wg  *sync.WaitGroup
	}{}
	a, err := New(d)
	if err != nil {
		return nil, err
	}
	toReturn.wg = &sync.WaitGroup{}
	toReturn.wg.Add(1)
	go func() {
		defer toReturn.wg.Done()
		err := a.app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("Waiting for API to start...")
	time.Sleep(3 * time.Second)

	toReturn.api = a
	defaultAPI = toReturn
	return defaultAPI.api, nil
}

func stopDefaultAPI() error {
	err := defaultAPI.api.Stop()
	if err != nil {
		return err
	}
	defaultAPI.wg.Wait()
	return nil
}

type dbWithLogin struct {
	db       db.BoilingDB
	token    string
	user     db.User
	password string
}

func cleanDBWithLogin() (*dbWithLogin, error) {
	d, err := cleanDB()
	if err != nil {
		return nil, err
	}

	err = d.SignUpUser("sometestuser", "sometestpw12345", "some@ex.am.ple.com")
	if err != nil {
		return nil, err
	}

	u, err := d.LoginAndGetUser("sometestuser", "sometestpw12345")
	if err != nil {
		return nil, err
	}

	tok, err := d.InsertTokenForUser(*u)
	if err != nil {
		return nil, err
	}

	err = d.UpdateUserAddPrivileges(u.ID, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	toReturn := &dbWithLogin{
		db:       d,
		token:    tok.Token,
		user:     *u,
		password: "sometestpw12345",
	}

	return toReturn, nil
}
