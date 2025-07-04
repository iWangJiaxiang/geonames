all: run

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf geonames

.PHONY: build
build: clean tidy
	@echo "Building..."
	@go build -o geonames main.go

.PHONY: run
run: build
	@echo "Running..."
	@./geonames

.PHONY: tidy
tidy:
	@echo "Tidying up..."
	@go mod tidy
