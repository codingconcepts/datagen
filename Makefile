.PHONY: test

run_cockroach:
	go run main.go -script ./examples/cockroach/script.sql --driver postgres --conn postgres://un:pw@localhost:26257/sandbox?sslmode=disable

run_mysql:
	go run main.go -script ./examples/mysql/script.sql --driver mysql --conn un:pw@/sandbox

test:
	go test ./... -v