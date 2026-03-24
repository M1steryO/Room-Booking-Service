package schedules_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	repmocks "github.com/avito-internships/test-backend-1-M1steryO/internal/repository/mocks"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/schedules"
)

type fixedClock struct{ now time.Time }

func (f fixedClock) Now() time.Time { return f.now }

func TestCreateForbiddenForUser(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Optional()
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	uc := schedules.NewSchedulesUsecase(
		roomsRepo,
		schedulesRepo,
		slotsRepo,
		fixedClock{now: time.Now().UTC()},
		7,
	)

	_, err := uc.Create(context.Background(), domain.RoleUser, "r1", []int{1}, "09:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %v", err)
	}
}

func TestCreateSuccessGeneratesSlots(t *testing.T) {
	var inserted int
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	schedulesRepo.CreateWithSlotsMock.Set(func(_ context.Context, s domain.Schedule, slots []domain.Slot) (domain.Schedule, error) {
		inserted = len(slots)
		return s, nil
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()
	uc := schedules.NewSchedulesUsecase(
		roomsRepo,
		schedulesRepo,
		slotsRepo,
		fixedClock{now: time.Date(2026, 3, 23, 8, 0, 0, 0, time.UTC)},
		2,
	)

	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1}, "09:00", "10:00")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if inserted == 0 {
		t.Fatal("expected generated slots to be passed to transactional create")
	}
}

func TestCreateRoomNotFound(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) {
		return domain.Room{}, domain.RoomNotFound()
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "missing-room", []int{1}, "09:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "ROOM_NOT_FOUND" {
		t.Fatalf("expected ROOM_NOT_FOUND, got %v", err)
	}
}

func TestCreateGetByIDRepoError(t *testing.T) {
	repoErr := errors.New("db read failed")
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) {
		return domain.Room{}, repoErr
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1}, "09:00", "10:00")
	if err == nil || !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestCreateWithSlotsRepoError(t *testing.T) {
	repoErr := errors.New("insert failed")
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) {
		return domain.Room{ID: "r1"}, nil
	})
	schedulesRepo.CreateWithSlotsMock.Set(func(_ context.Context, _ domain.Schedule, _ []domain.Slot) (domain.Schedule, error) {
		return domain.Schedule{}, repoErr
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1}, "09:00", "10:00")
	if err == nil || !errors.Is(err, repoErr) {
		t.Fatalf("expected CreateWithSlots error, got %v", err)
	}
}

func TestCreateEmptyRoomID(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, roomID string) (domain.Room, error) {
		if roomID != "" {
			t.Fatalf("expected empty roomID, got %q", roomID)
		}
		return domain.Room{}, domain.RoomNotFound()
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "", []int{1}, "09:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "ROOM_NOT_FOUND" {
		t.Fatalf("expected ROOM_NOT_FOUND for empty roomID, got %v", err)
	}
}

func TestCreateEmptyDaysValidation(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{}, "09:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %v", err)
	}
}

func TestCreateInvalidDayOfWeekValidation(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{0}, "09:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %v", err)
	}
}

func TestCreateInvalidStartTimeValidation(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1}, "aa:bb", "10:00")
	if err == nil || domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %v", err)
	}
}

func TestCreateEndBeforeOrEqualStartValidation(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.CreateWithSlotsMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1}, "10:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "INVALID_REQUEST" {
		t.Fatalf("expected INVALID_REQUEST, got %v", err)
	}
}

func TestCreateConflictFromRepository(t *testing.T) {
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	schedulesRepo.CreateWithSlotsMock.Set(func(_ context.Context, _ domain.Schedule, _ []domain.Slot) (domain.Schedule, error) {
		return domain.Schedule{}, domain.ScheduleExists()
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: time.Now().UTC()}, 7)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1}, "09:00", "10:00")
	if err == nil || domain.AsAppError(err).Code != "SCHEDULE_EXISTS" {
		t.Fatalf("expected SCHEDULE_EXISTS, got %v", err)
	}
}

func TestCreateRespectsPlanningHorizonTwoDays(t *testing.T) {
	now := time.Date(2026, 3, 23, 8, 0, 0, 0, time.UTC) // Monday
	roomsRepo := repmocks.NewRoomRepositoryMock(t)
	schedulesRepo := repmocks.NewScheduleRepositoryMock(t)
	slotsRepo := repmocks.NewSlotRepositoryMock(t)
	roomsRepo.GetByIDMock.Set(func(_ context.Context, _ string) (domain.Room, error) { return domain.Room{ID: "r1"}, nil })
	schedulesRepo.CreateWithSlotsMock.Set(func(_ context.Context, _ domain.Schedule, slots []domain.Slot) (domain.Schedule, error) {
		if len(slots) != 4 {
			t.Fatalf("expected 4 slots for 2-day horizon and 1h window, got %d", len(slots))
		}
		if got := slots[0].Start.UTC(); !got.Equal(time.Date(2026, 3, 23, 9, 0, 0, 0, time.UTC)) {
			t.Fatalf("unexpected first slot start: %v", got)
		}
		if got := slots[2].Start.UTC(); !got.Equal(time.Date(2026, 3, 24, 9, 0, 0, 0, time.UTC)) {
			t.Fatalf("unexpected third slot start: %v", got)
		}
		return domain.Schedule{ID: "s1", RoomID: "r1", DaysOfWeek: []int{1, 2}, StartTime: "09:00", EndTime: "10:00", CreatedAt: now}, nil
	})
	roomsRepo.CreateMock.Optional()
	roomsRepo.ListMock.Optional()
	schedulesRepo.GetByRoomIDMock.Optional()
	slotsRepo.BulkUpsertMock.Optional()
	slotsRepo.DateHasSlotsMock.Optional()
	slotsRepo.GetByIDMock.Optional()
	slotsRepo.ListAvailableByRoomAndDateMock.Optional()

	uc := schedules.NewSchedulesUsecase(roomsRepo, schedulesRepo, slotsRepo, fixedClock{now: now}, 2)
	_, err := uc.Create(context.Background(), domain.RoleAdmin, "r1", []int{1, 2}, "09:00", "10:00")
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
}
