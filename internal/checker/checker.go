package checker

import (
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
	startTime := time.Now()

	resp, err := client.Get(url)
	result.ResponseTimeMs = time.Since(startTime).Milliseconds()
	if err != nil {
		result.IsAvailable = false
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.IsAvailable = true
	} else {
		result.IsAvailable = false
		result.Error = http.StatusText(resp.StatusCode)
	}

	return result
}
