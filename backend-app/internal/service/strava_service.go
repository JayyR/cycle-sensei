package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	strava "github.com/obalunenko/strava-api/client"
)

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

func GetLoggedInAthleteActivities(stravaAuthToken string, p string, pp string) (string, error) {
	apiClient, err := getAPIClient(stravaAuthToken)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := createContext()
	defer cancel()

	pageInt, err := strconv.ParseInt(p, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	page := int32(pageInt)
	perPageInt, err := strconv.ParseInt(pp, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	perPage := int32(perPageInt)
	options := strava.GetLoggedInAthleteActivitiesOpts{
		Page:    &page,
		PerPage: &perPage,
	}

	activities, err := apiClient.Activities.GetLoggedInAthleteActivities(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	activitySummaries := make([]ActivitySummary, len(activities))
	for i, activity := range activities {
		activitySummaries[i] = ActivitySummary{
			ID:        activity.ID,
			Name:      activity.Name,
			StartDate: time.Time(activity.StartDate),
			SportType: string(activity.SportType),
		}
	}
	activityJSON, err := json.Marshal(activitySummaries)
	if err != nil {
		log.Fatal(err)
	}
	return string(activityJSON), nil
}

func GetLoggedInAthleteActivity(stravaAuthToken string, id string) (string, error) {
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

	options := strava.GetActivityByIdOpts{
		IncludeAllEfforts: func(b bool) *bool { return &b }(false),
	}

	activity, err := apiClient.Activities.GetActivityById(ctx, athleteID, options)
	if err != nil {
		log.Fatal(err)
	}

	activityWithoutSegmentEfforts := ActivityWithoutSegmentEfforts{
		ID:        activity.ID,
		Name:      activity.Name,
		StartDate: time.Time(activity.StartDate),
		// Copy other fields as needed
	}

	activityJSON, err := json.Marshal(activityWithoutSegmentEfforts)
	if err != nil {
		log.Fatal(err)
	}
	return string(activityJSON), nil
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

func LoadActivities() {
	apiClient, err := getAPIClient("")
	if err != nil {
		log.Fatal(err)
	}

	page := int32(1)
	perPage := int32(3)
	maxPages := 1
	maxActivities := maxPages * int(perPage)
	loadedActivities := 0

	for page <= int32(maxPages) {

		if loadedActivities >= maxActivities {
			break
		}
		options := strava.GetLoggedInAthleteActivitiesOpts{
			Page:    &page,
			PerPage: &perPage,
		}

		ctx, cancel := createContext()
		defer cancel()
		activities, err := apiClient.Activities.GetLoggedInAthleteActivities(ctx, options)
		if err != nil {
			log.Fatal(err)
		}

		if len(activities) == 0 {
			break
		}
		loadedActivities += len(activities)

		fmt.Printf("Page %d\n", page)
		for _, activity := range activities {
			fmt.Println(activity.Athlete.ID, activity.Name, activity.StartDate)
			GetActivity(activity.ID)
		}
		page++
	}
}

func GetActivity(id int64) {
	log.Printf("Getting activity %d", id)
	apiClient, err := getAPIClient("")
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
	fmt.Println(string(activityJSON))
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
