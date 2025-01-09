package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	strava "github.com/obalunenko/strava-api/client"
)

func GetLoggedInAthlete() (string, error) {

	apiClient, err := getAPIClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := createContext()
	defer cancel()

	athlete, err := apiClient.Athletes.GetLoggedInAthlete(ctx)
	if err != nil {
		log.Fatal(err)
	}

	athleteJSON, err := athlete.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	return string(athleteJSON), nil
}

func GetLoggedInAthleteZones() (string, error) {

	apiClient, err := getAPIClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := createContext()
	defer cancel()

	athleteZones, err := apiClient.Athletes.GetLoggedInAthleteZones(ctx)
	if err != nil {
		log.Fatal(err)
	}

	athleteZoneJSON, err := athleteZones.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	return string(athleteZoneJSON), nil
}

func GetLoggedInAthleteStats(id string) (string, error) {

	apiClient, err := getAPIClient()
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

	athleteStatsJSON, err := athleteStats.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	return string(athleteStatsJSON), nil
}

func LoadActivities() {
	apiClient, err := getAPIClient()
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
	apiClient, err := getAPIClient()
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

func getAPIClient() (*strava.APIClient, error) {
	if !stravaInitialized {
		return nil, fmt.Errorf("strava service is not initialized")
	}

	apiClient, err := strava.NewAPIClient(accessToken)
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
