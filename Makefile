CURRENT_DIR=$(shell pwd)

APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd

TAG=latest
ENV_TAG=latest
DOCKERFILE=Dockerfile

run:
	go run cmd/main.go
	
swag_init:
	swag init -g api/router.go  -o api/docs

migrate_up:
	migrate -path migrations/ -database "$(DB_URL)" up

migrate_down:
	migrate -path migrations/ -database "$(DB_URL)" down

migrate_force:
	migrate -path migrations/ -database "$(DB_URL)" force 6

create_migrate:
	./scripts/create_migration.sh

compose_down:  
	docker compose down

compose_up: compose_down
	docker compose up -d --build

crud:
	./scripts/crud.sh

create-repo:
	bash ./scripts/git-lab-hub-repo-creator.sh

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} . -f ${DOCKERFILE}
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

clear-image:
	docker rmi ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker rmi ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

swag-init:
	swag init -g api/api.go -o api/docs --parseDependency

linter:
	golangci-lint run