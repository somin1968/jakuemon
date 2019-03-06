package jakuemon

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter().StrictSlash(true)
	apiHandler(r)
	extHandler(r)
	http.Handle("/", r)
}
