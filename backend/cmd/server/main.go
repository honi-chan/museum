package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"museum/api/generated"
	"museum/internal/handler"
	"museum/internal/handler/middleware"
	"museum/internal/infrastructure/database"
	exhibitinfra "museum/internal/infrastructure/exhibit"
	exhibitioninfra "museum/internal/infrastructure/exhibition"
	"museum/internal/usecase"
)

func main() {
	ctx := context.Background()

	databaseURL :=
		os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal(
			"DATABASE_URL is required",
		)
	}

	// ----------------------------------------
	// Database
	// ----------------------------------------

	postgresPool, err :=
		database.NewPostgres(
			ctx,
			databaseURL,
		)

	if err != nil {
		log.Fatal(err)
	}

	defer postgresPool.Close()

	// ----------------------------------------
	// Infrastructure
	// ----------------------------------------

	exhibitionRepository :=
		exhibitioninfra.NewPostgresRepository(
			postgresPool,
		)

	idGenerator :=
		exhibitioninfra.UUIDGenerator{}

	// ----------------------------------------
	// UseCase
	// ----------------------------------------

	createExhibitionUseCase :=
		usecase.NewCreateExhibitionUseCase(
			exhibitionRepository,
			idGenerator,
		)

	getExhibitionUseCase :=
		usecase.NewGetExhibitionUseCase(
			exhibitionRepository,
		)

	listExhibitionsUseCase :=
		usecase.NewListExhibitionsUseCase(
			exhibitionRepository,
		)

	// ----------------------------------------
	// Handler
	// ----------------------------------------

	server :=
		handler.NewExhibitionHandler(
			createExhibitionUseCase,
			getExhibitionUseCase,
			listExhibitionsUseCase,
		)

	exhibitRepository := exhibitinfra.NewPostgresRepository(postgresPool)
	apiHandler := handler.NewAPIHandler(server, handler.NewExhibitHandler(
		usecase.NewCreateExhibitUseCase(exhibitRepository, idGenerator),
		usecase.NewListExhibitsUseCase(exhibitRepository, exhibitionRepository),
	))

	// OpenAPI Strict Server。
	strictHandler :=
		generated.NewStrictHandler(
			apiHandler,
			nil,
		)

	// OpenAPIによって生成されたHTTP Router。
	httpHandler :=
		generated.Handler(
			strictHandler,
		)

	// ----------------------------------------
	// Middleware
	// ----------------------------------------

	// BrowserからのFrontend → Backend通信を許可する。
	rootHandler :=
		middleware.CORS(
			httpHandler,
		)

	// ----------------------------------------
	// Server
	// ----------------------------------------

	log.Println(
		"server started on :8080",
	)

	if err := http.ListenAndServe(
		":8080",
		rootHandler,
	); err != nil {
		log.Fatal(err)
	}
}
