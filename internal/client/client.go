package client

import "context"

type ConferenceService interface {
	CreateLink(ctx context.Context, bookingID string) (string, error)
}
