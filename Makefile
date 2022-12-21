test:
	go test -v -count=1 ./...

test100:
	go test -v -count=100 ./...

race:
	go test -v -race -count=1 ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	# storager
	mockgen -source=internal/encoder/url_storager.go \
	-destination=internal/encoder/mocks/mock_url_storager.go

	# shortener
#	mockgen -source=internal/shortener/shortener.go \
#	-destination=internal/shortener/mocks/mock_shortener.go