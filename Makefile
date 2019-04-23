.PHONY: test

run_cockroach:
	go run main.go -script ./examples/cockroach/script.sql --driver postgres --conn postgres://root@localhost:26257/sandbox?sslmode=disable

run_mysql:
	go run main.go -script ./examples/mysql/script.sql --driver mysql --conn un:pw@/sandbox

test:
	go test ./... -v