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

	// 動作確認用。
	mux.HandleFunc("GET /health", health)

	// 最初のMUSEUM機能。
	mux.HandleFunc("POST /exhibitions", createExhibition)

	log.Println("server started on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func createExhibition(w http.ResponseWriter, r *http.Request) {
	var request createExhibitionRequest

	// JSON
	//
	// {
	//   "museum_id": "...",
	//   "title": "...",
	//   "description": "..."
	// }
	//
	// をGoのstructへ変換する。
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			`{"error":"invalid request body"}`,
			http.StatusBadRequest,
		)
		return
	}

	// HTTP用モデルから、
	// MUSEUMのCreateInputへ変換する。
	input := exhibition.CreateInput{
		MuseumID:    request.MuseumID,
		Title:       request.Title,
		Description: request.Description,
	}

	created, err := exhibition.Create(input)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
