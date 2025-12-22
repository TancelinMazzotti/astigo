.PHONY: swagger build up down restart rebuild build-debug up-debug down-debug restart-debug rebuild-debug up-monitoring down-monitoring

# Dev
swagger:
	swag init -g .\main.go -o .\docs\

# Standard
build:
	docker-compose --profile astigo build

up:
	docker-compose --profile astigo up -d

down:
	docker-compose --profile astigo down -v

restart: down up

rebuild: down build up

# Debug
build-debug:
	docker-compose --profile debug build

up-debug:
	docker-compose --profile debug up -d

down-debug:
	docker-compose --profile debug down -v

restart-debug: down up

rebuild-debug: down-debug build-debug up-debug

up-monitoring:
	docker-compose --profile monitoring up -d

down-monitoring:
	docker-compose --profile monitoring down -v
