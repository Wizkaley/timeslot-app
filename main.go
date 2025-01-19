package main

import (
	"fmt"
	"timeslot-app/app"
	_ "timeslot-app/docs"
	"timeslot-app/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	// swagger embed files
)

// @title           OpenAPI Time Slot App API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	app, err := app.InitApp()
	if err != nil {
		panic(err)
	}
	fmt.Println(app.Config)

	router := NewRouter(app)
	defer app.DB.Close()
	router.Run(":8000")
}

func NewRouter(app *app.App) *gin.Engine {

	r := gin.Default()
	r.Use(middleware.ErrorHandler)
	v1 := r.Group("/api/v1")
	{
		timeslot := v1.Group("/timeslots")
		timeslot.POST("", app.TimeslotService.CreateTimeSlot)
		timeslot.GET("/:username", app.TimeslotService.GetTimeSlotsByUserName)
		timeslot.GET("/recommend", app.TimeslotService.RecommendSlots)
		timeslot.DELETE("/:username", app.TimeslotService.DeleteTimeSlotsByUserName)
	}

	{
		users := v1.Group("/users")
		users.POST("", app.UserService.CreateUser)
	}

	{
		events := v1.Group("/events")
		events.POST("", app.EventService.CreateEvent)
		events.GET("/:username", app.EventService.GetEventsForUser)
		// events.GET("/{eventID}", app.EventService.GetEvent)
		events.DELETE("/:eventID", app.EventService.DeleteEvent)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
