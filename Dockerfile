FROM golang:1.20

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o go-test-lamoda ./cmd/stock/main.go

CMD ["./go-test-lamoda"]