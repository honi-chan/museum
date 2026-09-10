package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"museum/internal/domain"
	"museum/internal/repository"
)

func TestGetExhibitionUseCase_Execute(
	t *testing.T,
) {
	// Arrange
	expected := domain.Exhibition{
		ID:          "exhibition-001",
		MuseumID:    "museum-001",
		Title:       "THINGS I MADE",
		Description: "Things I created.",
		CreatedAt: time.Date(
			2026,
			9,
			10,
			0,
			0,
			0,
			0,
			time.UTC,
		),
	}

	repository :=
		&fakeExhibitionRepository{
			found: expected,
		}

	useCase :=
		NewGetExhibitionUseCase(
			repository,
		)

	// Act
	actual, err :=
		useCase.Execute(
			context.Background(),
			"exhibition-001",
		)

	// Assert
	if err != nil {
		t.Fatalf(
			"Execute() returned error: %v",
			err,
		)
	}

	if actual.ID != expected.ID {
		t.Errorf(
			"ID = %q, want %q",
			actual.ID,
			expected.ID,
		)
	}

	if actual.Title != expected.Title {
		t.Errorf(
			"Title = %q, want %q",
			actual.Title,
			expected.Title,
		)
	}
}

func TestGetExhibitionUseCase_Execute_NotFound(
	t *testing.T,
) {
	// Arrange
	repositoryMock :=
		&fakeExhibitionRepository{
			err: repository.ErrNotFound,
		}

	useCase :=
		NewGetExhibitionUseCase(
			repositoryMock,
		)

	// Act
	_, err :=
		useCase.Execute(
			context.Background(),
			"not-found",
		)

	// Assert
	if !errors.Is(
		err,
		repository.ErrNotFound,
	) {
		t.Fatalf(
			"error = %v, want ErrNotFound",
			err,
		)
	}
}
