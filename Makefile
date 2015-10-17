all: clean
	@echo "-> Building..."
	@go build
	@echo "-> Running..."
	@./c8vm

test:
	@echo "-> Testing..."
	@go test ./...

clean:
	@echo "-> Cleaning..."
	@rm -f ./c8vm

.PHONY: all clean
