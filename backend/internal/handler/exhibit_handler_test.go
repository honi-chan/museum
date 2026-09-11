package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"museum/api/generated"
	"museum/internal/domain"
	"museum/internal/repository"
	"museum/internal/usecase"
)

type exhibitStore struct {
	err   error
	items []domain.Exhibit
}

func (s *exhibitStore) Save(_ context.Context, item domain.Exhibit) (domain.Exhibit, error) {
	item.DisplayOrder = 9
	return item, s.err
}
func (s *exhibitStore) ListByExhibitionID(context.Context, string) ([]domain.Exhibit, error) {
	return s.items, s.err
}

type roomStore struct{ err error }

func (s roomStore) Save(context.Context, domain.Exhibition) error { return s.err }
func (s roomStore) FindByID(context.Context, string) (domain.Exhibition, error) {
	return domain.Exhibition{}, s.err
}
func (s roomStore) ListByMuseumID(context.Context, string) ([]domain.Exhibition, error) {
	return nil, s.err
}

type testID struct{}

func (testID) Generate() string { return "work" }

func TestExhibitHTTP(t *testing.T) {
	for _, tc := range []struct {
		name, method, body string
		storeErr, roomErr  error
		status             int
	}{
		{"create", "POST", `{"image_url":"https://example.com/a.jpg"}`, nil, nil, 201},
		{"empty", "POST", `{"image_url":" "}`, nil, nil, 400},
		{"missing field", "POST", `{}`, nil, nil, 400},
		{"null", "POST", `null`, nil, nil, 400},
		{"invalid JSON", "POST", `{`, nil, nil, 400},
		{"missing parent", "POST", `{"image_url":"https://example.com/a.jpg"}`, repository.ErrNotFound, nil, 404},
		{"storage error", "POST", `{"image_url":"https://example.com/a.jpg"}`, errors.New("secret database error"), nil, 500},
		{"empty list", "GET", "", nil, nil, 200},
		{"missing list parent", "GET", "", nil, repository.ErrNotFound, 404},
		{"list error", "GET", "", errors.New("secret database error"), nil, 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &exhibitStore{err: tc.storeErr}
			h := NewAPIHandler(nil, NewExhibitHandler(usecase.NewCreateExhibitUseCase(store, testID{}), usecase.NewListExhibitsUseCase(store, roomStore{err: tc.roomErr})))
			router := generated.Handler(generated.NewStrictHandler(h, nil))
			request := httptest.NewRequest(tc.method, "/exhibitions/room/exhibits", strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != tc.status {
				t.Fatalf("status %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
			if strings.Contains(response.Body.String(), "secret") {
				t.Fatal("internal details leaked")
			}
			if tc.status == 201 {
				var got generated.ExhibitResponse
				if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.DisplayOrder != 9 || got.ExhibitionId != "room" || got.Id != "work" || got.CreatedAt.Before(time.Now().Add(-time.Minute)) {
					t.Fatalf("mapping: %+v", got)
				}
			}
			if tc.status == 200 && response.Body.String() != "{\"exhibits\":[]}\n" {
				t.Fatalf("empty list must be array: %s", response.Body.String())
			}
		})
	}
}
