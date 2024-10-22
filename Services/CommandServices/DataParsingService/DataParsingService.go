package services

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type DataParsingService struct {
	logger *logrus.Logger
}

func NewDataParsingService(logger *logrus.Logger) *DataParsingService {
	return &DataParsingService{
		logger: logger,
	}
}

func (s *DataParsingService) ParseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("01-02-06", dateStr)
	if err != nil {
		s.logger.WithField("value", dateStr).Error("Failed to parse date")
		return time.Time{}, err
	}
	return date, nil
}

func (s *DataParsingService) ParseInt(intStr string) (int, error) {
	value, err := strconv.Atoi(intStr)
	if err != nil {
		s.logger.WithField("value", intStr).Error("Failed to parse int")
		return 0, err
	}
	return value, nil
}

func (s *DataParsingService) ParseFloat(floatStr string) (float64, error) {
	value, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		s.logger.WithField("value", floatStr).Error("Failed to parse float")
		return 0.0, err
	}
	return value, nil
}
