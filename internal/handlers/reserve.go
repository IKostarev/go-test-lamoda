package handlers

import (
	"encoding/json"
	"go-test-lamoda/internal/database"
	"net/http"
)

type ReserveRequest struct {
	UniqCodes []int `json:"uniq_codes"`
}

func ReserveToStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Неверный метод, необходим метод POST", http.StatusMethodNotAllowed)
		return
	}

	var req ReserveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.UniqCodes) == 0 {
		http.Error(w, "Передано 0 значений", http.StatusBadRequest)
		return
	}

	// начало транзакции
	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// обновление записей в таблице stocks
	stmt, err := tx.Prepare("UPDATE stocks SET feature = 'reserved' WHERE id = (SELECT stocks_id FROM stocks_products WHERE products_id = $1 AND feature = 'available' LIMIT 1 FOR UPDATE) RETURNING id")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	for _, id := range req.UniqCodes {
		_, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// подтверждение транзакции
	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
