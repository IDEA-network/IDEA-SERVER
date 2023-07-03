.PHONY: dev build

dev:
	docker-compose up dev
build:
	docker-compose up --build -d app 