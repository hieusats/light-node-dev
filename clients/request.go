package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Layer-Edge/light-node/utils"
	"github.com/go-resty/resty/v2"
)

// Default timeout in seconds if environment variable is not set
const DEFAULT_TIMEOUT = 100

// RequestOptions contains options for the request
type RequestOptions struct {
	Proxy   string // Proxy URL in format http://user:pass@host:port or socks5://user:pass@host:port
	Timeout int    // Timeout in seconds
}

func PostRequest[T any, R any](url string, requestData T, options ...RequestOptions) (*R, error) {
	client := resty.New()

	// Get timeout from environment variable or use default
	timeout := DEFAULT_TIMEOUT
	envTimeout := utils.GetEnv("API_REQUEST_TIMEOUT", "100")
	if t, err := strconv.Atoi(envTimeout); err == nil {
		timeout = t
	}

	// Apply options if provided
	if len(options) > 0 {
		// Override timeout if specified in options
		if options[0].Timeout > 0 {
			timeout = options[0].Timeout
		}

		// Set proxy if provided
		if options[0].Proxy != "" {
			client.SetProxy(options[0].Proxy)
		}
	}

	// Set default headers, timeout
	client.
		SetTimeout(time.Second*time.Duration(timeout)).
		SetHeader("Authorization", "Bearer your-token-here")

	// Make request
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestData).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	// Check status code
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode(), string(resp.Body()))
	}

	// Parse response into generic type
	var response R
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &response, nil
}
