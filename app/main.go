package main

import (
	"crypto/tls"
	"encoding/json"

	"fmt"

	"github.com/go-redis/redis"

	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var redisHost = os.Getenv("REDIS_HOST") // This should include the port which is most of the time 6379
var redisPassword = os.Getenv("REDIS_PASSWORD")
var redisTLSEnabled = os.Getenv("REDIS_TLS")
var redisTLSEnabledFlag = false
var client *redis.Client

type Input struct {
	ID     string
	Value  string
	Stored bool
}

type Result struct {
	ID        string
	Input     string
	Output    string
	Processed bool
}

type Results struct {
	Results []Result
}

var results Results

var inputs Inputs

type Inputs struct {
	Inputs []Input
}

func main() {

	if redisTLSEnabled != "" && redisTLSEnabled != "false" {
		redisTLSEnabledFlag = true
	}

	if !redisTLSEnabledFlag {
		client = redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: redisPassword,
			DB:       0,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: redisPassword,
			DB:       0,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		})
	}

	r := mux.NewRouter()
	r.HandleFunc("/", UIHandler).Methods("GET")
	r.HandleFunc("/info", InfoHandler).Methods("GET")
	r.HandleFunc("/store", StoreHandler).Methods("POST")
	r.HandleFunc("/values", GetValuesHandler).Methods("GET")

	r.HandleFunc("/clear", ClearHandler).Methods("DELETE")

	log.Printf("Strange app Started in port 8080!")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func InfoHandler(writer http.ResponseWriter, request *http.Request) {
	respondWithJSON(writer, http.StatusOK, "{ 'app': 'OK' }")
}

func StoreHandler(writer http.ResponseWriter, request *http.Request) {
	value := request.URL.Query().Get("value")
	err := client.LPush("values", string(value)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

}

func UIHandler(writer http.ResponseWriter, request *http.Request) {
	fileBytes, err := ioutil.ReadFile(os.Getenv("KO_DATA_PATH") + "/index.html")
	if err != nil {
		panic(err)
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html")
	writer.Write(fileBytes)
}

func ClearHandler(writer http.ResponseWriter, request *http.Request) {

	deleted, err := client.Del("values").Result()
	if err != nil {
		fmt.Println(err)
	}
	respondWithJSON(writer, http.StatusOK, deleted)
}

func GetValuesHandler(writer http.ResponseWriter, request *http.Request) {

	resultsFromRedis, err := client.LRange("values", 0, -1).Result()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	var evaluations []string
	for _, r := range resultsFromRedis {
		evaluations = append(evaluations, r)
	}

	respondWithJSON(writer, http.StatusOK, evaluations)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
