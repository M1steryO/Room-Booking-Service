//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

type tokenResponse struct {
	Token string `json:"token"`
}

type roomResponse struct {
	Room struct {
		ID string `json:"id"`
	} `json:"room"`
}

type slotsResponse struct {
	Slots []struct {
		ID string `json:"id"`
	} `json:"slots"`
}

type bookingResponse struct {
	Booking struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"booking"`
}

func TestCreateRoomScheduleAndBooking(t *testing.T) {
	baseURL := mustBaseURL(t)

	adminToken := dummyLogin(t, baseURL, "admin")
	userToken := dummyLogin(t, baseURL, "user")

	roomID := createRoom(t, baseURL, adminToken)
	createSchedule(t, baseURL, adminToken, roomID)

	date := nextWeekday(time.Now().UTC(), time.Monday).Format("2006-01-02")
	slotID := firstSlot(t, baseURL, userToken, roomID, date)
	if slotID == "" {
		t.Fatal("expected non-empty slot id")
	}

	bookingID := createBooking(t, baseURL, userToken, slotID)
	if bookingID == "" {
		t.Fatal("expected non-empty booking id")
	}
}

func TestCancelBooking(t *testing.T) {
	baseURL := mustBaseURL(t)

	adminToken := dummyLogin(t, baseURL, "admin")
	userToken := dummyLogin(t, baseURL, "user")

	roomID := createRoom(t, baseURL, adminToken)
	createSchedule(t, baseURL, adminToken, roomID)

	date := nextWeekday(time.Now().UTC(), time.Tuesday).Format("2006-01-02")
	slotID := firstSlot(t, baseURL, userToken, roomID, date)
	bookingID := createBooking(t, baseURL, userToken, slotID)

	body, _ := json.Marshal(map[string]any{})
	request, _ := http.NewRequest(http.MethodPost, baseURL+"/bookings/"+bookingID+"/cancel", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+userToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("cancel booking: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	var payload bookingResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode cancel response: %v", err)
	}

	if payload.Booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status, got %s", payload.Booking.Status)
	}
}

func mustBaseURL(t *testing.T) string {
	t.Helper()

	value := os.Getenv("E2E_BASE_URL")
	if value == "" {
		t.Skip("E2E_BASE_URL is required")
	}
	return value
}

func dummyLogin(t *testing.T, baseURL string, role string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"role": role})
	response, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("dummy login: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from dummyLogin, got %d", response.StatusCode)
	}

	var payload tokenResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode dummy login: %v", err)
	}

	return payload.Token
}

func createRoom(t *testing.T, baseURL string, token string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"name":        "Room " + time.Now().UTC().Format("150405.000000"),
		"description": "e2e",
		"capacity":    4,
	})

	request, _ := http.NewRequest(http.MethodPost, baseURL+"/rooms/create", bytes.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("create room: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}

	var payload roomResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode room response: %v", err)
	}

	return payload.Room.ID
}

func createSchedule(t *testing.T, baseURL string, token string, roomID string) {
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
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}
}

func firstSlot(t *testing.T, baseURL string, token string, roomID string, date string) string {
	t.Helper()

	request, _ := http.NewRequest(http.MethodGet, baseURL+"/rooms/"+roomID+"/slots/list?date="+date, nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("list slots: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}

	var payload slotsResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode slots response: %v", err)
	}

	if len(payload.Slots) == 0 {
		t.Fatal("expected non-empty slots list")
	}

	return payload.Slots[0].ID
}

func createBooking(t *testing.T, baseURL string, token string, slotID string) string {
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
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}

	var payload bookingResponse
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
