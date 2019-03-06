package jakuemon

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/appengine"
	// "google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	JSON_FILE_PATH   string = "jakuemon-235b6f83a3e3.json"
	GSAPI_SCOPE      string = "https://www.googleapis.com/auth/spreadsheets"
	SPREADSHEET_ID   string = "1EphMrjBOswkOQNqgXDgUPTUwptYnFAnGMLu3v_FEHi8"
	MEMCACHE_EXPIRED int    = 60 * 60 * 12
)

var rangeDict map[string]string = map[string]string{
	"recents": "最新公演",
	"topics":  "お知らせ",
	"events":  "公演情報",
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
	category, ok := rangeDict[mux.Vars(r)["category"]]
	if ok == false {
		respond(ctx, w, http.StatusBadRequest, errorResponse{
			Message: "引数が不正です。",
			Debug:   nil,
		})
		return
	}
	var articles []map[string]string
	_, err := memcache.JSON.Get(ctx, category, &articles)
	if err == nil {
		respond(ctx, w, http.StatusOK, articles)
		return
	}
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
	resp, err := service.Spreadsheets.Values.Get(SPREADSHEET_ID, category).Do()
	if err != nil {
		respond(ctx, w, http.StatusBadGateway, errorResponse{
			Message: "スプレッドシートのデータが取得できませんでした。",
			Debug:   fmt.Sprintf("%v", err),
		})
		return
	}
	rows := resp.Values
	if len(rows) == 0 {
		respond(ctx, w, http.StatusOK, nil)
		return
	}
	header := rows[0]
	articles = make([]map[string]string, len(rows)-1)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		var article = map[string]string{}
		for j, col := range row {
			article[header[j].(string)] = col.(string)
		}
		articles[i-1] = article
	}
	memcache.JSON.Set(ctx, &memcache.Item{
		Key:        category,
		Object:     articles,
		Expiration: time.Duration(MEMCACHE_EXPIRED) * time.Second,
	})
	respond(ctx, w, http.StatusOK, articles)
}

func apiHandler(r *mux.Router) {
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/sheets/{category}/", apiSheetListHandler).Methods("GET")
}
