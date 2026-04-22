package rooms_test

import (
	"context"
	"errors"
	"testing"

	"github.com/M1steryO/Room-Booking-Service/internal/domain"
	repmocks "github.com/M1steryO/Room-Booking-Service/internal/repository/mocks"
	roomsuc "github.com/M1steryO/Room-Booking-Service/internal/usecase/rooms"
)

func TestCreateForbiddenForUser(t *testing.T) {
	repo := repmocks.NewRoomRepositoryMock(t)
	repo.CreateMock.Optional()
	repo.ListMock.Optional()
	repo.GetByIDMock.Optional()
	uc := roomsuc.NewRoomsUsecase(repo)

	_, err := uc.Create(context.Background(), domain.RoleUser, "A", nil, nil)
	if err == nil {
		t.Fatal("expected forbidden error")
	}
	if domain.AsAppError(err).Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN, got %s", domain.AsAppError(err).Code)
	}
}

func TestCreateSuccess(t *testing.T) {
	var created domain.Room
	description := "Quiet room"
	capacity := 8
	repo := repmocks.NewRoomRepositoryMock(t)
	repo.CreateMock.Set(func(_ context.Context, room domain.Room) (domain.Room, error) {
		created = room
		room.CreatedAt = created.CreatedAt
		return room, nil
	})
	repo.ListMock.Optional()
	repo.GetByIDMock.Optional()
	uc := roomsuc.NewRoomsUsecase(repo)

	room, err := uc.Create(context.Background(), domain.RoleAdmin, "Blue", &description, &capacity)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if created.Name != "Blue" {
		t.Fatalf("unexpected name passed to repo: %s", created.Name)
	}
	if created.Description == nil || *created.Description != description {
		t.Fatalf("unexpected description passed to repo: %+v", created.Description)
	}
	if created.Capacity == nil || *created.Capacity != capacity {
		t.Fatalf("unexpected capacity passed to repo: %+v", created.Capacity)
	}
	if created.ID == "" {
		t.Fatal("expected generated room id to be passed to repo")
	}

	if room.ID != created.ID || room.Name != created.Name {
		t.Fatalf("expected usecase to return repo result: got=%+v repo=%+v", room, created)
	}
}

func TestList(t *testing.T) {
	expected := []domain.Room{{ID: "r1", Name: "R1"}}
	repo := repmocks.NewRoomRepositoryMock(t)
	repo.ListMock.Set(func(_ context.Context) ([]domain.Room, error) { return expected, nil })
	repo.CreateMock.Optional()
	repo.GetByIDMock.Optional()
	uc := roomsuc.NewRoomsUsecase(repo)

	got, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestCreateRepoError(t *testing.T) {
	repoErr := errors.New("insert failed")
	repo := repmocks.NewRoomRepositoryMock(t)
	repo.CreateMock.Set(func(_ context.Context, _ domain.Room) (domain.Room, error) {
		return domain.Room{}, repoErr
	})
	repo.ListMock.Optional()
	repo.GetByIDMock.Optional()
	uc := roomsuc.NewRoomsUsecase(repo)

	_, err := uc.Create(context.Background(), domain.RoleAdmin, "Blue", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestListRepoError(t *testing.T) {
	repoErr := errors.New("list failed")
	repo := repmocks.NewRoomRepositoryMock(t)
	repo.ListMock.Set(func(_ context.Context) ([]domain.Room, error) { return nil, repoErr })
	repo.CreateMock.Optional()
	repo.GetByIDMock.Optional()
	uc := roomsuc.NewRoomsUsecase(repo)

	_, err := uc.List(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
