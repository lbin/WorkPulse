run:
	ENV=dev go run ./cmd/server

fmt:
	go fmt ./...

test:
	go test ./...
