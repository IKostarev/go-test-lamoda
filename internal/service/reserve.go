package service

import (
	"database/sql"
	"log"
)

func (s *StockService) ReserveToStock(product *Product, result *bool) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != sql.ErrTxDone && err != nil {
			log.Println("Transaction rollback error:", err)
		}
	}()

	row := tx.QueryRow("SELECT count FROM products WHERE name=$1 AND size=$2 AND uniqcode=$3", product.Name, product.Size, product.UniqCode)
	var count int
	if err := row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			count = 0
		} else {
			return err
		}
	}

	if _, err := tx.Exec(`UPDATE stocks SET feature = 'reserved' WHERE id = (SELECT stocks_id FROM stocks_products WHERE products_id = $1 AND feature = 'available' LIMIT 1 FOR UPDATE) RETURNING id`,
		product.Name, product.Size, product.UniqCode, product.Count); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	*result = true
	return nil
}
