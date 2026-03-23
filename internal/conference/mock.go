package conference

import (
	"context"
	"fmt"
	"time"
)

type Service interface {
	CreateLink(ctx context.Context, bookingID string) (string, error)
}

type MockService struct {
	BaseURL string
	Timeout time.Duration
}

func NewMockService(baseURL string, timeout time.Duration) *MockService {
	return &MockService{
		BaseURL: baseURL,
		Timeout: timeout,
	}
}

func (s *MockService) CreateLink(ctx context.Context, bookingID string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(minDuration(s.Timeout, 50*time.Millisecond)):
		return fmt.Sprintf("%s/%s", s.BaseURL, bookingID), nil
	}
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
