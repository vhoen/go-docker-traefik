package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

func main() {
	Timeout := time.Minute * 1

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup

		forecast := ""
		temperature := ""

		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := http.Get("http://go-weather-forecast:8000/go-weather-forecast/forecast")
			if err != nil {
				logrus.Fatal(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Fatal(err)
			}

			data := map[string]string{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				logrus.Fatal(err)
			}

			forecast = data["data"]
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := http.Get("http://go-weather-temperature:8000/go-weather-temperature/temperature")
			if err != nil {
				logrus.Fatal(err)
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Fatal(err)
			}

			data := map[string]string{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				logrus.Fatal(err)
			}

			temperature = data["data"]
		}()

		wg.Wait()

		output := fmt.Sprintf("Forecast: %s\n", forecast)
		output += fmt.Sprintf("Temp: %s\n", temperature)

		w.Write([]byte(output))
	}).Methods(http.MethodGet, http.MethodOptions)

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
