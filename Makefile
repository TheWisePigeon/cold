build:
	@go build -o cold .

install-deps:
	@go mod download

run:
	@./cold
