package redis

import (
	"context"

	"github.com/garyburd/redigo/redis"
)

//Set a key-value pair to the redis database. returns an "OK" string if saved successfully
func Set(ctx context.Context, key string, value string, expire int) (ret string, err error) {
	args := []interface{}{key, value}
	if expire > 0 {
		args = append(args, "EX", expire)
	} else {
		//set default to 7 days
		args = append(args, "EX", 604800)
	}
	Del(ctx, key)
	ret, err = redis.String(doBest(ctx, "SET", args...))
	return ret, err
}

//Get a key-value pair from the redis database
func Get(ctx context.Context, key string) (value string, err error) {
	value = ""
	responses, err := doAll(ctx, "GET", key)
	if err != nil {
		return value, err
	}

	for _, res := range responses {
		value, err = redis.String(res.reply, res.err)
		if value != "" {
			return value, nil
		}
	}

	return value, nil
}

//Del : Deletes a key in the redis database. Returns the number of keys deleted.
func Del(ctx context.Context, key string) (deleted int, err error) {
	deleted = 0
	responses, err := doAll(ctx, "DEL", key)
	if err != nil {
		return deleted, err
	}

	for _, res := range responses {
		n, _ := redis.Int(res.reply, res.err)
		if n != 0 {
			deleted += n
		}
	}
	return deleted, nil
}

//Expire : sets the TTL of a key in the redis database. Returns 1 if timeout set successfully, else returns 0.
func Expire(ctx context.Context, key string, ttl int) (set int, err error) {
	set = 0
	responses, err := doAll(ctx, "EXPIRE", key, ttl)
	if err != nil {
		return set, err
	}

	for _, res := range responses {
		n, _ := redis.Int(res.reply, res.err)
		if n != 0 {
			set += n
		}
	}
	return set, nil
}

//Exists : check if a key exists in the redis database. Returns the number of redis instances the key exists in otherwise returns zero.
func Exists(ctx context.Context, key string) (set int, err error) {
	set = 0
	responses, err := doAll(ctx, "EXISTS", key)
	if err != nil {
		return set, err
	}

	for _, res := range responses {
		n, _ := redis.Int(res.reply, res.err)
		if n != 0 {
			set += n
		}
	}
	return set, nil
}
