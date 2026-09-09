package handler

import (
	"errors"
	"log"
	"net/http"

	"museum/internal/handler/httpx"
	"museum/internal/usecase"
)

// ExhibitionHandler は
// Exhibitionに関するHTTP通信を担当する。
//
// ビジネスルールはここには書かない。
//
// Handlerの責務は、
//
// HTTP Request
// ↓
// UseCase Input
// ↓
// UseCase実行
// ↓
// HTTP Response
//
// への変換だけにする。
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

// createExhibitionRequest は
// HTTP Request専用のデータ構造。
//
// DomainモデルをそのままHTTPに公開しない。
type createExhibitionRequest struct {
	MuseumID    string `json:"museum_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// createExhibitionResponse は
// HTTP Response専用のデータ構造。
//
// APIとして何を公開するかは
// Handler側で明示的に決める。
type createExhibitionResponse struct {
	ID          string `json:"id"`
	MuseumID    string `json:"museum_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Create はPOST /exhibitionsを処理する。
func (h *ExhibitionHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request createExhibitionRequest

	// -----------------------------
	// Request Decode
	// -----------------------------

	// HTTP Request BodyをJSONから構造体へ変換する。
	//
	// JSON処理そのものはhttpxへ委譲する。
	if err := httpx.DecodeJSON(
		r,
		&request,
	); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			errors.New("invalid request body"),
		)

		return
	}

	// -----------------------------
	// HTTP → UseCase
	// -----------------------------

	// HTTP専用Requestから
	// UseCase専用Inputへ変換する。
	input := usecase.CreateExhibitionInput{
		MuseumID:    request.MuseumID,
		Title:       request.Title,
		Description: request.Description,
	}

	// -----------------------------
	// UseCase
	// -----------------------------

	created, err := h.createUseCase.Execute(
		r.Context(),
		input,
	)

	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	// -----------------------------
	// Domain → HTTP
	// -----------------------------

	response := createExhibitionResponse{
		ID:          created.ID,
		MuseumID:    created.MuseumID,
		Title:       created.Title,
		Description: created.Description,
	}

	// -----------------------------
	// Response
	// -----------------------------

	if err := httpx.WriteJSON(
		w,
		http.StatusCreated,
		response,
	); err != nil {
		log.Printf(
			"failed to write response: %v",
			err,
		)
	}
}

// writeError はHandler内部で使う
// エラー出力の補助関数。
//
// JSON書き込み自体が失敗した場合だけ
// Server側のログとして記録する。
func writeError(
	w http.ResponseWriter,
	status int,
	err error,
) {
	if writeErr := httpx.WriteError(
		w,
		status,
		err,
	); writeErr != nil {
		log.Printf(
			"failed to write error response: %v",
			writeErr,
		)
	}
}
