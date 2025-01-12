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

	e.GET("/athlete", func(c echo.Context) error {
		athleteJSON, err := services.GetLoggedInAthlete()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete data")
		}
		return c.String(http.StatusOK, athleteJSON)
	})

	e.GET("/athlete/zones", func(c echo.Context) error {
		athleteZonesJSON, err := services.GetLoggedInAthleteZones()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete Zone data")
		}
		return c.String(http.StatusOK, athleteZonesJSON)
	})

	e.GET("/athlete/:id/stats", func(c echo.Context) error {
		id := c.Param("id")
		athleteZonesJSON, err := services.GetLoggedInAthleteStats(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete Zone data")
		}
		return c.String(http.StatusOK, athleteZonesJSON)
	})

	e.GET("/athlete/activities", func(c echo.Context) error {
		page := c.QueryParam("page")
		if page == "" || !isNumeric(page) {
			page = "1"
		}
		perPage := c.QueryParam("perPage")
		if perPage == "" || !isNumeric(perPage) {
			perPage = "10"
		}
		athleteActivitiesJSON, err := services.GetLoggedInAthleteActivities(page, perPage)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete activities")
		}
		return c.String(http.StatusOK, athleteActivitiesJSON)
	})

	e.GET("/athlete/activities/:id", func(c echo.Context) error {
		id := c.Param("id")
		athleteActivitiesJSON, err := services.GetLoggedInAthleteActivity(id)

		if err != nil {
			return c.String(http.StatusInternalServerError, "Error fetching athlete activities")
		}
		return c.String(http.StatusOK, athleteActivitiesJSON)
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
