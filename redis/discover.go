package redis

import (
	"os"
	"strings"
)

func getRedisInstances() (items []*Instance) {
	data := os.Environ()
	for _, item := range data {
		key, val := getEnvKeyVal(item)
		if strings.HasPrefix(key, "REDIS") && strings.HasSuffix(key, "URL") {
			items = append(items, &Instance{
				Key:            key,
				URL:            val,
				PoolActive:     false,
				ConnActive:     false,
				AvailableSpace: 0,
			})
		}
	}
	return items
}

func getEnvKeyVal(item string) (key string, val string) {
	splits := strings.SplitN(item, "=", 2)
	key = splits[0]
	val = splits[1]
	return
}
