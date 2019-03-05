package jakuemon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
)

const (
	JSON_FILE_PATH string = "jakuemon-235b6f83a3e3.json"
	GSAPI_SCOPE    string = "https://www.googleapis.com/auth/spreadsheets"
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

func getClient(ctx context.Context) (*http.Client, error) {
	json, err := ioutil.ReadFile(JSON_FILE_PATH)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(json, GSAPI_SCOPE)
	if err != nil {
		return nil, err
	}
	return config.Client(ctx), nil
}

func apiSheetListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client, err := getClient(ctx)
	if err != nil {
		respond(ctx, w, http.StatusBadGateway, errorResponse{
			Message: "認証に失敗しました。",
			Debug:   fmt.Sprintf("%v", err),
		})
		return
	}
	service, err := sheets.New(client)
	if err != nil {
		respond(ctx, w, http.StatusBadGateway, errorResponse{
			Message: "認証に失敗しました。",
			Debug:   fmt.Sprintf("%v", err),
		})
		return
	}
	vars := mux.Vars(r)
	resp, err := service.Spreadsheets.Values.Get(vars["id"], vars["range"]).Do()
	if err == nil {
		respond(ctx, w, http.StatusOK, resp.Values)
	} else {
		respond(ctx, w, http.StatusBadGateway, errorResponse{
			Message: "スプレッドシートのデータが取得できませんでした。",
			Debug:   fmt.Sprintf("%v", err),
		})
	}
}

func init() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/api/sheets/{id}/{range}/", apiSheetListHandler).Methods("GET")
	http.Handle("/", r)
}
