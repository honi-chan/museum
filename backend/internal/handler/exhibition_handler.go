package handler

import (
	"context"
	"errors"

	"museum/api/generated"
	"museum/internal/repository"
	"museum/internal/usecase"
)

// ExhibitionHandler は
// 展示室のHTTP Request / Response変換を担当する。
//
// HTTP Request / Responseの型は
// OpenAPIから自動生成されるため、
// Handlerでは手動定義しない。
type ExhibitionHandler struct {
	createUseCase *usecase.CreateExhibitionUseCase
	getUseCase    *usecase.GetExhibitionUseCase
	listUseCase   *usecase.ListExhibitionsUseCase
}

// NewExhibitionHandler は
// Exhibition関連のUseCaseを受け取り
// Handlerを生成する。
func NewExhibitionHandler(
	createUseCase *usecase.CreateExhibitionUseCase,
	getUseCase *usecase.GetExhibitionUseCase,
	listUseCase *usecase.ListExhibitionsUseCase,
) *ExhibitionHandler {
	return &ExhibitionHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
	}
}

// --------------------------------------------------
// Create Exhibition
// --------------------------------------------------

// CreateExhibition は
// OpenAPIの
//
//	operationId: createExhibition
//
// に対応するHandler。
//
// Request / Response型はすべて
// oapi-codegenが生成している。
func (h *ExhibitionHandler) CreateExhibition(
	ctx context.Context,
	request generated.CreateExhibitionRequestObject,
) (
	generated.CreateExhibitionResponseObject,
	error,
) {
	// request.BodyもOpenAPIから生成された型。
	//
	// json.NewDecoderは不要。
	if request.Body == nil {
		return generated.CreateExhibition400JSONResponse{
			Error: "request body is required",
		}, nil
	}

	// --------------------------------------------------
	// API Model → UseCase Input
	// --------------------------------------------------

	input := usecase.CreateExhibitionInput{
		MuseumID: request.Body.MuseumId,
		Title:    request.Body.Title,
	}

	// descriptionはOpenAPI上optionalなので
	// pointerとして生成される可能性がある。
	if request.Body.Description != nil {
		input.Description =
			*request.Body.Description
	}

	// --------------------------------------------------
	// UseCase
	// --------------------------------------------------

	created, err := h.createUseCase.Execute(
		ctx,
		input,
	)

	if err != nil {
		// 現時点ではUseCase errorを
		// Bad Requestとして扱う。
		//
		// 後でDomain Errorを導入して、
		// 400 / 404 / 409 / 500を自動分類する。
		return generated.CreateExhibition400JSONResponse{
			Error: err.Error(),
		}, nil
	}

	// --------------------------------------------------
	// Domain → API Response
	// --------------------------------------------------

	return generated.CreateExhibition201JSONResponse{
		Id:          created.ID,
		MuseumId:    created.MuseumID,
		Title:       created.Title,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
	}, nil
}

// --------------------------------------------------
// Get Exhibition
// --------------------------------------------------

// GetExhibition は
//
// GET /exhibitions/{id}
//
// を処理する。
//
// Path parameterやHTTP Response型は
// OpenAPIから自動生成される。
func (
	h *ExhibitionHandler,
) GetExhibition(
	ctx context.Context,
	request generated.GetExhibitionRequestObject,
) (
	generated.GetExhibitionResponseObject,
	error,
) {
	// --------------------------------------------------
	// UseCase
	// --------------------------------------------------

	exhibition, err :=
		h.getUseCase.Execute(
			ctx,
			request.Id,
		)

	if err != nil {
		// Repository共通のNotFoundを
		// HTTP 404へ変換する。
		if errors.Is(
			err,
			repository.ErrNotFound,
		) {
			return generated.GetExhibition404JSONResponse{
				Error: "exhibition not found",
			}, nil
		}

		// 想定外エラーはerrorとして返す。
		//
		// Strict Serverが
		// Internal Server Errorとして処理する。
		return nil, err
	}

	// --------------------------------------------------
	// Domain → API Response
	// --------------------------------------------------

	return generated.GetExhibition200JSONResponse{
		Id:          exhibition.ID,
		MuseumId:    exhibition.MuseumID,
		Title:       exhibition.Title,
		Description: exhibition.Description,
		CreatedAt:   exhibition.CreatedAt,
	}, nil
}

// ListMuseumExhibitions は
// Museumに所属する展示室一覧を返す。
func (
	h *ExhibitionHandler,
) ListMuseumExhibitions(
	ctx context.Context,
	request generated.ListMuseumExhibitionsRequestObject,
) (
	generated.ListMuseumExhibitionsResponseObject,
	error,
) {
	exhibitions, err :=
		h.listUseCase.Execute(
			ctx,
			request.MuseumId,
		)

	if err != nil {
		return nil, err
	}

	response :=
		make(
			[]generated.ExhibitionResponse,
			0,
			len(exhibitions),
		)

	for _, exhibition := range exhibitions {

		response = append(
			response,
			generated.ExhibitionResponse{
				Id:           exhibition.ID,
				MuseumId:     exhibition.MuseumID,
				Title:        exhibition.Title,
				Description:  exhibition.Description,
				DisplayOrder: exhibition.DisplayOrder,
				CreatedAt:    exhibition.CreatedAt,
			},
		)
	}

	return generated.ListMuseumExhibitions200JSONResponse{
		Exhibitions: response,
	}, nil
}

// --------------------------------------------------
// Health
// --------------------------------------------------

// GetHealth はAPIの生存確認を行う。
func (h *ExhibitionHandler) GetHealth(
	ctx context.Context,
	request generated.GetHealthRequestObject,
) (
	generated.GetHealthResponseObject,
	error,
) {
	return generated.GetHealth200JSONResponse{
		Status: "ok",
	}, nil
}
