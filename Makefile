.PHONY: test tidy fmt lint

test:
	go test -v ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

lint:
	go vet ./...