.PHONY: fmt lint

fmt:
	gofmt -w .
	goimports -local github.com/ChefBingbong/viem-go -w .

lint: fmt
	golangci-lint run