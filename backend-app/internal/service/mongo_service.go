package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ...existing code...
/*
func SaveDetailedAthlete(athlete *strava.DetailedAthlete) error {
	if !mongoInitialized {
		return fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athletes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, athlete)
	if err != nil {
		return fmt.Errorf("Failed to insert athlete: %v", err)
	}

	log.Println("Athlete saved to MongoDB")
	return nil
}
*/

func SaveCustomAthlete(athlete CustomAthlete) error {
	if !mongoInitialized {
		return fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athletes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, athlete)
	if err != nil {
		return fmt.Errorf("Failed to insert athlete: %v", err)
	}

	log.Println("Athlete saved to MongoDB")
	return nil
}

func SaveAthleteJSON(athleteJSON []byte) error {
	if !mongoInitialized {
		return fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athletes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var athlete map[string]interface{}
	err := json.Unmarshal(athleteJSON, &athlete)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal athlete JSON: %v", err)
	}

	_, err = collection.InsertOne(ctx, athlete)
	if err != nil {
		return fmt.Errorf("Failed to insert athlete: %v", err)
	}

	log.Println("Athlete saved to MongoDB")
	return nil
}

func SaveAthleteZonesJSON(athleteZoneJSON []byte) error {
	if !mongoInitialized {
		return fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athlete_zones")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var athleteZones map[string]interface{}
	err := json.Unmarshal(athleteZoneJSON, &athleteZones)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal athlete zones JSON: %v", err)
	}

	_, err = collection.InsertOne(ctx, athleteZones)
	if err != nil {
		return fmt.Errorf("Failed to insert athlete zones: %v", err)
	}

	log.Println("Athlete zones saved to MongoDB")
	return nil
}

func SaveAthleteStatsJSON(athleteStatsJSON []byte) error {
	if !mongoInitialized {
		return fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athlete_stats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var athleteStats map[string]interface{}
	err := json.Unmarshal(athleteStatsJSON, &athleteStats)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal athlete stats JSON: %v", err)
	}

	_, err = collection.InsertOne(ctx, athleteStats)
	if err != nil {
		return fmt.Errorf("Failed to insert athlete stats: %v", err)
	}

	log.Println("Athlete stats saved to MongoDB")
	return nil
}

// ...existing code...
