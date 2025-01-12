package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var stravaInitialized bool
var accessToken string

var openaiInitialized bool
var openaiAPIKey string

var mongoClient *mongo.Client
var mongoInitialized bool

func Initialize() {

	log.Println("Initializing services")

	if stravaInitialized && openaiInitialized && mongoInitialized {
		return
	}

	accessToken = getEnvVar("STRAVA_ACCESS_TOKEN")
	stravaInitialized = true

	openaiAPIKey = getEnvVar("OPENAI_API_KEY")
	openaiInitialized = true

	initializeMongoDB()

	log.Println("Services initialized")

}

func initializeMongoDB() {
	dbURL := getEnvVar("MONGO_DB_URL")
	dbUsername := getEnvVar("MONGO_DB_USERNAME")
	dbPassword := getEnvVar("MONGO_DB_PASSWORD")

	dbURL = fmt.Sprintf(dbURL, dbUsername, dbPassword)
	log.Printf("DB Connection %s", dbURL)

	clientOptions := options.Client().ApplyURI(dbURL)

	var err error
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	mongoInitialized = true
	log.Println("MongoDB initialized")
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
