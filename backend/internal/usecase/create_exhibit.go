package usecase

import (
	"context"
	"errors"
	"museum/internal/domain"
	"museum/internal/repository"
	"strings"
	"time"
)

var ErrExhibitImageRequired = errors.New("exhibit image_url is required")

type CreateExhibitInput struct {
	ExhibitionID string
	Title        string
	Caption      string
	ImageURL     string
}

type CreateExhibitUseCase struct {
	repository  repository.ExhibitRepository
	idGenerator IDGenerator
}

func NewCreateExhibitUseCase(repo repository.ExhibitRepository, ids IDGenerator) *CreateExhibitUseCase {
	return &CreateExhibitUseCase{repository: repo, idGenerator: ids}
}

func (u *CreateExhibitUseCase) Execute(ctx context.Context, input CreateExhibitInput) (domain.Exhibit, error) {
	imageURL := strings.TrimSpace(input.ImageURL)
	if imageURL == "" {
		return domain.Exhibit{}, ErrExhibitImageRequired
	}
	return u.repository.Save(ctx, domain.Exhibit{
		ID: u.idGenerator.Generate(), ExhibitionID: input.ExhibitionID,
		Title: input.Title, Caption: input.Caption, ImageURL: imageURL, CreatedAt: time.Now(),
	})
}
