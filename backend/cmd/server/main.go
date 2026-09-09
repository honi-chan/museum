package main

import (
	"log"
	"net/http"

	"museum/internal/handler"
	exhibitioninfra "museum/internal/infrastructure/exhibition"
	"museum/internal/usecase"
)

func main() {
	// --------------------------------
	// Infrastructure
	// --------------------------------

	// 現在はMemoryRepository。
	//
	// 将来ここをPostgreSQLRepositoryへ変更しても、
	// HandlerやUseCaseは変更しない。
	exhibitionRepository :=
		exhibitioninfra.NewMemoryRepository()

	// ID生成の具体実装。
	idGenerator :=
		exhibitioninfra.UUIDGenerator{}

	// --------------------------------
	// UseCase
	// --------------------------------

	createExhibitionUseCase :=
		usecase.NewCreateExhibitionUseCase(
			exhibitionRepository,
			idGenerator,
		)

	// --------------------------------
	// Handler
	// --------------------------------

	exhibitionHandler :=
		handler.NewExhibitionHandler(
			createExhibitionUseCase,
		)

	// --------------------------------
	// Router
	// --------------------------------

	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		health,
	)

	mux.HandleFunc(
		"POST /exhibitions",
		exhibitionHandler.Create,
	)

	log.Println("server started on :8080")

	if err := http.ListenAndServe(
		":8080",
		mux,
	); err != nil {
		log.Fatal(err)
	}
}

// health はサーバー生存確認用。
func health(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	_, _ = w.Write(
		[]byte(`{"status":"ok"}`),
	)
}
