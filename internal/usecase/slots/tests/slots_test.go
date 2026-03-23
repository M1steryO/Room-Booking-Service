package slots_test

import (
	"context"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	repmocks "github.com/avito-internships/test-backend-1-M1steryO/internal/repository/mocks"
	slotsuc "github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/slots"
)

func TestListAvailableNoScheduleReturnsEmpty(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	schedulesRepo.GetByRoomIDMock.Set(func(_ context.Context, _ string) (domain.Schedule, error) {
		return domain.Schedule{}, domain.NotFound("schedule not found")
	})
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateMock.Optional()
	uc := slotsuc.NewSlotsUsecase(roomsRepo, schedulesRepo, slotsRepo)

	got, err := uc.ListAvailable(context.Background(), "r1", time.Now().UTC())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty list, got %+v", got)
	}
}

func TestListAvailableGeneratesWhenMissing(t *testing.T) {
	var upsertCalled bool
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	schedulesRepo.GetByRoomIDMock.Set(func(_ context.Context, _ string) (domain.Schedule, error) {
		return domain.Schedule{RoomID: "r1", DaysOfWeek: []int{1}, StartTime: "09:00", EndTime: "10:00"}, nil
	})
	slotsRepo.DateHasSlotsMock.Set(func(_ context.Context, _ string, _ time.Time) (bool, error) { return false, nil })
	slotsRepo.BulkUpsertMock.Set(func(_ context.Context, slots []domain.Slot) error {
		upsertCalled = len(slots) > 0
		return nil
	})
	slotsRepo.ListAvailableByRoomAndDateMock.Set(func(_ context.Context, _ string, _ time.Time) ([]domain.Slot, error) {
		return []domain.Slot{{ID: "s1"}}, nil
	})
	slotsRepo.GetByIDMock.Optional()
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateMock.Optional()
	uc := slotsuc.NewSlotsUsecase(roomsRepo, schedulesRepo, slotsRepo)

	date := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	got, err := uc.ListAvailable(context.Background(), "r1", date)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !upsertCalled {
		t.Fatal("expected bulk upsert to be called")
	}
	if len(got) != 1 || got[0].ID != "s1" {
		t.Fatalf("unexpected slots: %+v", got)
	}
}
