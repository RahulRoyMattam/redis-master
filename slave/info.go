package slave

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rahulroymattam/redis-master/redis"
)

//GetInfo about the redis server
func GetInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	section := vars["section"]
	ctx := r.Context()
	ret := redis.Info(ctx, section)
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		http.Error(w, "Error encoding result to json: "+err.Error(), http.StatusInternalServerError)
	}
}

//GetKeys saved in the redis server
func GetKeys(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pattern := vars["pattern"]
	ctx := r.Context()
	ret, err := redis.Keys(ctx, pattern)
	if err != nil {
		http.Error(w, "Redis KEYS failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		http.Error(w, "Error encoding result to json: "+err.Error(), http.StatusInternalServerError)
	}
}

//FlushAllKeys is to clear the redis datastore
func FlushAllKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status, err := redis.FlushAll(ctx)
	if err != nil {
		http.Error(w, "Redis KEYS failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "Error encoding result to json: "+err.Error(), http.StatusInternalServerError)
	}
}

//GetRedisInstances which are available to the current redis-master instance
func GetRedisInstances(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(redis.Instances); err != nil {
		http.Error(w, "Error encoding result to json: "+err.Error(), http.StatusInternalServerError)
	}
}

//GetBestRedisInstance which is available to the current redis-master instance
//This endpoint is critical to the load balancing element of the redis-master master instance.
func GetBestRedisInstance(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(redis.BestInstance); err != nil {
		http.Error(w, "Error encoding result to json: "+err.Error(), http.StatusInternalServerError)
	}
}
