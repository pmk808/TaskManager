package ImportTaskService

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"taskmanager/Repository/CommandRepository/interfaces"
	"taskmanager/Services/CommandServices/ImportTaskService/schemas"
	validationInterfaces "taskmanager/Services/CommandServices/ImportTaskService/validation/interfaces"

	"github.com/sirupsen/logrus"
)

type importService struct {
	repo      interfaces.TaskCommandRepository
	validator validationInterfaces.Validator
	logger    *logrus.Logger
	directory string
}

func NewImportService(
	repo interfaces.TaskCommandRepository,
	validator validationInterfaces.Validator,
	logger *logrus.Logger,
	directory string,
) *importService {
	return &importService{
		repo:      repo,
		validator: validator,
		logger:    logger,
		directory: directory,
	}
}

func (s *importService) Import() (*schemas.TaskImportResponse, error) {
	stats := &schemas.TaskImportStats{
		StartTime: time.Now(),
	}
	defer func() {
		stats.EndTime = time.Now()
		stats.DurationMS = stats.EndTime.Sub(stats.StartTime).Milliseconds()
	}()

	entries, errs := s.readEntriesFromCSV()
	if len(errs) > 0 {
		return s.createErrorResponse(errs, stats), fmt.Errorf("failed to read entries from CSV")
	}

	stats.TotalProcessed = len(entries)
	s.logger.WithField("entry_count", len(entries)).Info("Entries read from CSV")

	if err := s.validator.ValidateBatch(entries); err != nil {
		s.logger.WithError(err).Error("Validation failed")
		return s.createErrorResponse([]error{err}, stats), fmt.Errorf("validation failed: %w", err)
	}

	s.logger.Info("All entries passed validation")

	if err := s.repo.BulkCreateTasks(entries); err != nil {
		s.logger.WithError(err).Error("Failed to import entries")
		return s.createErrorResponse([]error{err}, stats), fmt.Errorf("failed to import entries: %w", err)
	}

	stats.SuccessCount = len(entries)
	s.logger.WithFields(logrus.Fields{
		"duration":    stats.DurationMS,
		"entry_count": len(entries),
	}).Info("Data import process completed successfully")

	return &schemas.TaskImportResponse{
		Success:      true,
		Message:      "Import completed successfully",
		ImportedAt:   time.Now(),
		TotalEntries: len(entries),
		Stats:        stats,
	}, nil
}

func (s *importService) validateCSVHeaders(headers []string) error {
	expectedHeaders := []string{
		"Name",
		"Email",
		"Age",
		"Address",
		"Phone Number",
		"Department",
		"Position",
		"Salary",
		"Hire Date",
	}

	if len(headers) < len(expectedHeaders) {
		return fmt.Errorf("CSV is missing required columns. Expected: %v", expectedHeaders)
	}

	headers[0] = strings.TrimPrefix(headers[0], "\uFEFF")

	for i, expected := range expectedHeaders {
		if !strings.EqualFold(headers[i], expected) {
			return fmt.Errorf("invalid header at column %d. Expected: %s, Got: %s",
				i+1, expected, headers[i])
		}
	}

	return nil
}

func (s *importService) readEntriesFromCSV() ([]schemas.TaskImportEntry, []error) {
	var errors []error
	files, err := filepath.Glob(filepath.Join(s.directory, "*.csv"))
	if err != nil {
		return nil, []error{fmt.Errorf("failed to find CSV files: %w", err)}
	}

	if len(files) == 0 {
		return nil, []error{fmt.Errorf("no CSV files found in directory: %s", s.directory)}
	}

	filePath := files[0]
	s.logger.WithField("file", filePath).Info("Reading CSV file")

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read CSV file: %w", err)}
	}

	// Remove BOM if present
	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})

	// Create reader from cleaned content
	reader := csv.NewReader(bytes.NewReader(content))
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true // Allow lazy quotes

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read CSV rows: %w", err)}
	}

	if len(rows) == 0 {
		return nil, []error{fmt.Errorf("CSV file is empty")}
	}

	// Validate headers
	if err := s.validateCSVHeaders(rows[0]); err != nil {
		return nil, []error{fmt.Errorf("invalid CSV format: %w", err)}
	}

	s.logger.WithField("total_rows", len(rows)).Info("Total rows found in CSV")

	var entries []schemas.TaskImportEntry
	// Skip header row
	for i, row := range rows[1:] {
		s.logger.WithFields(logrus.Fields{
			"row_number": i + 2,
			"row_data":   row,
		}).Debug("Processing row")

		if len(row) < 9 {
			err := fmt.Errorf("row %d has insufficient columns (expected 9, got %d)", i+2, len(row))
			s.logger.WithError(err).Error("Invalid row")
			errors = append(errors, err)
			continue
		}

		entry, err := s.parseEntry(row, i+2)
		if err != nil {
			s.logger.WithError(err).WithField("row_number", i+2).Error("Failed to parse row")
			errors = append(errors, err)
			continue
		}
		entries = append(entries, entry)
	}

	if len(errors) > 0 {
		s.logger.WithField("error_count", len(errors)).Error("Encountered errors while reading CSV")
		return entries, errors
	}

	s.logger.WithField("entry_count", len(entries)).Info("Finished reading entries from CSV")
	return entries, nil
}

func (s *importService) parseEntry(row []string, rowNum int) (schemas.TaskImportEntry, error) {
	age, err := strconv.Atoi(row[2])
	if err != nil {
		return schemas.TaskImportEntry{}, fmt.Errorf("invalid age at row %d: %w", rowNum, err)
	}

	salary, err := strconv.ParseFloat(row[7], 64)
	if err != nil {
		return schemas.TaskImportEntry{}, fmt.Errorf("invalid salary at row %d: %w", rowNum, err)
	}

	var hireDate time.Time
	dateFormats := []string{
		"02/01/2006", // DD/MM/YYYY
		"2006-01-02", // YYYY-MM-DD
		"01/02/2006", // MM/DD/YYYY
		"2006/01/02", // YYYY/MM/DD
	}

	var parseErr error
	for _, format := range dateFormats {
		hireDate, parseErr = time.Parse(format, row[8])
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return schemas.TaskImportEntry{}, fmt.Errorf("invalid hire date at row %d: date must be in DD/MM/YYYY format", rowNum)
	}

	return schemas.TaskImportEntry{
		Name:        strings.TrimSpace(row[0]),
		Email:       strings.TrimSpace(row[1]),
		Age:         age,
		Address:     strings.TrimSpace(row[3]),
		PhoneNumber: strings.TrimSpace(row[4]),
		Department:  strings.TrimSpace(row[5]),
		Position:    strings.TrimSpace(row[6]),
		Salary:      salary,
		HireDate:    hireDate,
	}, nil
}

func (s *importService) createErrorResponse(errs []error, stats *schemas.TaskImportStats) *schemas.TaskImportResponse {
	errorMessages := make([]string, len(errs))
	for i, err := range errs {
		errorMessages[i] = err.Error()
	}

	stats.ErrorCount = len(errs)

	return &schemas.TaskImportResponse{
		Success:      false,
		Message:      "Import failed",
		ImportedAt:   time.Now(),
		TotalEntries: stats.TotalProcessed,
		Errors:       errorMessages,
		Stats:        stats,
	}
}
