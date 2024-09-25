
all: test

test:
	go test -v ./...

N := 20000000

bench:
	go run cmd/gen/main.go -n $(N) | go run cmd/sort/main.go > /dev/null
