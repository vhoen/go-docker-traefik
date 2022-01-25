package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

func main() {
	Timeout := time.Minute * 1

	router := mux.NewRouter()
	router = router.PathPrefix("/go-weather-forecast").Subrouter()
	router.StrictSlash(true)
	router.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(map[string]string{
			"data": "Cloudy with a chance of Meatball",
		})
		w.Write(resp)
	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal("Not Found")
		w.Write(resp)
	})

	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: Timeout,
		ReadTimeout:  Timeout,
		IdleTimeout:  Timeout,
		Handler:      router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatal(err)
		}
	}()

	// Process signals channel
	sigChannel := make(chan os.Signal, 1)

	// Graceful shutdown via SIGINT
	signal.Notify(sigChannel, os.Interrupt)

	logrus.Info("Service running...")
	<-sigChannel // Block until SIGINT received

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	server.Shutdown(ctx)

	logrus.Info("Http Service shutdown")

}
