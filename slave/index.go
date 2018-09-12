package slave

import (
	"fmt"
	"net/http"
)

// Index : Method catering to the default route.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, welcome to redis-master Cache index route.")
}
