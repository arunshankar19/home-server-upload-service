serve:
	@docker compose up --build -d upload-service

migrate-up:
	@docker compose run --rm migrate up

migrate-down:
	@docker compose run --rm migrate down $(COUNT)

migrate-create:
	@docker compose run --rm migrate create -ext sql -dir /migrations -seq=false $(NAME)

app-logs:
	@docker compose logs --tail 100 --follow upload-service
