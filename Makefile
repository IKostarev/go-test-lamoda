.PHONY: all build clean db migrate test

all: build

build:
  go build -o go-test-lamoda cmd/server/server.go

clean:
  rm -f go-test-lamoda

db:
  docker-compose up -d db

migrate:
  docker-compose run --rm migrate

test:
  go test ./...