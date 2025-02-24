test: 
	go test -v -race ./...

update: 
	go get -u -t ./...

tidy: 
	go mod tidy

install-lint:
	mkdir -p $(ROOT)/bin
	echo "*"> $(ROOT)/bin/.gitignore
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(ROOT)/bin v1.63.4

lint: 
	# bin/golangci-lint run ./...
	golangci-lint run ./...