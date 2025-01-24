package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ...existing code...

const mongoNotInitializedErr = "MongoDB is not initialized"

func SaveCustomAthlete(athlete CustomAthlete) error {
	if !mongoInitialized {
		return errors.New(mongoNotInitializedErr)
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athletes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, athlete)
	if err != nil {
		return fmt.Errorf("failed to insert athlete: %v", err)
	}

	log.Println("Athlete saved to MongoDB")
	return nil
}

func SaveAthleteJSON(athleteJSON []byte) error {
	if !mongoInitialized {
		return errors.New(mongoNotInitializedErr)
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athletes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var athlete map[string]interface{}
	err := json.Unmarshal(athleteJSON, &athlete)
	if err != nil {
		return fmt.Errorf("failed to unmarshal athlete JSON: %v", err)
	}

	_, err = collection.InsertOne(ctx, athlete)
	if err != nil {
		return fmt.Errorf("failed to insert athlete: %v", err)
	}

	log.Println("Athlete saved to MongoDB")
	return nil
}

func SaveAthleteZonesJSON(athleteZoneJSON []byte) error {
	if !mongoInitialized {
		return errors.New(mongoNotInitializedErr)
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athlete_zones")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var athleteZones map[string]interface{}
	err := json.Unmarshal(athleteZoneJSON, &athleteZones)
	if err != nil {
		return fmt.Errorf("failed to unmarshal athlete zones JSON: %v", err)
	}

	_, err = collection.InsertOne(ctx, athleteZones)
	if err != nil {
		return fmt.Errorf("failed to insert athlete zones: %v", err)
	}

	log.Println("Athlete zones saved to MongoDB")
	return nil
}

func SaveAthleteStatsJSON(athleteStatsJSON []byte) error {
	if !mongoInitialized {
		return errors.New(mongoNotInitializedErr)
	}

	collection := mongoClient.Database("cycle_sensei").Collection("athlete_stats")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var athleteStats map[string]interface{}
	err := json.Unmarshal(athleteStatsJSON, &athleteStats)
	if err != nil {
		return fmt.Errorf("failed to unmarshal athlete stats JSON: %v", err)
	}

	_, err = collection.InsertOne(ctx, athleteStats)
	if err != nil {
		return fmt.Errorf("failed to insert athlete stats: %v", err)
	}

	log.Println("Athlete stats saved to MongoDB")
	return nil
}

func SaveActivities(athleteId string, activities []string) error {
	if !mongoInitialized {
		return errors.New(mongoNotInitializedErr)
	}

	collection := mongoClient.Database("cycle_sensei").Collection("activities")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, activityJSON := range activities {
		var activity map[string]interface{}
		err := json.Unmarshal([]byte(activityJSON), &activity)
		if err != nil {
			return fmt.Errorf("failed to unmarshal activity JSON: %v", err)
		}
		activity["athlete_id"] = athleteId

		_, err = collection.InsertOne(ctx, activity)
		if err != nil {
			return fmt.Errorf("failed to insert activity: %v", err)
		}
	}

	log.Println("Activities saved to MongoDB")
	return nil
}

func GetLastSyncedDate(athleteId string) (int64, error) {
	if !mongoInitialized {
		return 0, errors.New(mongoNotInitializedErr)
	}

	collection := mongoClient.Database("cycle_sensei").Collection("activities")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert athleteId to int
	athleteIdInt, err := strconv.Atoi(athleteId)
	if err != nil {
		return 0, fmt.Errorf("invalid athlete ID: %v", err)
	}

	filter := map[string]interface{}{
		"athlete.id": athleteIdInt,
	}
	options := options.FindOne().SetSort(map[string]interface{}{"start_date": -1})

	var result struct {
		StartDate string `bson:"start_date"`
	}
	err = collection.FindOne(ctx, filter, options).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No activities found for athlete:", athleteId)
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get last synced date: %v", err)
	}

	// Parse the start_date string to a timestamp
	startDate, err := time.Parse(time.RFC3339, result.StartDate)
	if err != nil {
		return 0, fmt.Errorf("failed to parse start_date: %v", err)
	}

	return startDate.Unix(), nil
}

func UpdateLastSyncedDate(athleteId string, lastSyncedDate int64) error {
	if !mongoInitialized {
		return errors.New(mongoNotInitializedErr)
	}

	collection := mongoClient.Database("cycle_sensei").Collection("sync_info")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, map[string]interface{}{"athlete_id": athleteId}, map[string]interface{}{
		"$set": map[string]interface{}{
			"last_synced_date": lastSyncedDate,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("failed to update last synced date: %v", err)
	}

	log.Println("Last synced date updated in MongoDB")
	return nil
}

func updateSyncStatus(athleteId string, status string) {
	if !mongoInitialized {
		log.Println(mongoNotInitializedErr)
		return
	}

	collection := mongoClient.Database("cycle_sensei").Collection("sync_info")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, map[string]interface{}{"athlete_id": athleteId}, map[string]interface{}{
		"$set": map[string]interface{}{
			"sync_status": status,
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		log.Println("Failed to update sync status:", err)
	}
}

func GetActivitiesFromDB(athleteId string, p string, pp string) ([]map[string]interface{}, error) {
	if !mongoInitialized {
		return nil, fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("activities")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	page, err := strconv.Atoi(p)
	if err != nil {
		return nil, fmt.Errorf("invalid page number: %v", err)
	}
	perPage, err := strconv.Atoi(pp)
	if err != nil {
		return nil, fmt.Errorf("invalid per page number: %v", err)
	}

	// Convert athleteId to int
	athleteIdInt, err := strconv.Atoi(athleteId)
	if err != nil {
		return nil, fmt.Errorf("invalid athlete ID: %v", err)
	}

	filter := map[string]interface{}{
		"athlete.id": athleteIdInt,
		"sport_type": map[string]interface{}{
			"$in": []string{"VirtualRide", "Ride", "GravelRide", "MountainBikeRide"},
		},
	}

	findOptions := options.Find().
		SetSort(map[string]interface{}{"start_date": -1}).
		SetProjection(map[string]interface{}{
			"id":         1,
			"name":       1,
			"start_date": 1,
			"sport_type": 1,
			"_id":        0,
		})

	if page != -1 && perPage != -1 {
		skip := (page - 1) * perPage
		limit := int64(perPage)
		findOptions.SetSkip(int64(skip)).SetLimit(limit)
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find activities: %v", err)
	}
	defer cursor.Close(ctx)

	var activities []map[string]interface{}
	for cursor.Next(ctx) {
		var activity map[string]interface{}
		err := cursor.Decode(&activity)
		if err != nil {
			return nil, fmt.Errorf("failed to decode activity: %v", err)
		}
		activities = append(activities, activity)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return activities, nil
}

func GetSelectedActivitiesFromDB(athleteId string, activityIds []int64) ([]map[string]interface{}, error) {
	if !mongoInitialized {
		return nil, fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("activities")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert athleteId to int
	athleteIdInt, err := strconv.Atoi(athleteId)
	if err != nil {
		return nil, fmt.Errorf("invalid athlete ID: %v", err)
	}

	filter := map[string]interface{}{
		"athlete.id": athleteIdInt,
		"id": map[string]interface{}{
			"$in": activityIds,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find activities: %v", err)
	}
	defer cursor.Close(ctx)

	var activities []map[string]interface{}
	for cursor.Next(ctx) {
		var activity map[string]interface{}
		err := cursor.Decode(&activity)
		if err != nil {
			return nil, fmt.Errorf("failed to decode activity: %v", err)
		}
		activities = append(activities, activity)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return activities, nil
}

// ...existing code...
