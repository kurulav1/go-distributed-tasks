SHELL:=/bin/sh

up:
	docker compose up -d nats roach1 roach2 roach3
	docker compose up -d db-bootstrap
	docker compose up -d api worker ui

down:
	docker compose down

reset:
	docker compose down -v --remove-orphans
	docker network prune -f
	docker compose up -d nats roach1 roach2 roach3
	docker compose up -d db-bootstrap
	docker compose up -d api worker ui

logs:
	docker compose logs -f

ps:
	docker compose ps
