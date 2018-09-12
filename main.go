package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rahulroymattam/redis-master/env"
	"github.com/rahulroymattam/redis-master/redis"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("cmd", "redis-master-main")
)

var appEnv = flag.String("env", "", "set the environment variable for debugging in local. Options include LOCAL, SLAVE, MASTER")

func main() {
	//initialize environment
	flag.Parse()
	if *appEnv != "" {
		os.Setenv("REDIS_MASTER_ENV", *appEnv)
	}
	env.SetEnvironmentVariables()

	initialize()
	defer destroy()

	port, isPortAvailable := os.LookupEnv("PORT")
	if isPortAvailable == false {
		port = "8000"
	}
	fmt.Println("Using port:", port)
	port = ":" + port
	router := NewRouter()
	http.Handle("/", router)
	log.Println(http.ListenAndServe(port, nil))
}

func initialize() {
	redis.Connect()
}

func destroy() {
	redis.Dispose()
}
