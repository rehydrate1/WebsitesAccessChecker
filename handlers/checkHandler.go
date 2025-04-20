package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"io"
	"sync"
	"github.com/rehydrate1/WebsitesAccessChecker/internal/checker"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получен запрос к /check")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела запроса: %v", err)
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var urls []string
	err = json.Unmarshal(bodyBytes, &urls)
	if err != nil {
		log.Printf("Ошибка парсинга JSON: %v. Тело: %s", err, string(bodyBytes))
		http.Error(w, "Неверный формат JSON. Ожидается массив URL строк.", http.StatusBadRequest)
		return
	}

	if len(urls) == 0 {
		http.Error(w, "Список URL для проверки пуст.", http.StatusBadRequest)
		return
	}
	log.Printf("Получено %d URL для проверки", len(urls))

	var wg sync.WaitGroup
	resultChan := make(chan checker.CheckResult, len(urls))

	for _, url := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()
			
			result := checker.CheckSite(u)
			resultChan <- result
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	finalResults := make([]checker.CheckResult, 0, len(urls))
	for result := range resultChan {
		finalResults = append(finalResults, result)
	}
	log.Printf("Проверка завершена. Собрано %d результатов.", len(finalResults))

	responseBytes, err := json.Marshal(finalResults)
	if err != nil {
		log.Printf("Ошибка маршалинга результатов в JSON: %v", err)
		http.Error(w, "Не удалось подготовить ответ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Printf("Ошибка отправки ответа клиенту: %v", err)
	}
}