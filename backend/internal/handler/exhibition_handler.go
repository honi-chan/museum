package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"museum/internal/usecase"
)

// ExhibitionHandler は
// Exhibitionに関するHTTP通信を担当する。
//
// ビジネスルールはここには書かない。
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
// HTTP Request専用モデル。
//
// DomainやUseCaseとは分離する。
type createExhibitionRequest struct {
	MuseumID    string `json:"museum_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// createExhibitionResponse は
// HTTP Response専用モデル。
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

	// JSON Requestを読み込む。
	if err := json.NewDecoder(
		r.Body,
	).Decode(&request); err != nil {
		http.Error(
			w,
			`{"error":"invalid request body"}`,
			http.StatusBadRequest,
		)

		return
	}

	// HTTP RequestをUseCase用Inputへ変換する。
	input := usecase.CreateExhibitionInput{
		MuseumID:    request.MuseumID,
		Title:       request.Title,
		Description: request.Description,
	}

	// ビジネス処理はUseCaseへ任せる。
	created, err := h.createUseCase.Execute(
		r.Context(),
		input,
	)

	if err != nil {
		http.Error(
			w,
			`{"error":"`+err.Error()+`"}`,
			http.StatusBadRequest,
		)

		return
	}

	// DomainをHTTP Responseへ変換する。
	response := createExhibitionResponse{
		ID:          created.ID,
		MuseumID:    created.MuseumID,
		Title:       created.Title,
		Description: created.Description,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	if err := json.NewEncoder(
		w,
	).Encode(response); err != nil {
		log.Printf(
			"failed to encode response: %v",
			err,
		)
	}
}
