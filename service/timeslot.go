package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"timeslot-app/models"
	"timeslot-app/repository"
	"timeslot-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
)

type TimeslotServiceImplementaion struct {
	TimeslotRepo repository.TimeslotRepo
	UserRepo     repository.UserRepo
}

func Init(db *pgx.Conn) *TimeslotServiceImplementaion {
	service := new(TimeslotServiceImplementaion)
	service.TimeslotRepo = repository.NewTimeslotRepository(db)
	service.UserRepo = repository.NewUserRepo(db)
	return service
}

// ShowAccount godoc
// @Summary      Create a time slot
// @Description  Create time slot for a user
// @Tags         Timeslots
// @Accept       json
// @Produce      json
// @Param        body   body    models.UserTimeSlotRequest   true  "Timeslot request body"
// @Success      200  {object}  models.ServiceMessage
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /timeslot [post]
func (ts *TimeslotServiceImplementaion) CreateTimeSlot(ctx *gin.Context) {
	// create time slot

	var userTimeSlot models.UserTimeSlotRequest
	if err := ctx.BindJSON(&userTimeSlot); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorHelper("Invalid request body", err))
		return
	}

	// check if the user with the given name exists

	userFromDB, err := ts.UserRepo.Get(userTimeSlot.UserName)
	if err != nil {
		fmt.Println("error is ::", err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("Error fetching user", err))
		return
	}

	if userTimeSlot.TimeSlots == nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorHelper("Invalid request body", errors.New("time slots not provided")))
		return
	}

	// timeSlots := []string{}
	userTimeSlots := make([]models.TimeSlot, 0)
	// if yes validate the time slots
	for _, timeSlot := range userTimeSlot.TimeSlots {
		// validate the time slot
		// if not valid return error
		// if valid save the time slot
		if !utils.ValidateTimeStamp(timeSlot) {
			ctx.JSON(http.StatusBadRequest, utils.ErrorHelper("Invalid time slot format", errors.New("invalid time slot format")))
			return
		}
		tsID, err := uuid.NewV4()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("error generating a new uuid", err))
			return
		}
		userTimeSlots = append(userTimeSlots, models.TimeSlot{
			ID:       tsID,
			UserID:   userFromDB.ID,
			TimeSlot: timeSlot,
		})

	}
	// save the time slot for the user.
	err = ts.TimeslotRepo.Create(userTimeSlots)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("Error creating time slots", err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Timeslot created successfully"})
}

// ShowAccount godoc
// @Summary      Get a time slot
// @Description  Get time slot for a user by name
// @Tags         Timeslots
// @Accept       json
// @Produce      json
// @Param        username   path   string   true  "Timeslot request body"
// @Success      200  {object}  models.TimeSlotResponse
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /:username [get]
func (ts *TimeslotServiceImplementaion) GetTimeSlotsByUserName(ctx *gin.Context) {
	userName := ctx.Param("username")
	if userName == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorHelper("Invalid request body", errors.New("username not provided")))
		return
	}

	timeSlots, err := ts.TimeslotRepo.GetTimeSlotsByUserName(userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("Error fetching time slots", err))
		return
	}
	if len(timeSlots) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "no time slots found for the user"})
		return
	} else {
		tsr := models.TimeSlotResponse{
			UserName:  userName,
			TimeSlots: timeSlots,
		}
		json.NewEncoder(ctx.Writer).Encode(tsr)
	}
}

// ShowAccount godoc
// @Summary      Recommend time slots
// @Description  Recommend time slots for the given organizer and participants
// @Tags         Timeslots
// @Accept       json
// @Produce      json
// @Param        body   body   	models.RecommendSlotsRequest   true "Recommendation request body"
// @Success      200  {object}  models.RecommendSlotsResponse
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /recommend [get]
func (ts *TimeslotServiceImplementaion) RecommendSlots(ctx *gin.Context) {
	var recommendSlotsRequest models.RecommendSlotsRequest
	if err := ctx.BindJSON(&recommendSlotsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorHelper("Invalid request body", err))
		return
	}

	organizer := recommendSlotsRequest.Organizer
	participants := recommendSlotsRequest.Participants
	eventDuration := time.Duration(recommendSlotsRequest.EventDuration) * time.Minute
	// convert int to duration in minutes

	matchedSlots, partialMatchedSlots, err := ts.RecommendSlotsReconciler(ctx, organizer, participants, eventDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("Error recommending slots", err))
		return
	}
	resp := models.RecommendSlotsResponse{
		MatchedSlots: matchedSlots,
		PartialSlots: partialMatchedSlots,
	}
	err = json.NewEncoder(ctx.Writer).Encode(resp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("Error marshalling response", err))
		return
	}
}

