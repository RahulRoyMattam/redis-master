package redis

//Set a key value pair to the redis database
import (
	"context"
	"strings"

	"github.com/garyburd/redigo/redis"
)

//Info gets redis server level details corresponding to the section type requested.
func Info(ctx context.Context, section string) (info []*InfoDisplay) {
	for _, inst := range Instances {
		ret, err := redis.String(inst.doSafe(ctx, "INFO", section))
		info = append(info, &InfoDisplay{
			Key:       inst.Key,
			URL:       inst.URL,
			RedisInfo: parseAllInfo(ret),
			Error:     err,
		})
	}
	return info
}

//Keys gets all the keys stored across the redis instances matching the pattern specified.
func Keys(ctx context.Context, pattern string) (keys []string, err error) {
	responses, err := doAll(ctx, "KEYS", pattern)
	if err != nil {
		return nil, err
	}

	for _, res := range responses {
		values, err := redis.Values(res.reply, res.err)
		if err != nil {
			return nil, err
		}

		for len(values) > 0 {
			var element string
			values, err = redis.Scan(values, &element)
			if err != nil {
				return nil, err
			}
			keys = append(keys, element)
		}
	}
	return keys, nil
}

//FlushAll : Delete all the keys of all the existing databases, not just the currently selected one. This command never fails.
func FlushAll(ctx context.Context) (value string, err error) {
	responses, err := doAll(ctx, "FLUSHALL")
	if err != nil {
		return value, err
	}
	for _, res := range responses {
		value, err = redis.String(res.reply, res.err)
		if value != "ok" {
			return value, nil
		}
	}
	return value, nil
}

//InfoDisplay is a collection of redis info of available redis clusters.
type InfoDisplay struct {
	Key       string                       `json:"key"`
	URL       string                       `json:"url"`
	RedisInfo map[string]map[string]string `json:"info"`
	Error     error                        `json:"error"`
}

//parseAllInfo parses redis info string to a map for viewing by human
func parseAllInfo(info string) (redis map[string]map[string]string) {
	var lines = strings.Split(info, "\r\n")
	redis = make(map[string]map[string]string)
	current := ""
	for _, line := range lines {
		if len(line) > 0 {
			if strings.HasPrefix(line, "#") {
				current = line
				redis[current] = make(map[string]string)
			} else {
				splits := strings.SplitN(line, ":", 2)
				redis[current][splits[0]] = splits[1]
			}
		}
	}
	return redis
}

//parseInfo parses redis info for use by redis-master.
func parseInfo(info string) (redis map[string]string) {
	var lines = strings.Split(info, "\r\n")
	redis = make(map[string]string)
	for _, line := range lines {
		if len(line) > 0 {
			if strings.HasPrefix(line, "#") {
				continue
			} else {
				splits := strings.SplitN(line, ":", 2)
				redis[splits[0]] = splits[1]
			}
		}
	}
	return redis
}
