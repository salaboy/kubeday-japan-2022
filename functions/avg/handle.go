package function

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

var redisHost = os.Getenv("REDIS_HOST") // This should include the port which is most of the time 6379
var redisPassword = os.Getenv("REDIS_PASSWORD")
var redisTLSEnabled = os.Getenv("REDIS_TLS")
var redisTLSEnabledFlag = false
var client *redis.Client

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
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

	resultsFromRedis, err := client.LRange("values", 0, -1).Result()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	var total int
	var count int
	for _, r := range resultsFromRedis {
		intVar, _ := strconv.Atoi(r)
		total += intVar
		count++
	}

	var avg float64
	avg = float64(total / count)
	respondWithJSON(res, http.StatusOK, avg)

	//respondWithJSON(res, http.StatusOK, "OK")

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
