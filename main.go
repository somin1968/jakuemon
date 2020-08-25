package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const SECRET_YAML_PATH = "secret.yaml"

var setting struct {
	Sendgrid struct {
		ApiKey string `yaml:"apikey"`
		Sender string `yaml:"sender"`
	} `yaml:"sendgrid"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	buf, err := ioutil.ReadFile(SECRET_YAML_PATH)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(buf, &setting)
	r := mux.NewRouter().StrictSlash(true)
	apiHandler(r)
	http.Handle("/", r)
	loggingHandler := handlers.CombinedLoggingHandler(os.Stderr, r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, loggingHandler); err != nil {
		log.Fatal(err)
	}
}
