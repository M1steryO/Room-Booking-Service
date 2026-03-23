package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

type tokenResponse struct {
	Token string `json:"token"`
}

type roomEnvelope struct {
	Room struct {
		ID string `json:"id"`
	} `json:"room"`
}

type slotsEnvelope struct {
	Slots []struct {
		ID string `json:"id"`
	} `json:"slots"`
}

type bookingEnvelope struct {
	Booking struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"booking"`
}

func TestRoomScheduleBookingFlow(t *testing.T) {
	baseURL := mustIntegrationBaseURL(t)

	adminToken := mustDummyLogin(t, baseURL, "admin")
	userToken := mustDummyLogin(t, baseURL, "user")

	roomID := mustCreateRoom(t, baseURL, adminToken)
	mustCreateSchedule(t, baseURL, adminToken, roomID)

	date := nextWeekday(time.Now().UTC(), time.Wednesday).Format("2006-01-02")
	slotID := mustFirstSlot(t, baseURL, userToken, roomID, date)

	bookingID := mustCreateBooking(t, baseURL, userToken, slotID)
	if bookingID == "" {
		t.Fatal("expected non-empty booking id")
	}
}

func TestCancelBookingFlow(t *testing.T) {
	baseURL := mustIntegrationBaseURL(t)

	adminToken := mustDummyLogin(t, baseURL, "admin")
	userToken := mustDummyLogin(t, baseURL, "user")

	roomID := mustCreateRoom(t, baseURL, adminToken)
	mustCreateSchedule(t, baseURL, adminToken, roomID)

	date := nextWeekday(time.Now().UTC(), time.Thursday).Format("2006-01-02")
	slotID := mustFirstSlot(t, baseURL, userToken, roomID, date)
	bookingID := mustCreateBooking(t, baseURL, userToken, slotID)

	request, _ := http.NewRequest(http.MethodPost, baseURL+"/bookings/"+bookingID+"/cancel", bytes.NewReader([]byte(`{}`)))
	request.Header.Set("Authorization", "Bearer "+userToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("cancel booking: %v", err)
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Fatalf("close response body: %v", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	var payload bookingEnvelope
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode cancel response: %v", err)
	}

	if payload.Booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status, got %s", payload.Booking.Status)
	}
}
