package conference

import (
	"context"
	"fmt"
	"github.com/M1steryO/Room-Booking-Service/internal/client"
	"time"
)

type serv struct {
	BaseURL string
	Timeout time.Duration
}

func NewConferenceService(baseURL string, timeout time.Duration) client.ConferenceService {
	return &serv{
		BaseURL: baseURL,
		Timeout: timeout,
	}
}

func (s *serv) CreateLink(ctx context.Context, bookingID string) (string, error) {
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
