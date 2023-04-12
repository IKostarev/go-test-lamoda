package service

import "database/sql"

type Product struct {
	Name     string `json:"name"`
	Size     int    `json:"size"`
	UniqCode int    `json:"uniq_code"`
	Count    int    `json:"count"`
}

type StockService struct {
	Count int
	DB    *sql.DB
}
