qa: analyze test

analyze:
	@go vet
	@go tool staticcheck --checks=all

test:
	@go test -failfast -cover ./...


.PHONY: qa \
		analyze \
		test