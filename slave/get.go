package slave

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rahulroymattam/redis-master/redis"

	redigo "github.com/garyburd/redigo/redis"
)

//GetValue for a key in redis
func GetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	ctx := r.Context()
	ret, err := redis.Get(ctx, key)
	if err == redigo.ErrNil {
		http.Error(w, "KEY NOT FOUND", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Redis fetch failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, ret)
}

//CheckExists a key in redis
func CheckExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	ctx := r.Context()
	ret, err := redis.Exists(ctx, key)
	if err != nil {
		http.Error(w, "Redis fetch failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, strconv.Itoa(ret))
}
