package handler

import (
	"context"
	"errors"

	"museum/api/generated"
	"museum/internal/usecase"
)

// ExhibitionHandler は
// OpenAPIから生成されたStrictServerInterfaceを実装する。
//
// HTTP Request / Responseの型は
// OpenAPIから自動生成されるため、
// Handlerでは手動定義しない。
type ExhibitionHandler struct {
	createUseCase *usecase.CreateExhibitionUseCase
}

// NewExhibitionHandler はHandlerを生成する。
func NewExhibitionHandler(
	createUseCase *usecase.CreateExhibitionUseCase,
) *ExhibitionHandler {
	return &ExhibitionHandler{
		createUseCase: createUseCase,
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

// Compile-time check.
//
// ExhibitionHandlerがOpenAPIで要求される
// interfaceを満たしていなければ
// コンパイル時にエラーになる。
var _ generated.StrictServerInterface = (*ExhibitionHandler)(nil)

// errors importを将来使う予定がない場合は削除してOK。
var _ = errors.New
