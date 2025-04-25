document.addEventListener('DOMContentLoaded', () => {
    // Получаем ссылки на элементы, включая новую кнопку
    const urlsInput = document.getElementById('urlsInput');
    const checkButton = document.getElementById('checkButton');
    const loadHostsButton = document.getElementById('loadHostsButton'); // Получаем новую кнопку
    const resultsArea = document.getElementById('resultsArea');
    const loadingIndicator = document.getElementById('loadingIndicator');

    // --- Обработчик для кнопки "Проверить" (уже был, без изменений) ---
    checkButton.addEventListener('click', async () => {
        const urlsText = urlsInput.value.trim();
        const urls = urlsText.split('\n').map(url => url.trim()).filter(url => url !== '');

        if (urls.length === 0) {
            resultsArea.innerHTML = '<p>Пожалуйста, введите хотя бы один URL для проверки.</p>';
            return;
        }

        // Очистка предыдущих результатов и показ индикатора
        resultsArea.innerHTML = '';
        loadingIndicator.classList.remove('hidden');
        checkButton.disabled = true; // Блокируем обе кнопки
        loadHostsButton.disabled = true;

        try {
            // Отправляем запрос на бэкенд (эндпоинт /check)
            const response = await fetch('/check', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(urls),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Ошибка сервера: ${response.status} ${response.statusText}. ${errorText}`);
            }

            const results = await response.json();

            // Отображаем результаты
            displayResults(results);

        } catch (error) {
            console.error('Ошибка при проверке сайтов:', error);
            resultsArea.innerHTML = `<p class="status-error">Произошла ошибка при выполнении запроса: ${error.message}</p>`;
        } finally {
            // Скрываем индикатор загрузки и разблокируем обе кнопки
            loadingIndicator.classList.add('hidden');
            checkButton.disabled = false;
            loadHostsButton.disabled = false;
        }
    });

    // --- НОВЫЙ ОБРАБОТЧИК для кнопки "Загрузить из hosts" ---
    loadHostsButton.addEventListener('click', async () => {
        // Очистка предыдущих результатов и показ индикатора
        resultsArea.innerHTML = ''; // Очищаем область результатов
        loadingIndicator.classList.remove('hidden'); // Показываем индикатор загрузки
        checkButton.disabled = true; // Блокируем обе кнопки на время запроса
        loadHostsButton.disabled = true;

        try {
            // Отправляем GET запрос на эндпоинт для hosts файла
            // URL должен соответствовать пути, который ты зарегистрировал в main.go
            const response = await fetch('/get-hosts'); 

            if (!response.ok) {
                const errorText = await response.text();
                 throw new Error(`Ошибка сервера при загрузке hosts: ${response.status} ${response.statusText}. ${errorText}`);
            }

            // Ожидаем получить JSON массив строк (список доменов)
            // Пример: ["localhost", "my-local-site.local", ...]
            const domains = await response.json();

            // Убедимся, что получили массив
            if (!Array.isArray(domains)) {
                 throw new Error('Неверный формат данных от сервера: ожидается массив доменов.');
            }

            // Преобразуем массив доменов в строку, где каждый домен на новой строке
            const domainsText = domains.join('\n');

            // Помещаем полученные домены в текстовое поле для URL
            urlsInput.value = domainsText;

            // Опционально: Отобразить сообщение об успехе в области результатов
            resultsArea.innerHTML = '<p>Домены из hosts загружены в текстовое поле.</p>';


        } catch (error) {
            console.error('Ошибка при загрузке hosts:', error);
            // Отображаем ошибку в области результатов
            resultsArea.innerHTML = `<p class="status-error">Произошла ошибка при загрузке hosts: ${error.message}</p>`;
        } finally {
            // Скрываем индикатор загрузки и разблокируем обе кнопки
            loadingIndicator.classList.add('hidden');
            checkButton.disabled = false;
            loadHostsButton.disabled = false;
        }
    });


    // --- Функция отображения результатов (уже была, без изменений) ---
    function displayResults(results) {
        // ... (код функции displayResults) ...
        if (!results || results.length === 0) {
            resultsArea.innerHTML = '<p>Нет результатов для отображения.</p>';
            return;
        }

        resultsArea.innerHTML = ''; // Очищаем перед добавлением новых результатов

        results.forEach(result => {
            const resultDiv = document.createElement('div');
            resultDiv.classList.add('result-item');

            let statusText = '';
            let statusClass = '';
            let details = '';

            if (result.is_available) {
                statusText = `Доступен (Код: ${result.status_code})`;
                statusClass = 'status-ok';
                if (result.response_time_ms !== undefined && result.response_time_ms !== null) { // Проверка на undefined/null
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