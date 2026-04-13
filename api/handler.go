package api

import (
	"app/pkg/ratelimiter"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type RateLimiter interface {
	Allow(ctx context.Context) (ratelimiter.LimitResult, error)
}

func RateLimiterAPIHandler(limiter *ratelimiter.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		queryParams := req.URL.Query()
		userIDArr, ok := queryParams["userID"]

		if !ok {
			// TODO: Good to have a more general response dto
			// to include error detail
			w.WriteHeader(http.StatusNotFound)
			log.Println("not found userID")
			return
		}

		ctx := context.WithValue(req.Context(), "userID", userIDArr[0])
		req = req.WithContext(ctx)

		res, err := limiter.Allow(req.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		responseJson, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(responseJson)
	}
}
