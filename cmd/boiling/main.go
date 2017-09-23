package main

import (
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"

	"github.com/mutaborius/boiling-api/api"
	"github.com/mutaborius/boiling-api/db"
)

func main() {
	d, err := db.New("boiling", "boiling", "boiling")
	if err != nil {
		log.Fatal(err)
	}

	a, err := api.New(d)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run(iris.Addr(":8080"))
	if err != nil {
		log.Fatal(err)
	}
}
