package validation

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
)

type QueryValidator struct{}

func NewQueryValidator() *QueryValidator {
	return &QueryValidator{}
}

func (v *QueryValidator) ValidateClientParams(clientName, clientID string) error {
	if strings.TrimSpace(clientName) == "" {
		return fmt.Errorf("client name cannot be empty")
	}

	if strings.TrimSpace(clientID) == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	// Validate UUID format
	if _, err := uuid.Parse(clientID); err != nil {
		return fmt.Errorf("invalid client ID format: must be a valid UUID")
	}

	return nil
}