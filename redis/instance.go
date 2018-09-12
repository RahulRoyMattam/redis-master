package redis

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	pingPeriod  time.Duration = 10 * time.Second
	waitSleep                 = time.Second * 10
	waitTimeout               = time.Minute * 10
)

//Instance references to the available redis instances
type Instance struct {
	//Key for the redis instance
	Key string `json:"configKey"`

	//URL of the redis instance
	URL string `json:"url"`

	//PoolActive flag for monitoring pool health
	PoolActive bool `json:"poolActive"`
	//ConnActive flag for monitoring connection establishment health
	ConnActive bool `json:"connActive"`
	//AvailableSpace is a integer representative of available space in the redis cluster. Usually in terms of number of bytes free
	AvailableSpace int64 `json:"availableSpace"`

	pool *redis.Pool
}

func (i *Instance) initializePool() {
	i.PoolActive = true
	var err error

	maxActive, err := strconv.Atoi(os.Getenv("APP_MAX_REDIS_CONN"))
	if err != nil || maxActive == 0 {
		log.WithField("APP_MAX_REDIS_CONN", maxActive).Fatal("$APP_MAX_REDIS_CONN must be set to a non zero value")
	}

	i.pool, err = NewRedisPoolFromURL(i.URL, maxActive)
	if err != nil {
		log.WithFields(logrus.Fields{"key": i.Key, "url": i.URL, "err": err}).Error("Unable to create Redis pool")
		i.PoolActive = false
	}
	waited, err := waitForAvailability(i.URL, waitTimeout, "Waiting for redis to be available", wait)
	if !waited || err != nil {
		log.WithFields(logrus.Fields{"key": i.Key, "url": i.URL, "waitTimeout": waitTimeout, "err": err}).Error("Redis not available by timeout!")
		i.PoolActive = false
	}
}

//Dispose off all connections created
func (i *Instance) dispose() {
	i.pool.Close()
}

//doSafe function abstracts establishing the redis connection, calling the Do function and closing the connection safely.
func (i *Instance) doSafe(ctx context.Context, commandName string, args ...interface{}) (interface{}, error) {
	conn := i.initConnSafe()
	defer conn.Close()
	done := make(chan *struct {
		ret interface{}
		err error
	})
	go func() {
		ret, err := conn.Do(commandName, args...)
		done <- &struct {
			ret interface{}
			err error
		}{ret: ret, err: err}
	}()

	for {
		select {
		case res := <-done:
			return res.ret, res.err
		case <-ctx.Done():
			return nil, errors.New("redis-master : doSafe was cancelled")
		}
	}
}

//connSafe : get the redis connection safely
func (i *Instance) initConnSafe() (conn redis.Conn) {
	conn = i.pool.Get()
	if err := conn.Err(); err != nil {
		log.WithFields(logrus.Fields{"key": i.Key, "url": i.URL, "err": err}).Error("Initializing Redis Connection Failed.")
		i.ConnActive = false
		return
	}
	i.ConnActive = true
	return conn
}

func (i *Instance) ping() {
	pingTimer := time.NewTicker(pingPeriod)
	defer func() {
		pingTimer.Stop()
	}()

	for {
		select {
		case <-pingTimer.C:
			ctx := context.Background()
			info, err := redis.String(i.doSafe(ctx, "INFO", "memory"))
			if err != nil {
				log.Error("PING failed: " + err.Error())
				i.ConnActive = false
			}
			mem := parseInfo(info)
			max, merr := strconv.ParseInt(mem["maxmemory"], 10, 64)
			used, uerr := strconv.ParseInt(mem["used_memory"], 10, 64)

			if mem["maxmemory_policy"] != "allkeys-lru" {
				evictionErr := errors.New("redis-master : Redis Key eviction policy is not set to allkeys-lru")
				log.Error(evictionErr.Error())
			}

			if uerr != nil || merr != nil {
				if uerr != nil {
					err := errors.Wrap(uerr, "redis-master: Failed to convert to Int64")
					log.Error(err.Error())
				}
				if merr != nil {
					err := errors.Wrap(merr, "redis-master: Failed to convert to Int64")
					log.Error(err.Error())
				}
				i.ConnActive = false
			}
			i.AvailableSpace = max - used
			updateBalancer <- i
		}
	}
}

//wait for redis availability
func wait(_ time.Time, waitingMessage string) error {
	log.Info(waitingMessage)
	time.Sleep(waitSleep)
	return nil
}
