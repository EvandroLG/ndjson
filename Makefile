.PHONY: test tidy fmt

test:
	go test -v ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...
