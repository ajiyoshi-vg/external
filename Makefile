
all: test

test:
	go test -v ./...

N := 20000000

bench.dat:
	go run cmd/gen/main.go -n $(N) > $@

bench: bench.dat
	time cat $< | go run cmd/sort/main.go > /dev/null
