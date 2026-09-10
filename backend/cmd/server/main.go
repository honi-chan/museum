package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"museum/api/generated"
	"museum/internal/handler"
	"museum/internal/infrastructure/database"
	exhibitioninfra "museum/internal/infrastructure/exhibition"
	"museum/internal/usecase"
)

func main() {
	ctx := context.Background()

	// ==================================================
	// Configuration
	// ==================================================

	// DB情報をコードに直接書かない。
	//
	// Local / CI / Productionで
	// 同じコードを利用できるように環境変数から取得する。
	databaseURL := os.Getenv(
		"DATABASE_URL",
	)

	if databaseURL == "" {
		log.Fatal(
			"DATABASE_URL is required",
		)
	}

	// ==================================================
	// Database
	// ==================================================

	postgresPool, err :=
		database.NewPostgres(
			ctx,
			databaseURL,
		)

	if err != nil {
		log.Fatal(err)
	}

	defer postgresPool.Close()

	// ==================================================
	// Infrastructure
	// ==================================================

	// MemoryRepositoryから
	// PostgreSQLRepositoryへ変更。
	//
	// repository interfaceは変わらないため、
	// UseCase側の修正は不要。
	exhibitionRepository :=
		exhibitioninfra.NewPostgresRepository(
			postgresPool,
		)

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
	// OpenAPI Server
	// ==================================================

	strictHandler :=
		generated.NewStrictHandler(
			server,
			nil,
		)

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
