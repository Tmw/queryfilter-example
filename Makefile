
.PHONY: server call test

server:
	@go run main.go

call:
	@curl -s "localhost:8080/shirts?size=L&min_price=9&max_price=10&color=yellow"

test:
	@go test ./... -v
