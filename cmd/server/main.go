package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"shortner/internal/cache"
	"shortner/internal/database"
	"shortner/internal/handlers"
)

var incomingRequests atomic.Uint64
var outgoingResponses atomic.Uint64

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingRequests.Add(1)
		next.ServeHTTP(w, r)
		outgoingResponses.Add(1)
	})
}

func metrics() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		in := incomingRequests.Swap(0)
		out := outgoingResponses.Swap(0)

		fmt.Printf("RPS: incoming=%d outgoing=%d\n", in, out)
	}
}

func main() {
	database.InitDB()
	cache.InitRedis()
	defer database.DB.Close()

	go metrics()
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", handlers.ShortenHandler)
	mux.HandleFunc("/app.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/app.js")
	})
	mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/style.css")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		handlers.RedirectHandler(w, r)
	})

	handler := metricsMiddleware(mux)

	fmt.Println("Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}
}
