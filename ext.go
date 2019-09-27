package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func extCacheFlushHandler(w http.ResponseWriter, r *http.Request) {
	// err := memcache.Flush(ctx)
	// if err == nil {
	// 	respond(ctx, w, http.StatusOK, "Cache flushed, all keys dropped.")
	// } else {
	// 	respond(ctx, w, http.StatusBadGateway, "Failed.")
	// }
}

func extHandler(r *mux.Router) {
	s := r.PathPrefix("/ext").Subrouter()
	s.HandleFunc("/cache/flush/", extCacheFlushHandler).Methods("GET")
}