func (ts *TimeslotServiceImplementaion) RecommendSlotsReconciler(ctx *gin.Context, organizer string, participants []string, eventDuration time.Duration) (matching []models.TimeSlotStartAndEnd, partialMatches []models.MatchingEventSlots, err error) {
	// get organizers timeslots

	// prepare timeslot stant and end time and participants for easy reconciliation
	organizerParticipant, participantsV2, err := ts.PrepareParticipantsDataForRecommendation(organizer, participants)
	if err != nil {
		log.Printf("error preparing participants data:: %s", err)
		return []models.TimeSlotStartAndEnd{}, []models.MatchingEventSlots{}, err
	}

	matchedSlots := []models.TimeSlotStartAndEnd{}
	partialMatchSlots := []models.MatchingEventSlots{}

	// get the common time slots
	for _, initiatorTimeSlot := range organizerParticipant.TimeSlots {
		availableCount := 0
		availableParticipants := []string{}
		unavailableParticipants := []string{}

		for _, participant := range participantsV2 {
			found := false
			for _, participantTimeSlot := range participant.TimeSlots {
				// check if the time slots overlap
				if utils.CheckIfTimeSlotsOverlap(initiatorTimeSlot, participantTimeSlot, eventDuration) {
					found = true
					break
				}
			}
			if found {
				availableCount++
				availableParticipants = append(availableParticipants, participant.Name)
			} else {
				unavailableParticipants = append(unavailableParticipants, participant.Name)
			}
		}

		if availableCount == len(participants) {
			matchedSlots = append(matchedSlots, initiatorTimeSlot)
		} else {
			partialMatchSlots = append(partialMatchSlots, models.MatchingEventSlots{
				Slot:                    initiatorTimeSlot,
				AvailableParticipants:   availableParticipants,
				UnavailableParticipants: unavailableParticipants,
			})
		}
	}

	return matchedSlots, partialMatchSlots, nil

}

func (ts *TimeslotServiceImplementaion) PrepareParticipantsDataForRecommendation(organizer string, participants []string) (models.Participant, []models.Participant, error) {
	// get the time slots and prepare participant for organizer and participants
	organizerParticipant, err := ts.GetUserTimeSlotsAndConvertToParticipant(organizer)
	if err != nil {
		log.Printf("error fetching organizer details:: %s", err)
		return models.Participant{}, []models.Participant{}, err
	}

	participantsV2 := []models.Participant{}
	for _, participant := range participants {
		p, err := ts.GetUserTimeSlotsAndConvertToParticipant(participant)
		if err != nil {
			log.Printf("error fetching participant details:: %s", err)
			return models.Participant{}, []models.Participant{}, err
		}

		participantsV2 = append(participantsV2, p)
	}

	return organizerParticipant, participantsV2, nil
}

func (ts *TimeslotServiceImplementaion) GetUserTimeSlotsAndConvertToParticipant(userName string) (models.Participant, error) {

	timeslotsOrganizer, err := ts.TimeslotRepo.GetTimeSlotsByUserName(userName)
	if err != nil {

		return models.Participant{}, err
	}

	initiator := models.Participant{
		Name: userName,
	}
	for _, timeSlot := range timeslotsOrganizer {
		startTime, endTime, valid := utils.ValidateAndFormatTimeStamp(timeSlot)
		if !valid {
			continue
		}
		tSAE := models.TimeSlotStartAndEnd{
			StartTime: startTime,
			EndTime:   endTime,
		}
		initiator.TimeSlots = append(initiator.TimeSlots, tSAE)
	}

	return initiator, nil
}

// ShowAccount godoc
// @Summary      Delete a time slot
// @Description  Delete time slot for a user by name
// @Tags         Timeslots
// @Accept       json
// @Produce      json
// @Param        username   path   string   true  "Timeslot request body"
// @Param        body       body    models.DeleteTimeSlotRequest   true  "Delete time slot request body"
// @Success      200  {object}  models.ServiceMessage
// @Failure      400  {object}  models.ServiceError
// @Failure      500  {object}  models.ServiceError
// @Router       /:username [delete]
func (ts *TimeslotServiceImplementaion) DeleteTimeSlotsByUserName(ctx *gin.Context) {
	userName := ctx.Param("username")
	if userName == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorHelper("Invalid request body", errors.New("username not provided")))
		return
	}

	var timeslot models.DeleteTimeSlotRequest

	if err := ctx.BindJSON(&timeslot); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorHelper("Invalid request body", err))
		return
	}

	verifyTimeSlotExists, err := ts.TimeslotRepo.GetTimeSlotsByUserName(userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("Error fetching time slots", err))
		return
	}

	if len(verifyTimeSlotExists) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "no time slots found for the user"})
		return
	}

	if !utils.SearchString(verifyTimeSlotExists, timeslot.Timeslot) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Time slot not found for the user"})
		return
	}

	fmt.Println(userName, timeslot.Timeslot)
	err = ts.TimeslotRepo.DeleteTimeSlotsByUserName(userName, timeslot.Timeslot)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorHelper("Error deleting time slots", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Time slots deleted successfully"})
}
