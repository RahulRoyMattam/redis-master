package redis

import (
	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger

	//Instances is a Slice of all available redis instances
	Instances []*Instance
)

func init() {
	//initialize logging to sentry
	log = logrus.New()
	hook, err := logrus_sentry.NewWithClientSentryHook(raven.DefaultClient, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	})

	if err == nil {
		hook.StacktraceConfiguration.Enable = true
		log.Hooks.Add(hook)
	}
}

//Connect initializes the connection pool to all available redis instances
func Connect() {
	Instances = getRedisInstances()
	go loadbalance()
	initializeInstances()
}

//Dispose off all connections created
func Dispose() {
	for _, inst := range Instances {
		inst.dispose()
	}
}

func initializeInstances() {
	for _, inst := range Instances {
		inst.initializePool()
		if inst.PoolActive {
			go inst.ping()
		}
	}
}
