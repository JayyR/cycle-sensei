package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	strava "github.com/obalunenko/strava-api/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var rateLimiter = NewRateLimiter(6, 15*time.Minute)

func GetLoggedInAthlete(stravaAuthToken string) (string, error) {
	apiClient, err := getAPIClient(stravaAuthToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := createContext()
	defer cancel()

	athlete, err := apiClient.Athletes.GetLoggedInAthlete(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Marshal athlete to JSON and save to MongoDB
	athleteJSON, err := athlete.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	err = SaveAthleteJSON(athleteJSON)
	if err != nil {
		log.Fatal(err)
	}

	return string(athleteJSON), nil
}

func GetLoggedInAthleteZones(stravaAuthToken string) (string, error) {
	apiClient, err := getAPIClient(stravaAuthToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := createContext()
	defer cancel()

	athleteZones, err := apiClient.Athletes.GetLoggedInAthleteZones(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Marshal athlete zones to JSON and save to MongoDB
	athleteZoneJSON, err := athleteZones.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	err = SaveAthleteZonesJSON(athleteZoneJSON)
	if err != nil {
		log.Fatal(err)
	}

	return string(athleteZoneJSON), nil
}

func GetLoggedInAthleteStats(stravaAuthToken string, id string) (string, error) {
	apiClient, err := getAPIClient(stravaAuthToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := createContext()
	defer cancel()

	athleteID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	athleteStats, err := apiClient.Athletes.GetStats(ctx, athleteID)
	if err != nil {
		log.Fatal(err)
	}

	// Marshal athlete stats to JSON and save to MongoDB
	athleteStatsJSON, err := athleteStats.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	err = SaveAthleteStatsJSON(athleteStatsJSON)
	if err != nil {
		log.Fatal(err)
	}

	return string(athleteStatsJSON), nil
}

func GetLoggedInAthleteActivities(athleteId string, p string, pp string) (string, error) {
	log.Printf("Getting activities for athlete %s", athleteId)
	activities, err := GetActivitiesFromDB(athleteId, p, pp)
	if err != nil {
		return "", fmt.Errorf("Failed to get activities from DB: %v", err)
	}

	activityJSON, err := json.Marshal(activities)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal activities to JSON: %v", err)
	}
	log.Printf("number of activities loaded ", len(activities))
	return string(activityJSON), nil
}

func GetLoggedInAthleteActivity(stravaAuthToken string, activityId string) (string, error) {
	log.Printf("Getting activity %s", activityId)
	if !mongoInitialized {
		return "", fmt.Errorf("MongoDB is not initialized")
	}

	collection := mongoClient.Database("cycle_sensei").Collection("activities")
	ctx, cancel := createContext()
	defer cancel()

	activityID, err := strconv.Atoi(activityId)
	if err != nil {
		return "", fmt.Errorf("Invalid activity ID: %v", err)
	}

	filter := bson.M{
		"id": activityID,
	}

	var activity map[string]interface{}
	err = collection.FindOne(ctx, filter).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("No activity found for activity ID %d", activityID)
		}
		return "", fmt.Errorf("Failed to get activity: %v", err)
	}
	log.Println("Activity found")

	activityJSON, err := json.Marshal(activity)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal activity to JSON: %v", err)
	}
	log.Printf("Activity loaded %s", activityJSON)

	return string(activityJSON), nil
}

func RefreshActivities(token string, athleteId string) error {
	go func() {
		lastSyncedDate, err := GetLastSyncedDate(athleteId)
		if err != nil {
			log.Println("Error getting last synced date:", err)
			updateSyncStatus(athleteId, "failed")
			return
		}

		updateSyncStatus(athleteId, "processing")

		activities := loadActivities(token, 10, lastSyncedDate)
		err = SaveActivities(athleteId, activities)
		if err != nil {
			log.Println("Error saving activities:", err)
			updateSyncStatus(athleteId, "failed")
			return
		}

		err = UpdateLastSyncedDate(athleteId, time.Now().Unix())
		if err != nil {
			log.Println("Error updating last synced date:", err)
			updateSyncStatus(athleteId, "failed")
			return
		}

		updateSyncStatus(athleteId, "completed")
	}()

	return nil
}

func GetLoggedInAthleteSelectedActivities(athleteId string, activityIds []int64) (string, error) {
	log.Printf("Getting selected activities for athlete %s", athleteId)
	activities, err := GetSelectedActivitiesFromDB(athleteId, activityIds)
	if err != nil {
		return "", fmt.Errorf("Failed to get selected activities from DB: %v", err)
	}

	activityJSON, err := json.Marshal(activities)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal selected activities to JSON: %v", err)
	}
	log.Printf("Number of selected activities loaded: %d", len(activities))
	return string(activityJSON), nil
}

