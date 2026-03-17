up:
	docker compose up --build

down:
	docker compose down

reset:
	docker compose down -v

migrate-up:
	docker compose run --rm --profile tools migrate up

migrate-down:
	docker compose run --rm --profile tools migrate down 1

migrate-version:
	docker compose run --rm --profile tools migrate version