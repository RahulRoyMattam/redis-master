package env

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("cmd", "redis-master-env")
)

//SetEnvironmentVariables for the app based on the environment passed in as a command line argument.
//the environment for an app is set manually using the "environment" config var.
//Make sure you have the correct environment set before deploying.
func SetEnvironmentVariables() {
	switch env := os.Getenv("REDIS_MASTER_ENV"); env {
	case "LOCAL":
		setLocalVariables()
	default:
		log.WithField("ENVIRONMENT", env).Fatal("$REDIS_MASTER_ENV is not set to a valid property")
	}
}
