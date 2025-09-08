.PHONY: build up down rebuild build-debug up-debug down-debug rebuild-debug

# Standard
build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down -v

rebuild: down build up

# Debug
build-debug:
	docker-compose -f docker-compose.debug.yml build

up-debug:
	docker-compose -f docker-compose.debug.yml up -d

down-debug:
	docker-compose -f docker-compose.debug.yml down -v

rebuild-debug: down-debug build-debug up-debug
