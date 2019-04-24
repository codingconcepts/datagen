.PHONY: test

example:
	go run main.go -script ./examples/script.sql --driver postgres --conn postgres://root@localhost:26257/sandbox?sslmode=disable

test:
	go test ./... -v