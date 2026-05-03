XXXX ?= 8080

.PHONY: run stop clean test

run:
	@echo "Starting application on port $(XXXX)..."
	XXXX=$(XXXX) docker-compose up --build

stop:
	docker-compose down

clean:
	docker-compose down -v
	rm -f stock-server

test:
	go test -v ./...