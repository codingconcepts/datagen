.PHONY: test

example:
	go run main.go -script ./examples/script.sql --driver postgres --conn postgres://root@localhost:26257/sandbox?sslmode=disable

test:
	go test ./... -v ;\
	go test ./... -cover

bench:
	go test ./... -bench=.

cover:
	go test ./... -coverprofile=coverage.out -coverpkg=\
	github.com/ejdem86/datagen/internal/pkg/parse,\
	github.com/ejdem86/datagen/internal/pkg/random,\
	github.com/ejdem86/datagen/internal/pkg/runner;\
	go tool cover -html=coverage.out
