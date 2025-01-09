package service

import (
	"log"
	"os"
)

var stravaInitialized bool
var accessToken string

var openaiInitialized bool
var openaiAPIKey string

func Initialize() {

	log.Println("Initializing services")

	if stravaInitialized && openaiInitialized {
		return
	}

	accessToken = getEnvVar("STRAVA_ACCESS_TOKEN")
	stravaInitialized = true

	openaiAPIKey = getEnvVar("OPENAI_API_KEY")
	openaiInitialized = true

	log.Println("Services initialized")

}

func getEnvVar(key string) string {
	log.Printf("Loading [%s]", key)
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%q not set", key)
	}
	log.Printf("Loaded [%s: %s]", key, value)
	return value
}
