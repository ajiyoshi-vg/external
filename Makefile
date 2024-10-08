
all: test

test:
	go test -v ./...

N := 20000000

bench.dat:
	go run cmd/gen/main.go -n $(N) > $@

bench: bench.dat
	time cat $< | go run cmd/sort/main.go | go run cmd/sorted/main.go

sort_command: bench.dat
	sort --version
	time cat $< | sort --parallel 10 | go run cmd/sorted/main.go

sorted:
	echo "2\n1" | go run cmd/sorted/main.go
	echo "1\n2" | go run cmd/sorted/main.go
