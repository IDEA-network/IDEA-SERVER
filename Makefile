.PHONY: dev build down redis

dev:
	docker-compose up -d redis
	docker-compose up dev
build:
	docker-compose up -d redis
	docker-compose up --build -d app 
down:
	docker-compose down --rmi all
redis:
	docker-compose exec redis /bin/sh -c 'redis-cli'
	