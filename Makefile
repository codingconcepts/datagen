.PHONY: test

cockroachdb:
	docker run -d -p 26257:26257 -p 8080:8080 cockroachdb/cockroach:v2.1.6 start --insecure

example:
	go run main.go -script ./examples/script.sql --driver postgres --conn postgres://root@localhost:26257/sandbox?sslmode=disable

test:
	go test ./... -v ;\
	go test ./... -cover

bench:
	go test ./... -bench=.

cover:
	go test ./... -coverprofile=coverage.out -coverpkg=\
	github.com/codingconcepts/datagen/internal/pkg/parse,\
	github.com/codingconcepts/datagen/internal/pkg/random,\
	github.com/codingconcepts/datagen/internal/pkg/runner;\
	go tool cover -html=coverage.out

release:
	# linux
	GOOS=linux go build -ldflags "-X main.semver=${VERSION}" -o datagen ;\
	tar -zcvf datagen_${VERSION}_linux.tar.gz ./datagen ;\

	# macos
	GOOS=darwin go build -ldflags "-X main.semver=${VERSION}" -o datagen ;\
	tar -zcvf datagen_${VERSION}_macOS.tar.gz ./datagen ;\

	# windows
	GOOS=windows go build -ldflags "-X main.semver=${VERSION}" -o datagen ;\
	tar -zcvf datagen_${VERSION}_windows.tar.gz ./datagen ;\