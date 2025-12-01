.PHONY: all build run test clean

APP_NAME=airbridge

all: build

build:
	go build -o $(APP_NAME) main.go

# Allow passing arguments to `make run`
# Example: make run send
ARGS = $(filter-out $@,$(MAKECMDGOALS))

run:
	go run main.go $(ARGS)

test:
	go test -v ./...

clean:
	rm -f $(APP_NAME)

# Catch-all target to allow arguments like "send", "receive" to be passed to "run" without error
%:
	@:
