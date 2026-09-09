package main

import (
	"encoding/json"
	"log"
	"net/http"

	"museum/internal/exhibition"
)

// createExhibitionRequest はHTTP API専用の入力。
//
// exhibition.CreateInputとは分ける。
// HTTPのJSON仕様をドメイン側に持ち込まないため。
type createExhibitionRequest struct {
	MuseumID    string `json:"museum_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// createExhibitionResponse はHTTP API専用のレスポンス。
type createExhibitionResponse struct {
	ID          string `json:"id"`
	MuseumID    string `json:"museum_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	mux := http.NewServeMux()

	// Exhibitionの保存先。
	//
	// 現段階ではメモリ保存。
	// 後でPostgreSQL Repositoryへ差し替える。
	repository := exhibition.NewMemoryRepository()

	idGenerator := exhibition.UUIDGenerator{}

	mux.HandleFunc(
		"GET /health",
		health,
	)

	mux.HandleFunc(
		"POST /exhibitions",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			createExhibition(
				w,
				r,
				idGenerator,
				repository,
			)
		},
	)

	log.Println("server started on :8080")

	if err := http.ListenAndServe(
		":8080",
		mux,
	); err != nil {
		log.Fatal(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func createExhibition(
	w http.ResponseWriter,
	r *http.Request,
	idGenerator exhibition.IDGenerator,
	repository exhibition.Repository,
) {
	var request createExhibitionRequest

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

	input := exhibition.CreateInput{
		MuseumID:    request.MuseumID,
		Title:       request.Title,
		Description: request.Description,
	}

	// HTTP Handler自身では、
	// Exhibitionの作成方法や保存方法を知らない。
	//
	// 必要な依存関係を渡して、
	// CreateAndSave()へ処理を委譲する。
	created, err := exhibition.CreateAndSave(
		r.Context(),
		input,
		idGenerator,
		repository,
	)

	if err != nil {
		http.Error(
			w,
			`{"error":"`+err.Error()+`"}`,
			http.StatusBadRequest,
		)

		return
	}

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
