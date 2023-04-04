package handlers

import (
	"encoding/json"
	"go-test-lamoda/internal/database"
	"net/http"
)

type ReleaseRequest struct {
	UniqCodes []int `json:"uniq_codes"`
}

func ReleaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Неверный метод, необходим метод POST", http.StatusMethodNotAllowed)
		return
	}

	// получаем список уникальных кодов товаров для освобождения
	var codes ReleaseRequest
	err := json.NewDecoder(r.Body).Decode(&codes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	// обновляем значения поля "feature" в таблице "stocks" на "available" для указанных товаров
	for _, code := range codes.UniqCodes {
		_, err = tx.Exec("UPDATE stocks SET feature='available' WHERE feature='reserved' AND id=(SELECT stocks_id FROM stocks_products WHERE products_id=(SELECT id FROM products WHERE uniq_code=$1))", code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
