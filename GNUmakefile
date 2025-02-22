.PHONY: example fmt

example:
	go run example/main.go

fmt:
	go mod tidy
	go vet
	go fmt ./...
