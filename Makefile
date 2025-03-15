.PHONY: all setup test lint format

setup:
	@echo "Setting up dependencies..."
	go mod tidy

test:
	@echo "Running tests..."
	go test -v ./... -count=1

lint:
	@echo "Applying linter..."
	golangci-lint run

format:
	@echo "Formatting code with gofumpt..."
	gofumpt -l -w .

.PHONY: testuijsonview
testuijsonview:
	DEBUGLOG=1 go run ./uitests/jsonview/uitest.go

.PHONY: testuiqueryview
testuiqueryview:
	DEBUGLOG=1 go run ./uitests/queryview/uitest.go
