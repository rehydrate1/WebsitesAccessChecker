package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rehydrate1/WebsitesAccessChecker/internal/hostfile"
)

func GetHostsHandler(w http.ResponseWriter, r *http.Request) {
	domains, err := hostfile.GetDomainsFromHosts()
	if err != nil {
		log.Printf("Ошибка при чтении файла hosts: %v\n", err)
		http.Error(w, "Ошибка при чтении файла hosts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(domains)
	if err != nil {
		log.Printf("Ошибка преобразования списка доменов в JSON: %v\n", err)
		http.Error(w, "Внутренняя ошибка сервера при форматировании данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Printf("Ошибка при отправке JSON ответа: %v\n", err)
	}
}