func loadActivities(token string, max int, fetchAfter int64) []string {
	apiClient, err := getAPIClient(token)
	if err != nil {
		log.Fatal(err)
	}

	page := int32(1)
	perPage := int32(30)
	maxPages := max
	maxActivities := maxPages * int(perPage)
	loadedActivities := 0
	var activitiesArray []string // Initialize an array to store activity JSON

	for page <= int32(maxPages) {

		if loadedActivities >= maxActivities {
			break
		}
		options := strava.GetLoggedInAthleteActivitiesOpts{
			Page:    &page,
			PerPage: &perPage,
			After:   &fetchAfter,
		}

		ctx, cancel := createContext()
		defer cancel()

		// Use rate limiter to ensure we do not exceed the rate limit
		err = rateLimiter.Call(func() error {
			activities, err := apiClient.Activities.GetLoggedInAthleteActivities(ctx, options)
			if err != nil {
				return err
			}

			if len(activities) == 0 {
				return nil
			}
			loadedActivities += len(activities)

			fmt.Printf("Page %d\n", page)
			for _, activity := range activities {
				fmt.Println(activity.Athlete.ID, activity.Name, activity.StartDate)
				activityJSON := getActivity(activity.ID, token)
				fmt.Println(string(activityJSON))
				activitiesArray = append(activitiesArray, string(activityJSON)) // Add activity JSON to array
			}

			// Periodically save activities to MongoDB
			if len(activitiesArray) > 0 {
				err = SaveActivities(token, activitiesArray)
				if err != nil {
					log.Println("Error saving activities:", err)
					return err
				}
				activitiesArray = activitiesArray[:0] // Clear the array after saving
			}

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		page++
	}

	return activitiesArray
}

func getActivity(id int64, token string) []byte {
	log.Printf("Getting activity %d", id)
	apiClient, err := getAPIClient(token)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := createContext()
	defer cancel()

	activity, err := apiClient.Activities.GetActivityById(ctx, id)
	if err != nil {
		log.Fatal(err)
	}

	activityJSON, err := activity.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	return activityJSON
}

func getAPIClient(stravaAuthToken string) (*strava.APIClient, error) {
	if !stravaInitialized {
		return nil, fmt.Errorf("strava service is not initialized")
	}

	apiClient, err := strava.NewAPIClient(stravaAuthToken)
	if err != nil {
		log.Fatal(err)
	}

	if apiClient == nil {
		log.Fatal("apiClient is not initialized")
	}

	return apiClient, nil
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// Define a new struct that omits the segment_efforts field
type ActivityWithoutSegmentEfforts struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	// Add other fields as needed
}

// Define a new struct that includes only the required fields
type ActivitySummary struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	SportType string    `json:"sport_type"`
}

// Define a custom struct to store athlete data
type CustomAthlete struct {
	ID                    int64     `json:"id"`
	Username              string    `json:"username"`
	ResourceState         int32     `json:"resource_state"`
	Firstname             string    `json:"firstname"`
	Lastname              string    `json:"lastname"`
	City                  string    `json:"city"`
	State                 string    `json:"state"`
	Country               string    `json:"country"`
	Sex                   string    `json:"sex"`
	Premium               bool      `json:"premium"`
	Summit                bool      `json:"summit"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	FollowerCount         int64     `json:"follower_count"`
	FriendCount           int64     `json:"friend_count"`
	MeasurementPreference string    `json:"measurement_preference"`
	FTP                   int64     `json:"ftp"`
	Weight                float32   `json:"weight"`
}

type RateLimiter struct {
	requests  int
	resetTime time.Time
	limit     int
	interval  time.Duration
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:  0,
		resetTime: time.Now().Add(interval),
		limit:     limit,
		interval:  interval,
	}
}

func (r *RateLimiter) Call(fn func() error) error {
	if r.requests >= r.limit {
		waitTime := time.Until(r.resetTime)
		if waitTime > 0 {
			time.Sleep(waitTime)
		}
		r.reset()
	}
	r.requests++
	return fn()
}

func (r *RateLimiter) reset() {
	r.requests = 0
	r.resetTime = time.Now().Add(r.interval)
}
