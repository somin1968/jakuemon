package jakuemon

import (
	"context"
	"encoding/json"
	"google.golang.org/appengine/log"
	"net/http"
)

type errorResponse struct {
	Message string      `json:"message"`
	Debug   interface{} `json:"debug"`
}

func respond(ctx context.Context, w http.ResponseWriter, status int, attributes interface{}) {
	if status/100 >= 5 {
		log.Errorf(ctx, "%#v", attributes)
	} else if status/100 >= 3 {
		log.Infof(ctx, "%#v", attributes)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(attributes)
}
