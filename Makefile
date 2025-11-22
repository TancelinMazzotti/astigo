.PHONY: swagger build up down restart rebuild build-debug up-debug down-debug restart-debug rebuild-debug

# Dev
swagger:
	swag init -g .\main.go -o .\docs\

# Standard
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down -v

restart: down up

rebuild: down build up

# Debug
build-debug:
	docker-compose -f docker-compose.debug.yml build

up-debug:
	docker-compose -f docker-compose.debug.yml up -d

down-debug:
	docker-compose -f docker-compose.debug.yml down -v

restart-debug: down up

rebuild-debug: down-debug build-debug up-debug
