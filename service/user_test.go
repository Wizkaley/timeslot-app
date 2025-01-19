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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) UserExists(userName string) (bool, error) {
	args := m.Called(userName)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) Get(userName string) (models.User, error) {
	args := m.Called(userName)
	return args.Get(0).(models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		userService := &UserService{userRepo: mockUserRepo}

		router := gin.Default()
		router.POST("/user", userService.CreateUser)

		userReq := models.UserCreateRequest{Name: "John Doe"}
		userReqJSON, _ := json.Marshal(userReq)

		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userReqJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		mockUserRepo.On("Create", mock.Anything).Return(nil)

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		userService := &UserService{userRepo: mockUserRepo}

		router := gin.Default()
		router.POST("/user", userService.CreateUser)

		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte("{invalid json}")))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Error Creating User", func(t *testing.T) {
		mockUserRepo := new(MockUserRepo)
		userService := &UserService{userRepo: mockUserRepo}

		router := gin.Default()
		router.POST("/user", userService.CreateUser)

		userReq := models.UserCreateRequest{Name: "John Doe"}
		userReqJSON, _ := json.Marshal(userReq)

		req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(userReqJSON))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		mockUserRepo.On("Create", mock.Anything).Return(errors.New("Error creating user"))

		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		mockUserRepo.AssertExpectations(t)
	})
}
