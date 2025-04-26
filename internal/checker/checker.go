package checker

import (
	"fmt"
	"net/http"
	"time"
)

type CheckResult struct {
	URL            string `json:"url"`
	IsAvailable    bool   `json:"is_available"`
	StatusCode     int    `json:"status_code"`
	Error          string `json:"error"`
	ResponseTimeMs int64  `json:"response_time_ms"`
}

func CheckSite(url string) CheckResult {
	result := CheckResult{URL: url}
	client := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		result.IsAvailable = false
		result.Error = fmt.Sprintf("Ошибка при создании запроса: %v", err)
		return result
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 YaBrowser/25.2.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "ru,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Sec-Ch-Ua", `"Not A(Brand";v="8", "Chromium";v="132", "YaBrowser";v="25.2", "Yowser";v="2.5"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "Windows")
	
	startTime := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		result.IsAvailable = false
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.ResponseTimeMs = time.Since(startTime).Milliseconds()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.IsAvailable = true
	} else {
		result.IsAvailable = false
		result.Error = http.StatusText(resp.StatusCode)
	}

	return result
}
