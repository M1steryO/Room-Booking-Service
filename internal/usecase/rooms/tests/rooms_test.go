package rooms_test

import (
	"context"
	"testing"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	repmocks "github.com/avito-internships/test-backend-1-M1steryO/internal/repository/mocks"
	roomsuc "github.com/avito-internships/test-backend-1-M1steryO/internal/usecase/rooms"
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
	repo := repmocks.NewRoomRepositoryMock(t)
	repo.CreateMock.Set(func(_ context.Context, room domain.Room) (domain.Room, error) {
		created = room
		return room, nil
	})
	repo.ListMock.Optional()
	repo.GetByIDMock.Optional()
	uc := roomsuc.NewRoomsUsecase(repo)

	room, err := uc.Create(context.Background(), domain.RoleAdmin, "Blue", nil, nil)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if room.Name != "Blue" || created.Name != "Blue" {
		t.Fatal("expected room to be passed to repository")
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
