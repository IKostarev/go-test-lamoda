package main

import (
	"database/sql"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	_ "github.com/lib/pq"
)

type StockService struct {
	count int
	db    *sql.DB
}

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/stockdb?sslmode=disable")
	if err != nil {
		log.Fatal("Open Postgres error is: ", err)
	}
	defer db.Close()

	service := &StockService{db: db}
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Have error in listen server: ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error in connect Accepted: ", err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
