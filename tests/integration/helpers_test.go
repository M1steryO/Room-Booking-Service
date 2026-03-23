package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

func mustIntegrationBaseURL(t *testing.T) string {
	t.Helper()

	if value := os.Getenv("INTEGRATION_BASE_URL"); value != "" {
		return value
	}
	if value := os.Getenv("E2E_BASE_URL"); value != "" {
		return value
	}

	t.Skip("INTEGRATION_BASE_URL or E2E_BASE_URL is required")
	return ""
}

func mustDummyLogin(t *testing.T, baseURL string, role string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"role": role})
	response, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("dummy login: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from dummyLogin, got %d", response.StatusCode)
	}

	var payload tokenResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode dummy login: %v", err)
	}

	return payload.Token
}

func mustCreateRoom(t *testing.T, baseURL string, token string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"name":        "Room " + time.Now().UTC().Format("150405.000000"),
		"description": "integration",
		"capacity":    6,
	})

	request, _ := http.NewRequest(http.MethodPost, baseURL+"/rooms/create", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("create room: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}

	var payload roomEnvelope
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode room response: %v", err)
	}

	return payload.Room.ID
}

func mustCreateSchedule(t *testing.T, baseURL string, token string, roomID string) {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"daysOfWeek": []int{1, 2, 3, 4, 5},
		"startTime":  "09:00",
		"endTime":    "18:00",
	})

	request, _ := http.NewRequest(http.MethodPost, baseURL+"/rooms/"+roomID+"/schedule/create", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}
}

func mustFirstSlot(t *testing.T, baseURL string, token string, roomID string, date string) string {
	t.Helper()

	request, _ := http.NewRequest(http.MethodGet, baseURL+"/rooms/"+roomID+"/slots/list?date="+date, nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("list slots: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	var payload slotsEnvelope
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode slots response: %v", err)
	}

	if len(payload.Slots) == 0 {
		t.Fatal("expected non-empty slots list")
	}

	return payload.Slots[0].ID
}

func mustCreateBooking(t *testing.T, baseURL string, token string, slotID string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"slotId":               slotID,
		"createConferenceLink": true,
	})

	request, _ := http.NewRequest(http.MethodPost, baseURL+"/bookings/create", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}

	var payload bookingEnvelope
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode booking response: %v", err)
	}

	return payload.Booking.ID
}

func nextWeekday(from time.Time, weekday time.Weekday) time.Time {
	date := from.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	for {
		if date.Weekday() == weekday {
			return date
		}
		date = date.AddDate(0, 0, 1)
	}
}
