include .env
export

compose-up: ### Run docker-compose
	docker-compose up --build -d
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### remove docker volume
	docker volume rm pg-data
.PHONY: docker-rm-volume

swag: ### generate swagger docs
	swag init -g internal/app/app.go --parseInternal --parseDependency
