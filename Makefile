.PHONY: all setup test updategolden lint format build

setup:
	@echo "Setting up dependencies..."
	go mod tidy

test:
	@echo "Running tests..."
	go test -v ./... -count=1

updategolden:
	@echo "Update golden file..."
	go test -v ./... -count=1 -update

lint:
	@echo "Applying linter..."
	golangci-lint run

format:
	@echo "Formatting code with gofumpt..."
	gofumpt -l -w .

build:
	@echo "Building exec file..."
	goreleaser build --single-target --snapshot --clean --output ./jqcompletion

.PHONY: testuijsonview
testuijsonview:
	DEBUGLOG=1 go run ./uitests/jsonview/uitest.go

.PHONY: testuiqueryview
testuiqueryview:
	DEBUGLOG=1 go run ./uitests/queryview/uitest.go
