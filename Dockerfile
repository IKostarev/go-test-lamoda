FROM golang:latest

ENV GOPATH=/

RUN apt-get update && apt-get install -y \
    postgresql-client \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

COPY ./ ./

RUN go mod download
RUN go build -o go-test-lamoda ./cmd/stock/main.go

CMD ["./go-test-lamoda"]
