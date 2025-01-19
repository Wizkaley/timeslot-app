package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"timeslot-app/models"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTimeslotRepo struct {
	mock.Mock
}

func (m *MockTimeslotRepo) Create(timeSlots []models.TimeSlot) error {
	args := m.Called(timeSlots)
	return args.Error(0)
}

func (m *MockTimeslotRepo) GetTimeSlotsByUserName(userName string) ([]string, error) {
	args := m.Called(userName)
	return args.Get(0).([]string), args.Error(1)
}

func TestCreateTimeSlot(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockTimeslotRepo := new(MockTimeslotRepo)
		mockUserRepo := new(MockUserRepo)
		timeslotService := &TimeslotServiceImplementaion{
			TimeslotRepo: mockTimeslotRepo,
			UserRepo:     mockUserRepo,
		}

		router := gin.Default()
		router.POST("/timeslot", timeslotService.CreateTimeSlot)

		userID, _ := uuid.NewV4()
		mockUser := models.User{ID: userID, Name: "John Doe"}
		mockUserRepo.On("Get", "John Doe").Return(mockUser, nil)

		timeSlot := "02 Jan 2025 2-4 PM MST"
		userTimeSlotReq := models.UserTimeSlotRequest{
			UserName:  "John Doe",
			TimeSlots: []string{timeSlot},
		}
		userTimeSlotReqJSON, _ := json.Marshal(userTimeSlotReq)

		req, _ := http.NewRequest(http.MethodPost, "/timeslot", bytes.NewBuffer(userTimeSlotReqJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		mockTimeslotRepo.On("Create", mock.Anything).Return(nil)

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		mockTimeslotRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockTimeslotRepo := new(MockTimeslotRepo)
		mockUserRepo := new(MockUserRepo)
		timeslotService := &TimeslotServiceImplementaion{
			TimeslotRepo: mockTimeslotRepo,
			UserRepo:     mockUserRepo,
		}

		router := gin.Default()
		router.POST("/timeslot", timeslotService.CreateTimeSlot)

		req, _ := http.NewRequest(http.MethodPost, "/timeslot", bytes.NewBuffer([]byte("{invalid json}")))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockTimeslotRepo := new(MockTimeslotRepo)
		mockUserRepo := new(MockUserRepo)
		timeslotService := &TimeslotServiceImplementaion{
			TimeslotRepo: mockTimeslotRepo,
			UserRepo:     mockUserRepo,
		}

		router := gin.Default()
		router.POST("/timeslot", timeslotService.CreateTimeSlot)

		userTimeSlotReq := models.UserTimeSlotRequest{
			UserName:  "John Doe",
			TimeSlots: []string{"2023-10-10T10:00:00Z"},
		}
		userTimeSlotReqJSON, _ := json.Marshal(userTimeSlotReq)

		req, _ := http.NewRequest(http.MethodPost, "/timeslot", bytes.NewBuffer(userTimeSlotReqJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		mockUserRepo.On("Get", "John Doe").Return(models.User{}, errors.New("user not found"))

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Invalid Time Slot Format", func(t *testing.T) {
		mockTimeslotRepo := new(MockTimeslotRepo)
		mockUserRepo := new(MockUserRepo)
		timeslotService := &TimeslotServiceImplementaion{
			TimeslotRepo: mockTimeslotRepo,
			UserRepo:     mockUserRepo,
		}

		router := gin.Default()
		router.POST("/timeslot", timeslotService.CreateTimeSlot)

		userID, _ := uuid.NewV4()
		mockUser := models.User{ID: userID, Name: "John Doe"}
		mockUserRepo.On("Get", "John Doe").Return(mockUser, nil)

		userTimeSlotReq := models.UserTimeSlotRequest{
			UserName:  "John Doe",
			TimeSlots: []string{"invalid-time-slot"},
		}
		userTimeSlotReqJSON, _ := json.Marshal(userTimeSlotReq)

		req, _ := http.NewRequest(http.MethodPost, "/timeslot", bytes.NewBuffer(userTimeSlotReqJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error Creating Time Slots", func(t *testing.T) {
		mockTimeslotRepo := new(MockTimeslotRepo)
		mockUserRepo := new(MockUserRepo)
		timeslotService := &TimeslotServiceImplementaion{
			TimeslotRepo: mockTimeslotRepo,
			UserRepo:     mockUserRepo,
		}

		router := gin.Default()
		router.POST("/timeslot", timeslotService.CreateTimeSlot)

		userID, _ := uuid.NewV4()
		mockUser := models.User{ID: userID, Name: "John Doe"}
		mockUserRepo.On("Get", "John Doe").Return(mockUser, nil)

		timeSlot := "02 Jan 2025 4-6 PM MST"
		userTimeSlotReq := models.UserTimeSlotRequest{
			UserName:  "John Doe",
			TimeSlots: []string{timeSlot},
		}
		userTimeSlotReqJSON, _ := json.Marshal(userTimeSlotReq)

		req, _ := http.NewRequest(http.MethodPost, "/timeslot", bytes.NewBuffer(userTimeSlotReqJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		mockTimeslotRepo.On("Create", mock.Anything).Return(errors.New("error creating time slots"))

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		mockTimeslotRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetTimeSlotsByUserName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockTimeslotRepo := new(MockTimeslotRepo)
		timeslotService := &TimeslotServiceImplementaion{
			TimeslotRepo: mockTimeslotRepo,
		}

		router := gin.Default()
		router.GET("/timeslot/:username", timeslotService.GetTimeSlotsByUserName)

		timeSlots := []string{"2023-10-10T10:00:00Z", "2023-10-11T10:00:00Z"}
		mockTimeslotRepo.On("GetTimeSlotsByUserName", "John Doe").Return(timeSlots, nil)

		req, _ := http.NewRequest(http.MethodGet, "/timeslot/John Doe", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		expectedResponse := models.TimeSlotResponse{
			UserName:  "John Doe",
			TimeSlots: timeSlots,
		}
		var actualResponse models.TimeSlotResponse
		json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
		assert.Equal(t, expectedResponse, actualResponse)
		mockTimeslotRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockTimeslotRepo := new(MockTimeslotRepo)
		timeslotService := &TimeslotServiceImplementaion{
			TimeslotRepo: mockTimeslotRepo,
		}

		router := gin.Default()
		router.GET("/timeslot/:username", timeslotService.GetTimeSlotsByUserName)

		mockTimeslotRepo.On("GetTimeSlotsByUserName", "John Doe").Return(nil, errors.New("user not found"))

		req, _ := http.NewRequest(http.MethodGet, "/timeslot/John Doe", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		mockTimeslotRepo.AssertExpectations(t)
	})
}

// type MockTimeslotRepo struct {
// 	mock.Mock
// }

func (m *MockTimeslotRepo) RecommendSlots(req models.RecommendSlotsRequest) ([]models.TimeSlot, error) {
	args := m.Called(req)
	return args.Get(0).([]models.TimeSlot), args.Error(1)
}

func TestRecommendSlots(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockTimeslotRepo)
	timeslotService := NewTimeslotService(mockRepo)

	router := gin.Default()
	router.POST("/recommend-slots", timeslotService.RecommendSlots)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.RecommendSlotsRequest{
			UserID: "test-user-id",
			// Add other fields as necessary
		}
		reqJSON, _ := json.Marshal(reqBody)

		expectedSlots := []models.TimeSlot{
			{ID: "slot1", Time: "10:00 AM"},
			{ID: "slot2", Time: "11:00 AM"},
		}

		mockRepo.On("RecommendSlots", reqBody).Return(expectedSlots, nil)

		req, _ := http.NewRequest(http.MethodPost, "/recommend-slots", bytes.NewBuffer(reqJSON))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var response map[string][]models.TimeSlot
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedSlots, response["slots"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidRequestBody", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/recommend-slots", bytes.NewBuffer([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		reqBody := models.RecommendSlotsRequest{
			UserID: "test-user-id",
			// Add other fields as necessary
		}
		reqJSON, _ := json.Marshal(reqBody)

		mockRepo.On("RecommendSlots", reqBody).Return(nil, assert.AnError)

		req, _ := http.NewRequest(http.MethodPost, "/recommend-slots", bytes.NewBuffer(reqJSON))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockRepo.AssertExpectations(t)
	})
}
