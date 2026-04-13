package api

import (
	"app/pkg/ratelimiter"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type limiterAPIResponse struct {
	UserID  int    `json:"userID"`
	Message string `json:"message"`
}

type RateLimiter interface {
	Allow(ctx context.Context) (ratelimiter.LimitResult, error)
}

func RateLimiterAPIHandler(limiter *ratelimiter.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// queryParams := req.URL.Query()
		// userIDArr, ok := queryParams["userID"]

		// if !ok {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	response := limiterAPIResponse{
		// 		UserID:  0,
		// 		Message: "user id is not sent!",
		// 	}

		// 	responseJson, err := json.Marshal(response)
		// 	if err != nil {
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		log.Fatal()
		// 		return
		// 	}

		// 	w.Write(responseJson)
		// 	return
		// }

		//ctx := context.WithValue(context.Background(), "userID", userIDArr[0])
		res, err := limiter.Allow(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal()
			return
		}

		responseJson, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal()
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(responseJson)
	}
}

func RateLimitAPIHandler(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	userIDArr, ok := queryParams["userID"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response := limiterAPIResponse{
			UserID:  0,
			Message: "user id is not sent!",
		}

		responseJson, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal()
			return
		}

		w.Write(responseJson)
		return
	}

	userID, err := strconv.Atoi(userIDArr[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal()
		return
	}
	response := limiterAPIResponse{
		UserID:  int(userID),
		Message: "user retrieved Successfully!",
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal()
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
}
