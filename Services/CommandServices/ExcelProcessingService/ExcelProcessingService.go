package services

import (
	"fmt"
	"taskmanager/Services/CommandServices/DataParsingService/interfaces"
	"taskmanager/schemas"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type ExcelProcessingService struct {
	logger         *logrus.Logger
	dataParsingSvc interfaces.DataParsingService
}

func NewExcelProcessingService(logger *logrus.Logger, dataParsingSvc interfaces.DataParsingService) *ExcelProcessingService {
	return &ExcelProcessingService{
		logger:         logger,
		dataParsingSvc: dataParsingSvc,
	}
}

func (s *ExcelProcessingService) ReadTasksFromSpreadsheet(filePath string) ([]schemas.Task, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	var tasks []schemas.Task
	for _, row := range rows {
		age, err := s.dataParsingSvc.ParseInt(row[2])
		if err != nil {
			s.logger.WithError(err).WithField("row", row).Error("Failed to parse age")
			continue
		}

		salary, err := s.dataParsingSvc.ParseFloat(row[7])
		if err != nil {
			s.logger.WithError(err).WithField("row", row).Error("Failed to parse salary")
			continue
		}

		hireDate, err := s.dataParsingSvc.ParseDate(row[8])
		if err != nil {
			s.logger.WithError(err).WithField("row", row).Error("Failed to parse hire date")
			continue
		}

		task := schemas.Task{
			Name:        row[0],
			Email:       row[1],
			Age:         age,
			Address:     row[3],
			PhoneNumber: row[4],
			Department:  row[5],
			Position:    row[6],
			Salary:      salary,
			HireDate:    hireDate,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
