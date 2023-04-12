package main

import (
	"net"
	"net/rpc"
	"testing"
)

type mockStockService struct{}

func (s *mockStockService) GetCount(_ struct{}, count *int) error {
	*count = 10
	return nil
}

func (s *mockStockService) Release(_ *Product, released *bool) error {
	*released = true
	return nil
}

func (s *mockStockService) Reserve(_ *Product, reserved *bool) error {
	*reserved = true
	return nil
}

func TestClient(t *testing.T) {
	// Register mock implementation
	service := &mockStockService{}
	rpc.Register(service)

	// Start RPC server on random port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to listen on port: %v", err)
	}
	defer listener.Close()

	// Start server in goroutine
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Errorf("Error accepting connection: %v", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	// Dial client to mock server
	client, err := rpc.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer client.Close()

	// Call methods on mock server and check results
	var count int
	if err := client.Call("StockService.GetCount", struct{}{}, &count); err != nil {
		t.Errorf("Failed to call GetCount: %v", err)
	}
	if count != 10 {
		t.Errorf("Unexpected count value: got %d, want %d", count, 10)
	}

	var releaseResult bool
	if err := client.Call("StockService.Release", &Product{Name: "Product A", Size: 10, UniqCode: 111, Count: 2}, &releaseResult); err != nil {
		t.Errorf("Failed to call Release: %v", err)
	}
	if !releaseResult {
		t.Errorf("Unexpected release result: got %t, want %t", releaseResult, true)
	}

	var reserveResult bool
	if err := client.Call("StockReserve", &Product{Name: "Product A", Size: 10, UniqCode: 111, Count: 2}, &reserveResult); err != nil {
		t.Errorf("Failed to call Reserve: %v", err)
	}
	if !reserveResult {
		t.Errorf("Unexpected reserve result: got %t, want %t", reserveResult, true)
	}
}
