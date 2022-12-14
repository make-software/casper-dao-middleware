package main

import (
	"log"
	"net/http"
	"os"

	"casper-dao-middleware/apps/api/config"
	"casper-dao-middleware/apps/api/handlers"
	"casper-dao-middleware/apps/api/swagger"
	"casper-dao-middleware/internal/crdao/dao_event_parser"
	"casper-dao-middleware/internal/crdao/persistence"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	cfg *config.Env,
	entityManager persistence.EntityManager,
	daoContractPackageHashes dao_event_parser.DAOContractsMetadata,
) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healcheck - OK"))
	})

	reputationHandler := handlers.NewReputation(entityManager, daoContractPackageHashes)
	votingHandler := handlers.NewVoting(entityManager)

	router.Get("/accounts/{address}/total-reputation", reputationHandler.HandleGetTotalReputation)
	router.Get("/accounts/{address}/aggregated-reputation-changes", reputationHandler.HandleGetAggregatedReputationChange)
	router.Get("/accounts/{address}/votes", votingHandler.HandleGetAccountVotes)

	router.Get("/votings", votingHandler.HandleGetVotings)
	router.Get("/votings/{voting_id}/votes", votingHandler.HandleGetVotingVotes)

	swaggerHost := string(cfg.Addr)
	if envHost := os.Getenv("SWAGGER_HOST"); envHost != "" {
		swaggerHost = envHost
	}

	swagger.SwaggerInfo.Host = swaggerHost
	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))

	log.Printf("Swagger is available on: %s/swagger/index.html#/\n", swagger.SwaggerInfo.Host)

	return router
}
