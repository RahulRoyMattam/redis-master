# redis-master

The redis-master project is an attempt at replicating the behaviour of how Redis Datastore manages its redis slave instances. This project is written in golang and is a simplistic implementation of managing redis slave instances.

You can clone the project, host it as a go server and scale it as many go instances as you want.

You have to set the environment variables by editing the env package of the redis-master go project so that the redis-master instance can start tracking your redis slave instances.

# Index

* [Redis Commands Supported](#redis-commands-supported)
    * [GET](#get)
    * [SET](#set)
    * [DEL](#del)
    * [EXPIRE](#expire)
    * [EXISTS](#exists)
    * [INFO](#info)
    * [KEYS](#keys)
    * [FLUSHALL](#flushall)

# Redis Commands Supported

### [GET](https://redis.io/commands/get)

A GET key request forwarded to a redis-master instance triggers a concurrent request for the key to all available redis instances managed by the redis-master instance. The redis-master instance waits till all the redis slave instances reply with the request status for the GET request and returns an appropriate error message or the data requested.

URL :
`GET  /get/{key}`

### [SET](https://redis.io/commands/set)

A SET Key request saves your string value to the best redis slave instance which is available to the redis-master instance. The redis-master instance responds to your web request with the message received from REDIS on executing the SET command.

The redis-master instance first deletes any key value object stored under the key specified in the SET request before issuing the SET request. This is to prevent duplicate keys being stored to different redis slave instances, as duplicate keys significantly slows down redis-master GET request if the key is found in multiple Redis slave instances, especially if the object cached under the key is huge since the redis-master downloads the entire object from multiple redis slave instances for the GET request.

For safety concerns, if expire is set to a negative value or zero, expire by default is set to 7 days. You can change this setting by modifying how the Set() in the basic_crud.go file of the redis package manages the EX flag.

URL :
`POST   /set`

POST BODY :
```json
{
    "key" : "keyName",
    "value" : "valueString",
    "expire" : "(key expire time in integer seconds)"
}
```
Example :
```json
{
    "key" : "rahul",
    "value" : "hello world!",
    "expire" : 3600
}
```

(More documentation on the way... Please raise an issue for clarifications)