build: 
	go build -o bin/octgopus

run: build
	./bin/octgopus

test:
	go test -v ./...