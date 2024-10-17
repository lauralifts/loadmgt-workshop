package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 10000

	doPings()
}

func doPings() {
	hcUrl := os.Getenv("HC_URL")
	freq := os.Getenv("FREQUENCY")
	sleepytime, err := time.ParseDuration(freq)

	if err != nil {
		log.Fatalf("error: %+v", err)
	}

	for {
		res, err := http.Get(hcUrl)
		if err != nil {
			log.Printf("Error: %v", err)
		} else if res.StatusCode != http.StatusOK {
			log.Printf("Status %d received", res.StatusCode)
		} else {
			log.Printf("OK")
		}

		time.Sleep(sleepytime)
	}
}
