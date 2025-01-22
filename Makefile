build:
	@go build -o bin/groblin.exe cmd/groblin/main.go

dev:
	@go run cmd/groblin/main.go $(ARGS)

start:
	@./bin/groblin

clean:
	@go clean -cache -modcache
