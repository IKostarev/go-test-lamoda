package main

import (
	"log"
	"net/rpc/jsonrpc"
)

type Product struct {
	Name     string `json:"name"`
	Size     int    `json:"size"`
	UniqCode int    `json:"uniq_code"`
	Count    int    `json:"count"`
}

func main() {
	client, err := jsonrpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dial error: ", err)
	}
	defer client.Close()

	var count int
	if err := client.Call("StockService.GetCount", struct{}{}, &count); err != nil {
		log.Fatal("GetCount error: ", err)
	}
	log.Printf("Count: %d\n", count)

	var releaseResult bool
	if err := client.Call("StockService.Release", &Product{Name: "Product A", Size: 10, UniqCode: 111, Count: 2}, &releaseResult); err != nil {
		log.Fatal("Release error: ", err)
	}
	log.Printf("Release result: %t\n", releaseResult)

	var reserveResult bool
	if err := client.Call("StockReserve", &Product{Name: "Product A", Size: 10, UniqCode: 111, Count: 2}, &reserveResult); err != nil {
		log.Fatal("Reserve error: ", err)
	}
	log.Printf("Reserve result: %t\n", reserveResult)
}
