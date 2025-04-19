package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"io"
	"sync"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос к /check")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается. Используйте POST", http.StatusMethodNotAllowed)
		return
	}
}