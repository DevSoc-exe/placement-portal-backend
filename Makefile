BINARY_DIR=bin
BINARY_NAME=main

build:
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) cmd/main.go

run: build
	@./$(BINARY_DIR)/$(BINARY_NAME)

clean:
	@rm -rf $(BINARY_DIR)

rebuild: clean build run

production_build:
	@mkdir -p $(BINARY_DIR)
	@go build -tags netgo -ldflags '-s -w' -o $(BINARY_DIR)/$(BINARY_NAME) cmd/main.go
	@echo "Deployed $(BINARY_NAME) successfully"

deploy: 
	@./$(BINARY_DIR)/$(BINARY_NAME)
