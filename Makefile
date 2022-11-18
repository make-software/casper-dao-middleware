install-deps:
	go mod vendor

setup:
	@$(MAKE) install-deps
	git config --local core.hooksPath .githooks/

run-local-infra:
	docker compose -f infra/local/docker-compose.db.yaml --project-name casper_dao_middleware up -d

stop-local-infra:
	docker compose -f infra/local/docker-compose.db.yaml --project-name casper_dao_middleware down

sync-db:
	sh -ac '. apps/handler/.env; migrate -database "mysql://$$DATABASE_URI" -path internal/crdao/resources/migrations up'

sync-test-db:
	sh -ac '. internal/crdao/tests/.env.test; migrate -database "mysql://$$TEST_DATABASE_URI" -path internal/crdao/resources/migrations up'

swagger:
	cd ./apps/api/ && swag init --parseDependency --output swagger --overridesFile swagger/.swaggo

swagger-format:
	cd ./apps/api/ && swag fmt

help:
	@echo "Usage: make <command>"
	@echo
	@echo "Where <command> is one of:"
	@echo "  install-deps                      Install project dependencies to the vendor directory"
	@echo "  setup                             Install project dependencies for local development"
	@echo "  run-local-infra                   Runs project infra for local development"
	@echo "  stop-local-infra                   Stops project infra for local development and removes containers"
	@echo "  sync-db                           Actualises network store database for local development"
	@echo "  sync-test-db                      Actualises network store database for running tests locally"
	@echo "  swagger                      	    Generate swagger documentation based on comments in api/handlers"
	@echo "  swagger-format                    Run swagger comments formatting"


