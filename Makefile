build: 
	go build -o cmd/api/api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./...