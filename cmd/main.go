package main

import (
	"app/internal/config"
	"app/api"
	"fmt"
	"log"
	"net/http"
	"app/pkg/ratelimiter"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	store := ratelimiter.NewInMemoryStore()
	limiter := ratelimiter.NewLimiter(store)

	http.HandleFunc("/", api.RateLimiterAPIHandler(limiter));
	err = http.ListenAndServe("localhost:" + config.App.Port, nil);
	if err != nil {
		fmt.Println(err.Error())
	}
}