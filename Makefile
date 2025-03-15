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

.PHONY: testjsonview
testjsonview:
	go test ./jsonview --tags debug -v -count=1

.PHONY: testqueryviewonly
testqueryviewonly:
	DEBUGLOG=1 go test ./queryview --tags debug -v -count=1 -run OnlyView

.PHONY: testqueryviewquery
testqueryviewquery:
	DEBUGLOG=1 go test ./queryview --tags debug -v -count=1 -run Query
