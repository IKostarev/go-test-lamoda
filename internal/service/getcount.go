package service

func (s *StockService) GetCount(_ struct{}, count *int) error {
	row := s.DB.QueryRow(`SELECT COUNT(*) FROM products
					JOIN stocks_products ON stocks_products.products_id = products.id
					JOIN stocks ON stocks_products.stocks_id = stocks.id
					WHERE stocks.feature = 'available'`)
	if err := row.Scan(count); err != nil {
		return err
	}
	return nil
}
