package slave

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rahulroymattam/redis-master/redis"
)

//DeleteKey for a key-value entry in redis
func DeleteKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	ctx := r.Context()
	ret, err := redis.Del(ctx, key)
	if err != nil {
		http.Error(w, "Redis fetch failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, strconv.Itoa(ret))
}

//SetExpire for a key-value entry in redis
func SetExpire(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	ttl, err := strconv.Atoi(vars["ttl"])
	if err != nil {
		http.Error(w, "String conversion failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	ret, err := redis.Expire(ctx, key, ttl)
	if err != nil {
		http.Error(w, "Redis fetch failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, strconv.Itoa(ret))
}
