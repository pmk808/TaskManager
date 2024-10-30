package TaskQueryService

import (
	"context"
	"testing"
	"time"

	repoInterfaces "taskmanager/Repository/QueryRepository/interfaces"
	serviceInterfaces "taskmanager/Services/QueryServices/interfaces"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskQueryRepository is a mock implementation of TaskQueryRepository
type MockTaskQueryRepository struct {
	mock.Mock
}

func (m *MockTaskQueryRepository) GetActiveTasks(ctx context.Context, clientName string, clientID string) ([]repoInterfaces.TaskDTO, error) {
	args := m.Called(ctx, clientName, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repoInterfaces.TaskDTO), args.Error(1)
}

func (m *MockTaskQueryRepository) GetTaskStatusHistory(ctx context.Context, clientName string, clientID string) ([]repoInterfaces.TaskStatusDTO, error) {
	args := m.Called(ctx, clientName, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repoInterfaces.TaskStatusDTO), args.Error(1)
}

func TestGetActiveTasks(t *testing.T) {
	// Setup
	logger := logrus.New()
	mockRepo := new(MockTaskQueryRepository)
	service := NewTaskQueryService(mockRepo, logger)
	ctx := context.Background()

	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name       string
		clientName string
		clientID   string
		mockSetup  func()
		verify     func(*testing.T, *serviceInterfaces.TasksResponseDTO, error)
	}{
		{
			name:       "Valid request with data",
			clientName: "Test Client",
			clientID:   validUUID,
			mockSetup: func() {
				mockRepo.On("GetActiveTasks", ctx, "Test Client", validUUID).Return([]repoInterfaces.TaskDTO{
					{
						ID:         1,
						Name:       "John Doe",
						Email:      "john@example.com",
						Department: "IT",
						Position:   "Developer",
						IsActive:   true,
						ClientName: "Test Client",
						ClientID:   validUUID,
					},
				}, nil).Once()
			},
			verify: func(t *testing.T, response *serviceInterfaces.TasksResponseDTO, err error) {
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.Equal(t, 1, response.TotalCount)
				assert.Len(t, response.Tasks, 1)
				assert.Equal(t, "John Doe", response.Tasks[0].Name)
			},
		},
		{
			name:       "Empty client name",
			clientName: "",
			clientID:   validUUID,
			mockSetup:  func() {},
			verify: func(t *testing.T, response *serviceInterfaces.TasksResponseDTO, err error) {
				assert.NoError(t, err)
				assert.False(t, response.Success)
				assert.Contains(t, response.Message, "client name cannot be empty")
			},
		},
		{
			name:       "Invalid UUID",
			clientName: "Test Client",
			clientID:   "invalid-uuid",
			mockSetup:  func() {},
			verify: func(t *testing.T, response *serviceInterfaces.TasksResponseDTO, err error) {
				assert.NoError(t, err)
				assert.False(t, response.Success)
				assert.Contains(t, response.Message, "invalid client ID format")
			},
		},
		{
			name:       "Valid request with no data",
			clientName: "Test Client",
			clientID:   validUUID,
			mockSetup: func() {
				mockRepo.On("GetActiveTasks", ctx, "Test Client", validUUID).Return([]repoInterfaces.TaskDTO{}, nil).Once()
			},
			verify: func(t *testing.T, response *serviceInterfaces.TasksResponseDTO, err error) {
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.Equal(t, 0, response.TotalCount)
				assert.Empty(t, response.Tasks)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous mock calls
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.mockSetup()
			response, err := service.GetActiveTasks(ctx, tt.clientName, tt.clientID)
			tt.verify(t, response, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTaskStatusHistory(t *testing.T) {
	// Setup
	logger := logrus.New()
	mockRepo := new(MockTaskQueryRepository)
	service := NewTaskQueryService(mockRepo, logger)
	ctx := context.Background()

	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	now := time.Now()

	tests := []struct {
		name       string
		clientName string
		clientID   string
		mockSetup  func()
		verify     func(*testing.T, *serviceInterfaces.StatusHistoryResponseDTO, error)
	}{
		{
			name:       "Valid request with history",
			clientName: "Test Client",
			clientID:   validUUID,
			mockSetup: func() {
				mockRepo.On("GetTaskStatusHistory", ctx, "Test Client", validUUID).Return([]repoInterfaces.TaskStatusDTO{
					{
						TaskID:            1,
						Status:            "IN_PROGRESS",
						StatusDescription: "Task started",
						UpdatedBy:         "system",
						CreatedAt:         now,
					},
				}, nil).Once()
			},
			verify: func(t *testing.T, response *serviceInterfaces.StatusHistoryResponseDTO, err error) {
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.Equal(t, 1, response.TotalCount)
				assert.Len(t, response.History, 1)
				assert.Equal(t, "IN_PROGRESS", response.History[0].Status)
			},
		},
		{
			name:       "Empty client name",
			clientName: "",
			clientID:   validUUID,
			mockSetup:  func() {},
			verify: func(t *testing.T, response *serviceInterfaces.StatusHistoryResponseDTO, err error) {
				assert.NoError(t, err)
				assert.False(t, response.Success)
				assert.Contains(t, response.Message, "client name cannot be empty")
			},
		},
		{
			name:       "Invalid UUID",
			clientName: "Test Client",
			clientID:   "invalid-uuid",
			mockSetup:  func() {},
			verify: func(t *testing.T, response *serviceInterfaces.StatusHistoryResponseDTO, err error) {
				assert.NoError(t, err)
				assert.False(t, response.Success)
				assert.Contains(t, response.Message, "invalid client ID format")
			},
		},
		{
			name:       "Valid request with no history",
			clientName: "Test Client",
			clientID:   validUUID,
			mockSetup: func() {
				mockRepo.On("GetTaskStatusHistory", ctx, "Test Client", validUUID).Return([]repoInterfaces.TaskStatusDTO{}, nil).Once()
			},
			verify: func(t *testing.T, response *serviceInterfaces.StatusHistoryResponseDTO, err error) {
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.Equal(t, 0, response.TotalCount)
				assert.Empty(t, response.History)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous mock calls
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.mockSetup()
			response, err := service.GetTaskStatusHistory(ctx, tt.clientName, tt.clientID)
			tt.verify(t, response, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
