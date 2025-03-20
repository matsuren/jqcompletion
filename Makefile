.PHONY: all setup test lint format build

all: setup test lint format build

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

build:
	@echo "Building exec file..."
	goreleaser build --single-target --snapshot --clean --output ./jqcompletion

# Optional command
.PHONY: demo-gif
demo-gif: build
	@if ! command -v vhs >/dev/null 2>&1; then \
		echo "Error: VHS is not installed. Please install from https://github.com/charmbracelet/vhs"; \
		exit 1; \
	fi
	@echo "Generating demo GIF..."
	vhs .README/demo.tape

.PHONY: testuijsonview
testuijsonview:
	DEBUGLOG=1 go run ./uitests/jsonview/uitest.go

.PHONY: testuiqueryview
testuiqueryview:
	DEBUGLOG=1 go run ./uitests/queryview/uitest.go
