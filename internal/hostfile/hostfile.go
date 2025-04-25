package hostfile

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// возвращает старндартный для текущей ОС путь к hosts
func getHostsPath() string {
	switch runtime.GOOS {
	case "windows":
		return `C:\Windows\System32\drivers\etc\hosts`
	case "linux", "darwin":
		return `/etc/hosts`
	default:
		return ""
	}
}

// читает hosts и возвращает список доменных имён
func GetDomainsFromHosts() ([]string, error) {
	filePath := getHostsPath()
	if filePath == "" {
		return nil, fmt.Errorf("неизвестная ОС: %s", runtime.GOOS)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл hosts %s: %w", filePath, err)
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		commentIndex := strings.Index(line, "#")
		if commentIndex > 0 {
			line = line[:commentIndex]
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		for i := 1; i < len(fields); i++ {
			domain := "https://" + fields[i]
			domains = append(domains, domain)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения hosts: %w", err)
	}
	return domains, nil
}
