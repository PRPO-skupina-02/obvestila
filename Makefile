docs:
	swag fmt
	swag init --parseDependency -g api/api.go -o api/docs

test:
	godotenv go test ./...

test-update:
	godotenv go test ./... -update

install-cli-tools:
	go install github.com/joho/godotenv/cmd/godotenv@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
