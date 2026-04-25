package main

import (
	"app/internal/config"
	"app/api"
	"fmt"
	"net/http"
	"app/pkg/ratelimiter"
)

func main() {
	appConfig := config.GetConfig();

	store := ratelimiter.NewInMemoryStore()
	limiter := ratelimiter.NewLimiter(store)

	fmt.Println("port is: " + appConfig.Port)

	http.HandleFunc("/", api.RateLimiterAPIHandler(limiter));
	err := http.ListenAndServe("localhost:" + appConfig.Port, nil);
	if err != nil {
		fmt.Println(err.Error())
	}
}