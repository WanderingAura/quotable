run:
	go run ./cmd/web
psql:
	psql ${QUOTABLE_DB_DSN}
up:
	@echo 'Running up migrations...'
	migration -path ./migrations -database ${QUOTABLE_DB_DSN} up	