package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
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
