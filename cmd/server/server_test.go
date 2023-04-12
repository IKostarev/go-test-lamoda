package main

import (
	"database/sql"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"testing"
	"time"
)

func setupTest() (*StockService, net.Listener, error) {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/stockdb?sslmode=disable")
	if err != nil {
		return nil, nil, err
	}
	service := &StockService{db: db}
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go jsonrpc.ServeConn(conn)
		}
	}()

	return service, listener, nil
}

func teardownTest(service *StockService, listener net.Listener) {
	rpc.Register(service)
	service.db.Close()
	listener.Close()
}

func TestGetCount(t *testing.T) {
	service, listener, err := setupTest()
	if err != nil {
		t.Fatal("Failed to setup test:", err)
	}
	defer teardownTest(service, listener)

	// Wait for listener to start accepting connections
	time.Sleep(time.Millisecond * 50)

	client, err := jsonrpc.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal("Failed to connect to server:", err)
	}
	defer client.Close()

	var count int
	err = client.Call("StockService.GetCount", struct{}{}, &count)
	if err != nil {
		t.Fatal("GetCount failed:", err)
	}

	// Verify the result
	if count < 0 {
		t.Errorf("Expected count to be non-negative, but got %d", count)
	}
}

func TestMain(m *testing.M) {
	// Make sure we can connect to the database before running tests
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/stockdb?sslmode=disable")
	if err != nil {
		log.Fatal("Open Postgres error is: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	// Run tests
	code := m.Run()

	os.Exit(code)
}
