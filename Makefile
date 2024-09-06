
include .env
export


DATABASE_URL ?= postgresql://${POSTGRES_APP_USER}:${POSTGRES_APP_PASS}@localhost:5432/${POSTGRES_APP_DB}
#######################################################################################################################

DOCKER_COMPOSE_FILES ?= $(shell find docker -maxdepth 1 -type f -name "*.yaml" -exec printf -- '-f %s ' {} +; echo)

#######################################################################################################################

## ▸▸▸ Development commands ◂◂◂

.PHONY: help
help:			## Show this help
	@fgrep -h "## " $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/## //'

.PHONY: fmt
fmt:			## Format code
	find . -iname "*.go" | xargs gofmt -w

# .PHONY: build
# build:			## Build apps
# 	CGO_ENABLED=0 go build \
# 		-mod=mod \
# 		-tags='no_mysql no_sqlite3' \
# 		-o ./bin/goods$(shell go env GOEXE) cmd/goods/main.go

# 	CGO_ENABLED=0 go build \
# 		-mod=mod \
# 		-tags='no_mysql no_sqlite3' \
# 		-o ./bin/migrate$(shell go env GOEXE) cmd/migrate/main.go

# 	CGO_ENABLED=0 go build \
# 		-mod=mod \
# 		-tags='no_mysql no_sqlite3' \
# 		-o ./bin/connect$(shell go env GOEXE) cmd/connect/main.go

.PHONY: lint
lint:
	go vet ./...
	staticcheck ./...
	errcheck ./...
	golint ./...

.PHONY: test
test:
	go test ./... -p 1

.PHONY: testcov
testcov:
	# go test -v -coverpkg=./internal/... -coverprofile=profile.cov ./internal/...
	go test ./internal/... -coverpkg=./internal/... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func coverage.out

.PHONY: bufgen
bufgen:
	buf generate
	
.PHONY: mockgen
mockgen:
	mockgen -source=./internal/server/repository/user.go -destination=./mocks/mock_user.go -package=mocks
	mockgen -source=./internal/server/repository/repository.go -destination=./mocks/mock_repository.go -package=mocks
	mockgen -source=./internal/server/repository/meta.go -destination=./mocks/mock_meta.go -package=mocks
	mockgen -source=./tools/client/minio_client.go -destination=./mocks/minio_client.go -package=mocks
	mockgen -source=./internal/app/client/client.go -destination=./mocks/mock_app_client.go -package=mocks
	mockgen -source=./internal/client/client.go -destination=./mocks/mock_internal_client.go -package=mocks
	# mockgen -source=./gen/proto/api/service/v1/service_grpc.pb.go -destination=./mocks/mock_dataservice.go -package=mocks
	# mockgen -source=./gen/proto/api/user/v1/user_grpc.pb.go -destination=./mocks/mock_userservice.go -package=mocks


.PHONY: clean
clean:			## Remove generated artifacts
	@rm -rf ./bin
	@rm -rf ./docker/volume

#######################################################################################################################

## ▸▸▸ Docker commands ◂◂◂

.PHONY: config
config:			## Show Docker config
	docker compose ${DOCKER_COMPOSE_FILES} config

.PHONY: up
up:			## Run Docker services
	docker compose ${DOCKER_COMPOSE_FILES} up --detach

.PHONY: build
build:			## Build Docker services
	docker compose ${DOCKER_COMPOSE_FILES} build

.PHONY: down
down:			## Stop Docker services
	docker compose ${DOCKER_COMPOSE_FILES} down

.PHONY: ps
ps:			## Show Docker containers info
	docker ps --size --all --filter "name=datakeeper"

#######################################################################################################################

## ▸▸▸ Utils commands ◂◂◂

.PHONY: connect
connect:		## Connect to the database
	pgcli ${DATABASE_URL}

.PHONY: goose-status
goose-status:		## Dump the migration status for the current DB
	goose -dir migrations postgres ${DATABASE_URL} status

.PHONY: goose-up
goose-up:		## Migrate the DB to the most recent version available
	goose -dir migrations postgres ${DATABASE_URL} up

.PHONY: goose-down
goose-down:		## Roll back the version by 1
	goose -dir migrations postgres ${DATABASE_URL} down

#######################################################################################################################

.reqs:
	pip install --upgrade pip
	pip install pgcli
	pip install "psycopg[binary,pool]"
	go install github.com/pressly/goose/v3/cmd/goose@latest




