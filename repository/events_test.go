package repository

import (
	"testing"
	"timeslot-app/models"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventRepo struct {
	mock.Mock
}

func (m *MockEventRepo) CreateEvent(event models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventRepo) GetEvent(eventID string) (models.Event, error) {
	args := m.Called(eventID)
	return args.Get(0).(models.Event), args.Error(1)
}

func (m *MockEventRepo) DeleteEvent(eventID string) error {
	args := m.Called(eventID)
	return args.Error(0)
}

func (m *MockEventRepo) GetEventsForUser(username string) ([]models.Event, error) {
	args := m.Called(username)
	return args.Get(0).([]models.Event), args.Error(1)
}

func TestEventRepo(t *testing.T) {
	mockRepo := new(MockEventRepo)
	eventID, _ := uuid.NewV4()
	event := models.Event{ID: eventID, Title: "Test Event"}

	t.Run("CreateEvent", func(t *testing.T) {
		mockRepo.On("CreateEvent", event).Return(nil)

		err := mockRepo.CreateEvent(event)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetEvent", func(t *testing.T) {
		mockRepo.On("GetEvent", "1").Return(event, nil)

		result, err := mockRepo.GetEvent("1")
		assert.NoError(t, err)
		assert.Equal(t, event, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		mockRepo.On("DeleteEvent", "1").Return(nil)

		err := mockRepo.DeleteEvent("1")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetEventsForUser", func(t *testing.T) {
		events := []models.Event{event}
		mockRepo.On("GetEventsForUser", "testuser").Return(events, nil)

		result, err := mockRepo.GetEventsForUser("testuser")
		assert.NoError(t, err)
		assert.Equal(t, events, result)
		mockRepo.AssertExpectations(t)
	})
}
