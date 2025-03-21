package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	proxyList   []string
	proxyMutex  sync.Mutex
	proxyLoaded bool
)

// LoadProxiesFromFile loads proxies from proxy.txt file
func LoadProxiesFromFile() ([]string, error) {
	proxyMutex.Lock()
	defer proxyMutex.Unlock()

	if proxyLoaded && len(proxyList) > 0 {
		return proxyList, nil
	}

	file, err := os.Open("proxy.txt")
	if err != nil {
		// If proxy.txt doesn't exist, return empty list
		return nil, fmt.Errorf("failed to open proxy.txt: %v", err)
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" {
			proxies = append(proxies, proxy)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading proxy.txt: %v", err)
	}

	proxyList = proxies
	proxyLoaded = true
	return proxies, nil
}

// GetProxyByIndex returns a proxy by index
func GetProxyByIndex(index int) (string, error) {
	proxies, err := LoadProxiesFromFile()
	if err != nil {
		return "", err
	}

	if len(proxies) == 0 {
		return "", fmt.Errorf("no proxies available")
	}

	if index < 0 || index >= len(proxies) {
		return "", fmt.Errorf("proxy index out of range")
	}

	return proxies[index], nil
}

// FormatProxyURL converts a proxy string to the correct URL format for HTTP clients
// Input format: host:port:username:password
// Output format: http://username:password@host:port
func FormatProxyURL(proxy string) string {
	if proxy == "" {
		return ""
	}

	parts := strings.Split(proxy, ":")
	if len(parts) < 2 {
		// Invalid format, return as is
		return proxy
	}

	// Basic format: host:port
	if len(parts) == 2 {
		return fmt.Sprintf("http://%s:%s", parts[0], parts[1])
	}

	// Format with credentials: host:port:username:password
	host := parts[0]
	port := parts[1]

	var username, password string

	if len(parts) == 3 {
		// Format: host:port:username (no password)
		username = parts[2]
		password = ""
	} else if len(parts) >= 4 {
		// Format: host:port:username:password
		username = parts[2]
		password = parts[3]

		// If there are more parts, they might be part of the password (if it contains colons)
		if len(parts) > 4 {
			additionalParts := parts[4:]
			for _, part := range additionalParts {
				password += ":" + part
			}
		}
	}

	// Construct the proxy URL in the format http://username:password@host:port
	if password == "" {
		return fmt.Sprintf("http://%s@%s:%s", username, host, port)
	}
	return fmt.Sprintf("http://%s:%s@%s:%s", username, password, host, port)
}

// TestProxy kiểm tra xem proxy có hoạt động không bằng cách gửi yêu cầu đến ipinfo.io
func TestProxy(proxy string) (bool, string, error) {
	if proxy == "" {
		return false, "Không có proxy được cung cấp", nil
	}

	// Format proxy URL
	formattedProxy := FormatProxyURL(proxy)

	// Tạo HTTP client với proxy
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: func(_ *http.Request) (*url.URL, error) {
				return url.Parse(formattedProxy)
			},
		},
	}

	// Gửi yêu cầu đến ipinfo.io để kiểm tra IP
	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		return false, "", fmt.Errorf("lỗi khi kiểm tra proxy: %v", err)
	}
	defer resp.Body.Close()

	// Đọc phản hồi
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("lỗi khi đọc phản hồi: %v", err)
	}

	// Kiểm tra xem phản hồi có chứa thông tin IP không
	if !strings.Contains(string(body), "ip") {
		return false, "Phản hồi không chứa thông tin IP", nil
	}

	return true, string(body), nil
}

// CheckAllProxies kiểm tra tất cả các proxy trong file proxy.txt
func CheckAllProxies() (map[string]string, error) {
	proxies, err := LoadProxiesFromFile()
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)
	for i, proxy := range proxies {
		success, info, err := TestProxy(proxy)
		if err != nil {
			results[proxy] = fmt.Sprintf("Lỗi: %v", err)
		} else if success {
			results[proxy] = fmt.Sprintf("Hoạt động: %s", info)
		} else {
			results[proxy] = "Không hoạt động"
		}

		// Log kết quả
		log.Printf("Proxy %d: %s - %s", i+1, proxy, results[proxy])
	}

	return results, nil
}
