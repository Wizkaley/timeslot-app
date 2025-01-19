package repository

import (
	"testing"
	"timeslot-app/models"

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

func TestTimeslotRepo(t *testing.T) {
	mockRepo := new(MockTimeslotRepo)
	tID, _ := uuid.NewV4()
	ts := "02 Jan 2025 2-4 PM MST"
	timeSlot := models.TimeSlot{ID: tID, TimeSlot: ts}

	t.Run("Create", func(t *testing.T) {
		mockRepo.On("Create", []models.TimeSlot{timeSlot}).Return(nil)

		err := mockRepo.Create([]models.TimeSlot{timeSlot})
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetTimeSlotsByUserName", func(t *testing.T) {
		timeSlots := []string{"10:00 AM"}
		mockRepo.On("GetTimeSlotsByUserName", "testuser").Return(timeSlots, nil)

		result, err := mockRepo.GetTimeSlotsByUserName("testuser")
		assert.NoError(t, err)
		assert.Equal(t, timeSlots, result)
		mockRepo.AssertExpectations(t)
	})
}
