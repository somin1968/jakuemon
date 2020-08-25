package main

import (
	"fmt"
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
	buf, err := ioutil.ReadFile(SECRET_YAML_PATH)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(buf, &setting)
	r := mux.NewRouter().StrictSlash(true)
	apiHandler(r)
	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
