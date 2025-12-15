run:
	ENV=dev go run ./cmd/server

fmt:
	go fmt ./...

test:
	go test ./...

migrate:
	go run ./cmd/migrate
