package main

import (
	"log"
	"net/http"
	"unicode"

	services "github.com/cycleai/go-app/internal/service"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file from project root")
	}

	services.Initialize()

}

func isNumeric(s string) bool {
	if s == "-1" {
		return true
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {

		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/athlete/:id", func(c echo.Context) error {
		stravaAuthToken := c.Request().Header.Get("StravaAuthToken")
		id := c.Param("id")
		log.Print(id, stravaAuthToken)
		athleteJSON, err := services.GetLoggedInAthlete(stravaAuthToken)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete data")
		}
		return c.String(http.StatusOK, athleteJSON)
	})

	e.GET("/athlete/:id/zones", func(c echo.Context) error {
		stravaAuthToken := c.Request().Header.Get("StravaAuthToken")
		athleteZonesJSON, err := services.GetLoggedInAthleteZones(stravaAuthToken)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete Zone data")
		}
		return c.String(http.StatusOK, athleteZonesJSON)
	})

	e.GET("/athlete/:id/stats", func(c echo.Context) error {
		stravaAuthToken := c.Request().Header.Get("StravaAuthToken")
		id := c.Param("id")
		athleteZonesJSON, err := services.GetLoggedInAthleteStats(stravaAuthToken, id)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete Zone data")
		}
		return c.String(http.StatusOK, athleteZonesJSON)
	})

	e.GET("/athlete/:id/activities", func(c echo.Context) error {
		//stravaAuthToken := c.Request().Header.Get("StravaAuthToken")
		page := c.QueryParam("page")
		if page == "" || !isNumeric(page) {
			page = "1"
		}
		perPage := c.QueryParam("perPage")
		if perPage == "" || !isNumeric(perPage) {
			perPage = "10"
		}
		id := c.Param("id")
		athleteActivitiesJSON, err := services.GetLoggedInAthleteActivities(id, page, perPage)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete activities")
		}
		return c.String(http.StatusOK, athleteActivitiesJSON)
	})

	e.GET("/athlete/:athleteId/activities/:activityId", func(c echo.Context) error {
		stravaAuthToken := c.Request().Header.Get("StravaAuthToken")
		activityId := c.Param("activityId")
		athleteActivitiesJSON, err := services.GetLoggedInAthleteActivity(stravaAuthToken, activityId)
		log.Println(athleteActivitiesJSON)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete activities")
		}
		return c.String(http.StatusOK, athleteActivitiesJSON)
	})

	e.PATCH("/athlete/:athleteId/activities", func(c echo.Context) error {
		stravaAuthToken := c.Request().Header.Get("StravaAuthToken")
		athleteId := c.Param("athleteId")
		err := services.RefreshActivities(stravaAuthToken, athleteId)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error Refreshing athlete activities")
		}

		return c.NoContent(http.StatusOK)
	})

	e.POST("/embeddings", func(c echo.Context) error {
		request := new(EmbeddingRequest)
		if err := c.Bind(request); err != nil {
			return err
		}
		log.Println("Request:", request.Value)
		services.GetEmbeddings(request.Value)
		return c.String(http.StatusOK, "Embeddings generated")
	})

	e.GET("/vector-search", func(c echo.Context) error {
		query := c.QueryParam("query")
		result := services.ChatCompletionWithVector(query)
		return c.JSON(http.StatusOK, result)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

type EmbeddingRequest struct {
	Value string `json:"value" xml:"value" form:"value" query:"value"`
}
