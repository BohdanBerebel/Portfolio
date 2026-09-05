up:
	docker compose up --build

down:
	docker compose down

migrate-up:
	docker exec -it go_app migrate -path /migrations -database "postgres://admin:test@postgres:5432/portfolio?sslmode=disable" up

migrate-down:
	docker exec -it go_app migrate -path /migrations -database "postgres://admin:test@postgres:5432/portfolio?sslmode=disable" down