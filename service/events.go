package service

import (
	"errors"
	"net/http"
	"timeslot-app/models"
	"timeslot-app/repository"
	"timeslot-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
)

type EventService struct {
	EventRepo    repository.EventRepo
	TimeslotRepo repository.TimeslotRepo
	UserRepo     repository.UserRepo
}

func NewEventService(db *pgx.Conn) *EventService {
	return &EventService{
		EventRepo:    repository.NewEventRepository(db),
		UserRepo:     repository.NewUserRepo(db),
		TimeslotRepo: repository.NewTimeslotRepository(db),
	}
}

// ShowAccount godoc
// @Summary      Create a Event
// @Description  Create a new Event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        body   body   	models.EventRequest   true "Create Event request body"
// @Success      200  {object}  string "Event created successfully"
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /events [post]
func (es *EventService) CreateEvent(ctx *gin.Context) {
	var eventReq models.EventRequest
	if err := ctx.BindJSON(&eventReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := models.Event{}

	eventID, _ := uuid.NewV4()
	event.ID = eventID

	if eventReq.EventOwner == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Event created by is required"})
		return
	}

	user, err := es.UserRepo.Get(eventReq.EventOwner)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Event owner does not exist"})
		return
	}
	event.EventOwner = user.ID

	if eventReq.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Event title is required"})
		return
	}

	event.Title = eventReq.Title

	if eventReq.Participants == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Event participants are required"})
		return
	}

	if eventReq.EventTimeSlot == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Event time slot is required"})
		return
	}
	startTime, endTime, valid := utils.ValidateAndFormatTimeStamp(eventReq.EventTimeSlot)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time slot format"})
		return
	}

	// check if the user requesting the event time has a timeslot

	userTimeSlots, err := es.TimeslotRepo.GetTimeSlotsByUserName(eventReq.EventOwner)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("error fetching user time slots")})
		return
	}
	if len(userTimeSlots) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("user does not have any time slots")})
		return
	}

	if !utils.SearchString(userTimeSlots, eventReq.EventTimeSlot) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("user does not have the requested time slot")})
		return
	}

	event.EventStartTime = startTime
	event.EventEndTime = endTime
	event.Participants = eventReq.Participants

	err = es.EventRepo.CreateEvent(event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Event created successfully"})
}

// func (es *EventService) GetEvents(ctx *gin.Context) {

// 	username := ctx.Query("username")
// 	events, err := es.EventRepo.GetEventsForUser(username)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"events": events})
// }

// ShowAccount godoc
// @Summary      Delete a Event
// @Description  Delete a Event
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        eventID   path   string   true  "Event ID"
// @Success      200  {object}  string "Event deleted successfully"
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /events/{eventID} [delete]
func (es *EventService) DeleteEvent(ctx *gin.Context) {
	eventID := ctx.Param("eventID")
	if eventID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}

	err := es.EventRepo.DeleteEvent(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// ShowAccount godoc
// @Summary      Get a Event
// @Description  Get Event by ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        eventID   path   string   true  "Event ID"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /events/{eventID} [get]
func (es *EventService) GetEvent(ctx *gin.Context) {
	eventID := ctx.Param("eventID")
	if eventID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Event ID is required"})
		return
	}

	event, err := es.EventRepo.GetEvent(eventID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"event": event})
}

// ShowAccount godoc
// @Summary      Get Events for a user
// @Description  Get Events for a user
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        username   path   string   true  "Username"
// @Success      200  {object}  []models.Event
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /events/{username} [get]
func (es *EventService) GetEventsForUser(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	events, err := es.EventRepo.GetEventsForUser(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"events": events})
}
