document.addEventListener('DOMContentLoaded', () => {
    const urlsInput = document.getElementById('urlsInput');
    const checkButton = document.getElementById('checkButton');
    const resultsArea = document.getElementById('resultsArea');
    const loadingIndicator = document.getElementById('loadingIndicator');

    checkButton.addEventListener('click', async () => {
        const urlsText = urlsInput.value.trim();
        const urls = urlsText.split('\n').map(url => url.trim()).filter(url => url !== ''); // Получаем непустые URL

        if (urls.length === 0) {
            resultsArea.innerHTML = '<p>Пожалуйста, введите хотя бы один URL для проверки.</p>';
            return;
        }

        // Очистка предыдущих результатов и показ индикатора
        resultsArea.innerHTML = '';
        loadingIndicator.classList.remove('hidden');
        checkButton.disabled = true; // Блокируем кнопку на время запроса

        try {
            // Отправляем запрос на бэкенд (эндпоинт /check)
            const response = await fetch('/check', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                // Отправляем массив строк в теле запроса в формате JSON
                body: JSON.stringify(urls), 
            });

            if (!response.ok) {
                // Если сервер вернул ошибку (не 2xx статус)
                const errorText = await response.text();
                throw new Error(`Ошибка сервера: ${response.status} ${response.statusText}. ${errorText}`);
            }

            // Ожидаем JSON с результатами от бэкенда
            // Формат ответа ожидается таким:
            // [
            //   { "url": "...", "status_code": 200, "is_available": true, "error": "", "response_time_ms": 123 },
            //   { "url": "...", "status_code": 0, "is_available": false, "error": "connection refused", "response_time_ms": 0 },
            //   ...
            // ]
            const results = await response.json();

            // Отображаем результаты
            displayResults(results);

        } catch (error) {
            console.error('Ошибка при проверке сайтов:', error);
            resultsArea.innerHTML = `<p class="status-error">Произошла ошибка при выполнении запроса: ${error.message}</p>`;
        } finally {
            // Скрываем индикатор загрузки и разблокируем кнопку
            loadingIndicator.classList.add('hidden');
            checkButton.disabled = false;
        }
    });

    function displayResults(results) {
        if (!results || results.length === 0) {
            resultsArea.innerHTML = '<p>Нет результатов для отображения.</p>';
            return;
        }

        results.forEach(result => {
            const resultDiv = document.createElement('div');
            resultDiv.classList.add('result-item');

            let statusText = '';
            let statusClass = '';
            let details = '';

            if (result.is_available) {
                statusText = `Доступен (Код: ${result.status_code})`;
                statusClass = 'status-ok';
                if (result.response_time_ms) {
                     details = `Время ответа: ${result.response_time_ms} мс`;
                }
            } else {
                statusText = 'Недоступен';
                statusClass = 'status-error';
                 if (result.error) {
                     details = `<span class="error-message">Ошибка: ${result.error}</span>`;
                } else if (result.status_code > 0) {
                     details = `Код ответа: ${result.status_code}`;
                } else {
                     details = `<span class="error-message">Неизвестная ошибка сети</span>`;
                }
            }

            resultDiv.innerHTML = `
                <strong>${result.url}:</strong> 
                <span class="${statusClass}">${statusText}</span>
                ${details ? `<br/>${details}` : ''}
            `;
            resultsArea.appendChild(resultDiv);
        });
    }
});