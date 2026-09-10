package main

import (
	"log"
	"net/http"

	"museum/api/generated"
	"museum/internal/handler"
	exhibitioninfra "museum/internal/infrastructure/exhibition"
	"museum/internal/usecase"
)

func main() {
	// ==================================================
	// Infrastructure
	// ==================================================

	// 現在はMemoryRepository。
	//
	// 後でPostgreSQLRepositoryへ差し替える。
	exhibitionRepository :=
		exhibitioninfra.NewMemoryRepository()

	idGenerator :=
		exhibitioninfra.UUIDGenerator{}

	// ==================================================
	// UseCase
	// ==================================================

	createExhibitionUseCase :=
		usecase.NewCreateExhibitionUseCase(
			exhibitionRepository,
			idGenerator,
		)

	// ==================================================
	// Handler
	// ==================================================

	server :=
		handler.NewExhibitionHandler(
			createExhibitionUseCase,
		)

	// ==================================================
	// OpenAPI Strict Server
	// ==================================================

	// HandlerをStrict Serverでラップする。
	//
	// OpenAPIで定義されたRequest / Responseを
	// HTTPへ変換する処理は自動生成コードが担当する。
	strictHandler :=
		generated.NewStrictHandler(
			server,
			nil,
		)

	// OpenAPIから生成されたrouterを使用する。
	//
	// ここで
	//
	// mux.HandleFunc(...)
	//
	// をAPIごとに書かなくてよい。
	httpHandler :=
		generated.Handler(
			strictHandler,
		)

	// ==================================================
	// HTTP Server
	// ==================================================

	log.Println(
		"server started on :8080",
	)

	if err := http.ListenAndServe(
		":8080",
		httpHandler,
	); err != nil {
		log.Fatal(err)
	}
}
