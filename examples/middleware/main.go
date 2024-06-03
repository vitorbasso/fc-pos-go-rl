package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"ratel/pkg/ratel"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	godotenv.Load()
	redisAddr := "localhost:6379"
	addr, ok := os.LookupEnv("REDIS_ADDR")
	if ok {
		redisAddr = addr
	}
	ratel := ratel.Middleware(
		ratel.WithEnvRules(),
		ratel.WithRedisRequestStore(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		}),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: ratel(mux),
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
