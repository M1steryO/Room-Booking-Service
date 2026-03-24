package bookings_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	repmocks "github.com/avito-internships/test-backend-1-M1steryO/internal/repository/mocks"
	bookingsuc "github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/bookings"
)

type fixedClock struct{ now time.Time }

func (f fixedClock) Now() time.Time { return f.now }

type conferenceStub struct {
	createFn func(ctx context.Context, bookingID string) (string, error)
}

func (s conferenceStub) CreateLink(ctx context.Context, bookingID string) (string, error) {
	return s.createFn(ctx, bookingID)
}

func TestCreateForbiddenForAdmin(t *testing.T) {
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListAllMock.Optional()
	bookingsRepo.ListByUserMock.Optional()
	bookingsRepo.CancelByOwnerMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", nil }},
		fixedClock{now: time.Now().UTC()},
	)

	_, err := uc.Create(context.Background(), "u1", domain.RoleAdmin, "s1", false)
	if err == nil || domain.AsAppError(err).Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %v", err)
	}
}

func TestCreatePastSlotRejected(t *testing.T) {
	now := time.Date(2026, 3, 23, 10, 0, 0, 0, time.UTC)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Slot, error) {
		return domain.Slot{ID: "s1", Start: now.Add(-time.Hour)}, nil
	})
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListAllMock.Optional()
	bookingsRepo.ListByUserMock.Optional()
	bookingsRepo.CancelByOwnerMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", nil }},
		fixedClock{now: now},
	)

	_, err := uc.Create(context.Background(), "u1", domain.RoleUser, "s1", false)
	if err == nil || domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %v", err)
	}
}

func TestCancelDelegates(t *testing.T) {
	var called bool
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.CancelByOwnerMock.Set(func(_ context.Context, bookingID, userID string) (domain.Booking, error) {
		called = bookingID == "b1" && userID == "u1"
		return domain.Booking{ID: bookingID, UserID: userID}, nil
	})
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListAllMock.Optional()
	bookingsRepo.ListByUserMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", errors.New("unused") }},
		fixedClock{now: time.Now().UTC()},
	)

	_, err := uc.Cancel(context.Background(), "u1", domain.RoleUser, "b1")
	if err != nil {
		t.Fatalf("Cancel error: %v", err)
	}
	if !called {
		t.Fatal("expected CancelByOwner to be called")
	}
}

func TestListAllForbiddenForUser(t *testing.T) {
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListAllMock.Optional()
	bookingsRepo.ListByUserMock.Optional()
	bookingsRepo.CancelByOwnerMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", nil }},
		fixedClock{now: time.Now().UTC()},
	)

	_, _, err := uc.ListAll(context.Background(), domain.RoleUser, 1, 20)
	if err == nil || domain.AsAppError(err).Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %v", err)
	}
}

func TestListAllInvalidPagination(t *testing.T) {
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListAllMock.Optional()
	bookingsRepo.ListByUserMock.Optional()
	bookingsRepo.CancelByOwnerMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", nil }},
		fixedClock{now: time.Now().UTC()},
	)

	_, _, err := uc.ListAll(context.Background(), domain.RoleAdmin, 0, 20)
	if err == nil || domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %v", err)
	}
}

func TestListAllSuccess(t *testing.T) {
	expected := []domain.Booking{{ID: "b1"}}
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.ListAllMock.Set(func(_ context.Context, page int, pageSize int) ([]domain.Booking, int, error) {
		if page != 2 || pageSize != 10 {
			t.Fatalf("unexpected pagination: page=%d pageSize=%d", page, pageSize)
		}
		return expected, 1, nil
	})
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListByUserMock.Optional()
	bookingsRepo.CancelByOwnerMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", nil }},
		fixedClock{now: time.Now().UTC()},
	)

	got, total, err := uc.ListAll(context.Background(), domain.RoleAdmin, 2, 10)
	if err != nil {
		t.Fatalf("ListAll error: %v", err)
	}
	if total != 1 || len(got) != 1 || got[0].ID != "b1" {
		t.Fatalf("unexpected result: total=%d bookings=%+v", total, got)
	}
}

func TestListMineForbiddenForAdmin(t *testing.T) {
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListAllMock.Optional()
	bookingsRepo.ListByUserMock.Optional()
	bookingsRepo.CancelByOwnerMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", nil }},
		fixedClock{now: time.Now().UTC()},
	)

	_, err := uc.ListMine(context.Background(), "u1", domain.RoleAdmin)
	if err == nil || domain.AsAppError(err).Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %v", err)
	}
}

func TestListMineSuccess(t *testing.T) {
	now := time.Date(2026, 3, 23, 10, 0, 0, 0, time.UTC)
	expected := []domain.Booking{{ID: "b1", UserID: "u1"}}
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	bookingsRepo := repmocks.NewBookingRepositoryMock(t)
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	bookingsRepo.ListByUserMock.Set(func(_ context.Context, userID string, from time.Time) ([]domain.Booking, error) {
		if userID != "u1" {
			t.Fatalf("unexpected userID: %s", userID)
		}
		if !from.Equal(now) {
			t.Fatalf("unexpected from time: %v", from)
		}
		return expected, nil
	})
	bookingsRepo.CreateMock.Optional()
	bookingsRepo.ListAllMock.Optional()
	bookingsRepo.CancelByOwnerMock.Optional()
	bookingsRepo.GetByIDMock.Optional()
	uc := bookingsuc.NewBookingsUsecase(
		slotsRepo,
		bookingsRepo,
		conferenceStub{createFn: func(_ context.Context, _ string) (string, error) { return "", nil }},
		fixedClock{now: now},
	)

	got, err := uc.ListMine(context.Background(), "u1", domain.RoleUser)
	if err != nil {
		t.Fatalf("ListMine error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "b1" {
		t.Fatalf("unexpected result: %+v", got)
	}
}
