package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"museum/internal/domain"
	"museum/internal/repository"
)

// CreateExhibitionInput は
// 展示室作成ユースケースへの入力。
//
// HTTP Requestではない。
//
// UIやCLIなど、どこから呼び出されても
// 同じ形式で扱える。
type CreateExhibitionInput struct {
	MuseumID    string
	Title       string
	Description string
}

// CreateExhibitionUseCase は
// 「展示室を作成する」という処理を担当する。
type CreateExhibitionUseCase struct {
	repository  repository.ExhibitionRepository
	idGenerator IDGenerator
}

// NewCreateExhibitionUseCase は
// UseCaseに必要な依存関係を受け取る。
//
// 依存関係はここで一度だけ注入する。
func NewCreateExhibitionUseCase(
	repository repository.ExhibitionRepository,
	idGenerator IDGenerator,
) *CreateExhibitionUseCase {
	return &CreateExhibitionUseCase{
		repository:  repository,
		idGenerator: idGenerator,
	}
}

// Execute は展示室作成処理を実行する。
func (u *CreateExhibitionUseCase) Execute(
	ctx context.Context,
	input CreateExhibitionInput,
) (domain.Exhibition, error) {
	// MUSEUMのルール:
	// Exhibitionにはタイトルが必要。
	title := strings.TrimSpace(input.Title)

	if title == "" {
		return domain.Exhibition{},
			errors.New("exhibition title is required")
	}

	// ID生成方法はUseCase自身では知らない。
	id := u.idGenerator.Generate()

	exhibition := domain.Exhibition{
		ID:          id,
		MuseumID:    input.MuseumID,
		Title:       title,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}

	// 保存先の実装もUseCaseは知らない。
	if err := u.repository.Save(
		ctx,
		exhibition,
	); err != nil {
		return domain.Exhibition{}, err
	}

	return exhibition, nil
}
