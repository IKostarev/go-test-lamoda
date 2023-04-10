package handlers

import (
	"fmt"
	"go-test-lamoda/pkg/database"
	"io"
	"net/http"
)

func GetCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Неверный метод, необходим метод GET", http.StatusMethodNotAllowed)
		return
	}

	// начинаем транзакцию
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // откатываем транзакцию в случае ошибки

	// блокируем таблицу "products", чтобы другие клиенты не могли изменять данные
	_, err = tx.Exec("LOCK TABLE products IN SHARE MODE")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// выполняем запрос на подсчет количества товаров
	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM products 
    						JOIN stocks_products ON stocks_products.products_id = products.id 
    						JOIN stocks ON stocks_products.stocks_id = stocks.id
                			WHERE stocks.feature = 'available'`).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// выводим результат
	io.WriteString(w, fmt.Sprintf("Количество товаров на складе: %d", count))
}
