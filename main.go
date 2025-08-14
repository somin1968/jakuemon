package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

const SECRET_YAML_PATH = "secret.yaml"

var setting struct {
	Brevo struct {
		ApiKey string `yaml:"apikey"`
		Sender string `yaml:"sender"`
	} `yaml:"brevo"`
}

func main() {
	buf, err := os.ReadFile(SECRET_YAML_PATH)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(buf, &setting)
	r := mux.NewRouter().StrictSlash(true)
	apiHandler(r)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
