package usecase

import (
	"context"
	"errors"
	"museum/internal/domain"
	"museum/internal/repository"
	"testing"
)

type fakeExhibitRepository struct {
	saved domain.Exhibit
	items []domain.Exhibit
	err   error
	calls int
}

func (r *fakeExhibitRepository) Save(_ context.Context, item domain.Exhibit) (domain.Exhibit, error) {
	r.calls++
	item.DisplayOrder = 7
	r.saved = item
	return item, r.err
}
func (r *fakeExhibitRepository) ListByExhibitionID(_ context.Context, _ string) ([]domain.Exhibit, error) {
	r.calls++
	return r.items, r.err
}
func TestCreateExhibit(t *testing.T) {
	for _, url := range []string{"", " \n\t"} {
		repo := &fakeExhibitRepository{}
		_, err := NewCreateExhibitUseCase(repo, fixedIDGenerator{id: "work"}).Execute(context.Background(), CreateExhibitInput{ImageURL: url})
		if !errors.Is(err, ErrExhibitImageRequired) || repo.calls != 0 {
			t.Fatalf("empty image accepted: %v", err)
		}
	}
	repo := &fakeExhibitRepository{}
	got, err := NewCreateExhibitUseCase(repo, fixedIDGenerator{id: "work"}).Execute(context.Background(), CreateExhibitInput{ExhibitionID: "room", ImageURL: " https://example.com/a.jpg "})
	if err != nil || got.ID != "work" || got.ExhibitionID != "room" || got.DisplayOrder != 7 || got.ImageURL != "https://example.com/a.jpg" || got.Title != "" || got.Caption != "" || got.CreatedAt.IsZero() {
		t.Fatalf("unexpected result: %+v, %v", got, err)
	}
	for _, want := range []error{repository.ErrNotFound, errors.New("database unavailable")} {
		repo.err = want
		_, err = NewCreateExhibitUseCase(repo, fixedIDGenerator{}).Execute(context.Background(), CreateExhibitInput{ImageURL: "https://example.com/a.jpg"})
		if !errors.Is(err, want) {
			t.Fatalf("got %v, want %v", err, want)
		}
	}
}
func TestListExhibits(t *testing.T) {
	exhibits := &fakeExhibitRepository{items: []domain.Exhibit{{ID: "b", DisplayOrder: 1}, {ID: "a", DisplayOrder: 2}}}
	rooms := &fakeExhibitionRepository{}
	u := NewListExhibitsUseCase(exhibits, rooms)
	got, err := u.Execute(context.Background(), "room")
	if err != nil || len(got) != 2 || got[0].ID != "b" {
		t.Fatalf("order changed: %+v %v", got, err)
	}
	rooms.err = repository.ErrNotFound
	exhibits.calls = 0
	if _, err = u.Execute(context.Background(), "missing"); !errors.Is(err, repository.ErrNotFound) || exhibits.calls != 0 {
		t.Fatalf("missing room: %v", err)
	}
	rooms.err = nil
	exhibits.err = errors.New("offline")
	if _, err = u.Execute(context.Background(), "room"); !errors.Is(err, exhibits.err) {
		t.Fatalf("list error: %v", err)
	}
}
