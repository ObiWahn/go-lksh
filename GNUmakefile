.PHONY: example fmt

example:
	go run example/main.go

fmt:
	go fmt ./...
	go mod tidy
	go vet
