export XXXX ?= 8080

.PHONY: run stop clean test

run:
	@echo "Starting application on port $(XXXX)..."
	docker-compose up --build

stop:
	docker-compose down

clean:
	docker-compose down -v
	rm -f stock-server

test:
	@echo "Running tests inside Docker..."
	docker run --rm -v "$(PWD):/app" -w /app golang:1.26-alpine go test -v ./...