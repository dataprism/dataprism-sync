package main

import (
	"github.com/dataprism/dataprism-commons/api"
	"github.com/sirupsen/logrus"
	"flag"
	"strconv"
	"github.com/dataprism/dataprism-sync/sync"
	"github.com/dataprism/dataprism-commons/core"
	"os"
)

func main() {
	var port = flag.Int("p", 6400, "the port of the dataprism sync rest api")

	flag.Parse()

	platform, err := core.InitializePlatform()
	if err != nil {
		logrus.Error(err.Error())
		os.Exit(1)
	}

	// -- create the api
	API := api.CreateAPI("0.0.0.0:" + strconv.Itoa(*port))

	// -- route the api endpoints
	sync.CreateRoutes(platform, API)

	// -- start serving the api
	err = API.Start()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	logrus.Info("API listening on http://0.0.0.0:" + strconv.Itoa(*port) + "/v1")
}