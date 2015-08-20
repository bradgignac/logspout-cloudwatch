default: test

test:
	go test . -short

integration:
	go test . -v
