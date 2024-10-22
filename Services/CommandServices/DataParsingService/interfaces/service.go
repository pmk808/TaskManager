package interfaces

import "time"

type DataParsingService interface {
	ParseInt(value string) (int, error)
	ParseFloat(value string) (float64, error)
	ParseDate(value string) (time.Time, error)
}
