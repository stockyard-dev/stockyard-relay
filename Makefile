build:
	CGO_ENABLED=0 go build -o relay ./cmd/relay/

run: build
	./relay

test:
	go test ./...

clean:
	rm -f relay

.PHONY: build run test clean
