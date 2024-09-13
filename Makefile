COVFILE := x.cov

test:
	@go test -coverprofile $(COVFILE)  ./... \
	&& go tool cover -func=$(COVFILE)

run:
	@go run main.go -f 15s