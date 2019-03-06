package jakuemon

import (
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	// "google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"net/http"
)

func extCacheFlushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := memcache.Flush(ctx)
	if err == nil {
		fmt.Fprintln(w, "Cache flushed, all keys dropped.")
	} else {
		fmt.Fprintln(w, "Failed.")
	}
}

func extHandler(r *mux.Router) {
	s := r.PathPrefix("/ext").Subrouter()
	s.HandleFunc("/cache/flush/", extCacheFlushHandler).Methods("GET")
}
