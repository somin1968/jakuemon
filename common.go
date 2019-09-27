package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Message string      `json:"message"`
	Debug   interface{} `json:"debug"`
}

func respond(w http.ResponseWriter, status int, attributes interface{}) {
	if status/100 >= 5 {
		log.Printf("%#v", attributes)
	} else if status/100 >= 3 {
		log.Printf("%#v", attributes)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(attributes)
}
