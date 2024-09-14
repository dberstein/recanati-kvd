COVFILE := x.cov

test:
	@go test -v -coverprofile $(COVFILE)  ./... \
	&& go tool cover -func=$(COVFILE)

run:
	@go run main.go -f 15s