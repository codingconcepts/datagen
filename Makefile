.PHONY: test

run:
	go run main.go -script input.sql --driver postgres --conn postgres://root@localhost:26257/sandbox?sslmode=disable

test:
	go test ./... -v