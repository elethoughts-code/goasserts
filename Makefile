mocks:
	mockgen -destination="mocks/assertion/assert.go" -package=mocks github.com/elethoughts-code/goasserts/assertion PublicTB

lint:
	golangci-lint run -c .golangci.yml ./...