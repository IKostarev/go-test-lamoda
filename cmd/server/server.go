package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Product struct {
	Name     string `json:"name"`
	Size     int    `json:"size"`
	UniqCode int    `json:"uniq_code"`
	Count    int    `json:"count"`
}

type StockService struct {
	count int
	stock map[Product]int
}

func (s *StockService) GetCount(_ *struct{}, count *int) error {
	*count = s.count
	return nil
}

func (s *StockService) Release(product *Product, result *bool) error {
	if count, ok := s.stock[*product]; ok && count >= product.Count {
		s.stock[*product] -= product.Count
		*result = true
		return nil
	}

	*result = false
	return errors.New("error")
}

func (s *StockService) Reverse(product *Product, result *bool) error {
	if _, ok := s.stock[*product]; ok {
		s.stock[*product] += product.Count
	} else {
		s.stock[*product] = product.Count
	}

	*result = true
	return nil
}

func main() {
	stockService := &StockService{
		count: 10,
		stock: make(map[Product]int),
	}

	server := rpc.NewServer()
	if err := server.Register(stockService); err != nil {
		log.Fatal("Register error: ", err)
	}

	http.Handle("/jsonrpc", server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	defer listener.Close()

	log.Println("Starting server on :8080")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		go server.ServeConn(jsonrpc.NewServerCodec(conn))
	}
}
