
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

lint:
	go vet ./...
	staticcheck ./...
	errcheck ./...
	golint ./...

test:
	go test ./... -p 1
