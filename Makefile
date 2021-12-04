test:
	@go test -race ./...

test-unit:
	@go test -race ./internal/vwap

test-integration:
	@go test -race ./internal/websocket

upgrade:
	@echo "Upgrading dependencies..."
	@go get -u
	@go mod tidy
	
run:
	@go run main.go

build:
	@go build -o vwap main.go

clean:
	@rm -rf wvap