up:
	docker compose up -d

down:
	docker compose down

reset:
	docker compose down -v

migrate-up:
	docker compose --profile tools run --rm migrate up

migrate-down:
	docker compose run --rm --profile tools migrate down 1

migrate-version:
	docker compose run --rm --profile tools migrate version