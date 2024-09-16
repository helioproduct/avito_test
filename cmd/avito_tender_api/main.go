package main

import (
	"avito_api/internal/config"

	"avito_api/pkg/postgres"

	"log"

	uc "avito_api/internal/usecase"
	advancedBidPg "avito_api/internal/usecase/repo/advancedBid/postgres"
	bidPg "avito_api/internal/usecase/repo/bid/postgres"
	orgPg "avito_api/internal/usecase/repo/organization/postgres"
	tenderPg "avito_api/internal/usecase/repo/tender/postgres"
	userPg "avito_api/internal/usecase/repo/user/postgres"

	bidRoutes "avito_api/internal/server/handlers/bids"
	orgRoutes "avito_api/internal/server/handlers/organizations"
	pingRoutes "avito_api/internal/server/handlers/ping"
	tenderRoutes "avito_api/internal/server/handlers/tenders"
	userRoutes "avito_api/internal/server/handlers/users"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"net/http"
)

func main() {

	cfg := config.MustLoad()

	log.Println("Loaded configuration:", cfg)
	log.Println("Postgres connection URL:", cfg.ConnURL)
	log.Println("Server address:", cfg.Address)

	db, err := postgres.NewPostgresDB(cfg.ConnURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Could not establish connection to the database")
	}
	log.Println("Database connection successful")

	// repositories
	userRepo := userPg.NewPostgresUserRepo(db)
	orgRepo := orgPg.NewPostgresOrganizationRepo(db)
	tenderRepo := tenderPg.NewPostgresTenderRepo(db)
	bidRepo := bidPg.NewPostgresBidRepo(db)

	bidAdvancedRepo := advancedBidPg.NewPostgresBidAdvancedRepo(db, bidRepo)

	// usecases
	orgUC := uc.NewOrganizationUseCase(orgRepo)
	userUC := uc.NewUserUseCase(userRepo)
	tenderUC := uc.NewTenderUseCase(tenderRepo, orgRepo, userRepo)
	bidUC := uc.NewBidUseCase(bidAdvancedRepo, orgRepo, userRepo, tenderRepo)

	advancedBidUC := uc.NewAdvancedBidUseCase(bidUC, bidAdvancedRepo, orgRepo, tenderRepo, 3)

	// routes
	pingHandlers := pingRoutes.NewHandlers()
	userHandlers := userRoutes.NewHandlers(userUC)
	orgHandlers := orgRoutes.NewHandlers(orgUC)
	tenderHandlers := tenderRoutes.NewHandlers(tenderUC)
	// bidHandlers := bidRoutes.NewHandlers(bidUC)
	bidHandlers := bidRoutes.NewHandlers(advancedBidUC)

	router := chi.NewRouter()

	router.Route("/api/ping", func(r chi.Router) {
		pingRoutes.RegisterRoutes(r, pingHandlers)
	})

	router.Route("/api/users", func(r chi.Router) {
		userRoutes.RegisterRoutes(r, userHandlers)
	})

	router.Route("/api/organizations", func(r chi.Router) {
		orgRoutes.RegisterRoutes(r, orgHandlers)
	})

	router.Route("/api/tenders", func(r chi.Router) {
		tenderRoutes.RegisterRoutes(r, tenderHandlers)
	})

	router.Route("/api/bids", func(r chi.Router) {
		bidRoutes.RegisterRoutes(r, bidHandlers)
	})

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	log.Println("Starting server on", cfg.Address)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
