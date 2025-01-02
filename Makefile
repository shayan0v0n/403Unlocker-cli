lint:
	golangci-lint run 
build:
	go build -o 403unlocker cmd/403unlockercli/main.go 
test:
	go test ./...
