package timeservice

import (
	"time"

	timeservice "github.com/mohit/finance-go/internal/common/time_service_interface"
)

// Service provides the current UTC time.
type Service struct{}

// Ensure Service implements timeservice.IService.
var _ timeservice.IService = &Service{}

// New creates a new instance of the Service.
func New() *Service {
	return &Service{}
}

// NowUTC returns the current time in UTC.
func (t *Service) NowUTC() time.Time {
	return time.Now().UTC()
}
