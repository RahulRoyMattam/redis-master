package slave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rahulroymattam/redis-master/redis"
)

//SetValue to a key in redis
func SetValue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var m SetCommand
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read from body failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		http.Error(w, "Json Unmarshall failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	ret, err := redis.Set(ctx, m.Key, m.Value, m.Expire)
	if err != nil {
		http.Error(w, "Redis Save failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, ret)
}
