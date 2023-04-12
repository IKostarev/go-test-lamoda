package service

import (
	"database/sql"
	"errors"
	"log"
)

func (s *StockService) Release(product *Product, result *bool) error {
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
		return err
	}

	if count >= product.Count {
		if _, err := tx.Exec(`UPDATE stocks 
							SET feature='available' 
							WHERE feature='reserved' AND id=(SELECT stocks_id 
															 FROM stocks_products 
															 WHERE products_id=(SELECT id 
														FROM products WHERE uniq_code=$1))`, product.UniqCode); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		*result = true
		return nil
	}

	*result = false
	return errors.New("insufficient stock")
}
