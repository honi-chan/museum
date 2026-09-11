package handler

import (
	"context"
	"errors"
	"museum/api/generated"
	"museum/internal/domain"
	"museum/internal/repository"
	"museum/internal/usecase"
)

type ExhibitHandler struct {
	create *usecase.CreateExhibitUseCase
	list   *usecase.ListExhibitsUseCase
}

func NewExhibitHandler(create *usecase.CreateExhibitUseCase, list *usecase.ListExhibitsUseCase) *ExhibitHandler {
	return &ExhibitHandler{create: create, list: list}
}

func (h *ExhibitHandler) CreateExhibit(ctx context.Context, request generated.CreateExhibitRequestObject) (generated.CreateExhibitResponseObject, error) {
	if request.Body == nil {
		return generated.CreateExhibit400JSONResponse{Error: "request body is required"}, nil
	}
	input := usecase.CreateExhibitInput{ExhibitionID: request.ExhibitionId, ImageURL: request.Body.ImageUrl}
	if request.Body.Title != nil {
		input.Title = *request.Body.Title
	}
	if request.Body.Caption != nil {
		input.Caption = *request.Body.Caption
	}
	created, err := h.create.Execute(ctx, input)
	if errors.Is(err, usecase.ErrExhibitImageRequired) {
		return generated.CreateExhibit400JSONResponse{Error: err.Error()}, nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return generated.CreateExhibit404JSONResponse{Error: "exhibition not found"}, nil
	}
	if err != nil {
		return generated.CreateExhibit500JSONResponse{Error: "could not create exhibit"}, nil
	}
	return generated.CreateExhibit201JSONResponse(toExhibitResponse(created)), nil
}

func (h *ExhibitHandler) ListExhibits(ctx context.Context, request generated.ListExhibitsRequestObject) (generated.ListExhibitsResponseObject, error) {
	exhibits, err := h.list.Execute(ctx, request.ExhibitionId)
	if errors.Is(err, repository.ErrNotFound) {
		return generated.ListExhibits404JSONResponse{Error: "exhibition not found"}, nil
	}
	if err != nil {
		return generated.ListExhibits500JSONResponse{Error: "could not list exhibits"}, nil
	}
	response := make([]generated.ExhibitResponse, 0, len(exhibits))
	for _, exhibit := range exhibits {
		response = append(response, toExhibitResponse(exhibit))
	}
	return generated.ListExhibits200JSONResponse{Exhibits: response}, nil
}

func toExhibitResponse(exhibit domain.Exhibit) generated.ExhibitResponse {
	return generated.ExhibitResponse{Id: exhibit.ID, ExhibitionId: exhibit.ExhibitionID,
		Title: exhibit.Title, Caption: exhibit.Caption, ImageUrl: exhibit.ImageURL,
		DisplayOrder: exhibit.DisplayOrder, CreatedAt: exhibit.CreatedAt}
}
