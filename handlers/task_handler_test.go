package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "taskmanager/interfaces"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) ImportData() error {
	args := m.Called()
	return args.Error(0)
}

func TestImportData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := logrus.New()
	logger.SetOutput(io.Discard) // Discard logs in tests

	tests := []struct {
		name           string
		serviceError   error
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Successful import",
			serviceError:   nil,
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Data import completed successfully"},
		},
		{
			name:           "Import failure",
			serviceError:   errors.New("import failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Failed to import data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTaskService)
			mockService.On("ImportData").Return(tt.serviceError)

			handler := NewTaskHandler(mockService, logger)

			router := gin.New()
			router.POST("/import", handler.ImportData)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/import", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}
