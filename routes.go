package main

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/mux"
	"github.com/rahulroymattam/redis-master/slave"
)

// Route : struct to define the route parameters.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes : list of Route parameters for configuration.
type Routes []Route

// NewRouter : Initialize the router with the parameters provided.
func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	var routes Routes
	routes = redisMasterRoutes

	for _, route := range routes {
		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(raven.RecoveryHandler(route.HandlerFunc))
	}

	return r
}

// slaveRoutes : Initialize the routes for redis-master slave on initial setup
var redisMasterRoutes = Routes{
	Route{"Index", http.MethodGet, "/", slave.Index},
	Route{"GetValue", http.MethodGet, "/get/{key}", slave.GetValue},
	Route{"CheckExists", http.MethodGet, "/exists/{key}", slave.CheckExists},

	Route{"SetValue", http.MethodPost, "/set", slave.SetValue},

	Route{"DeleteKey", http.MethodDelete, "/del/{key}", slave.DeleteKey},
	Route{"SetExpire", http.MethodDelete, "/expire/{key}/{ttl}", slave.SetExpire},

	//Monitoring- Statging endpoints
	Route{"GetKeys", http.MethodGet, "/info/keys/{pattern}", slave.GetKeys},
	Route{"GetInfo", http.MethodGet, "/info/raw/{section}", slave.GetInfo},
	Route{"GetRedisInstances", http.MethodGet, "/info/connected-redis", slave.GetRedisInstances},
	Route{"GetBestRedisInstance", http.MethodGet, "/info/best-redis", slave.GetBestRedisInstance},
	Route{"FlushAllKeys", http.MethodDelete, "/info/flushall", slave.FlushAllKeys},
}
