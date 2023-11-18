lint:
	golangci-lint run

run:
	go run . --server

build:
	go build -o bin/mindia

docker:
	docker build -t mindia:latest .

release-check:
	goreleaser check

release-local:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean