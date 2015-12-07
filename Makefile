default: deps test

deps:
	go get -t .

test:
	go test . -short

integration:
	go test . -v -timeout 2h

.PHONY: test integration
