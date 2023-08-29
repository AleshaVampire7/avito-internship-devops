package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")

	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Fatalf("Error converting REDIS_DB to an integer: %s", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	r := mux.NewRouter()

	r.HandleFunc("/get_key/{key}", GetKeyHandler(client)).Methods("GET")
	r.HandleFunc("/set_key", SetKeyHandler(client)).Methods("POST")
	r.HandleFunc("/del_key", DeleteKeyHandler(client)).Methods("DELETE")

	http.Handle("/", r)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetKeyHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		log.Printf("Fetching value for key: %s", key)

		val, err := client.Get(ctx, key).Result()
		if err == redis.Nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("Key '%s' not found", key)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error fetching value for key '%s': %s", key, err)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"value": val})
	}
}

func SetKeyHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received a set_key request")

		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error decoding request data:", err)
			return
		}

		key := data["key"]
		value := data["value"]

		err := client.Set(ctx, key, value, 0).Err()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error setting key-value:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Println("Key-value set successfully")
	}
}

func DeleteKeyHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received a del_key request")

		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error decoding request data:", err)
			return
		}

		key, exists := data["key"]
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Missing key in request data")
			return
		}

		_, err := client.Del(ctx, key).Result()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error deleting key:", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Println("Key deleted successfully")
	}
}
