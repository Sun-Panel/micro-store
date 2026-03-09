package paddle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"sun-panel/lib/paddle/types"
)

const (
	DefaultBaseURL = "https://api.paddle.com/"
	SandboxBaseURL = "https://sandbox-api.paddle.com/"
)

var Webhook *WebhookType

type Client struct {
	client       *http.Client // HTTP client used to communicate with the API.
	commonServer service      // Reuse a single struct instead of allocating one for each service on the heap.
	BaseURL      *url.URL
	Token        string

	Transaction *TransactionService
}

type service struct {
	client *Client
}

func (c *Client) get(url string, v interface{}) (*http.Response, error) {
	url = c.BaseURL.Scheme + "://" + c.BaseURL.Host + url
	// fmt.Println("请求URL", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		// fmt.Println("3663666", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("555555252", err)
		return nil, err
	}
	// fmt.Println("返回数据", string(body))
	if err := checkResponse(resp, body); err != nil {
		// fmt.Println("555555", err)
		return resp, err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return resp, fmt.Errorf("unmarshal error %s", err.Error())
	}
	return resp, nil
}

func (c *Client) post(url string, data interface{}, v interface{}) (*http.Response, error) {
	url = c.BaseURL.Scheme + "://" + c.BaseURL.Host + url
	// fmt.Println("请求URL", url)
	// fmt.Println("请求数据", string(optionToJsonBytes(data)))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(optionToJsonBytes(data)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		// fmt.Println("3663666", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("555555252", err)
		return nil, err
	}
	// fmt.Println("返回数据", string(body))
	if err := checkResponse(resp, body); err != nil {
		// fmt.Println("555555", err)
		return resp, err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return resp, fmt.Errorf("unmarshal error %s", err.Error())
	}
	return resp, nil
}

func NewClient(baseURL, token string, httpClient *http.Client) *Client {
	url, _ := url.Parse(baseURL)
	return getClient(httpClient, url, token)
}

func getClient(httpClient *http.Client, baseURL *url.URL, token string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
		Token:   token,
	}

	c.commonServer.client = c
	c.Transaction = (*TransactionService)(&c.commonServer)

	return c
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Float64 is a helper routine that allocates a new float64 value
// to store v and returns a pointer to it.
func Float64(v float64) *float64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

func optionToJsonBytes(options interface{}) []byte {
	res, _ := json.Marshal(options)
	return res
}

// Check wether or not the API response contains an error
func checkResponse(r *http.Response, data []byte) error {
	errorResponse := &types.ErrorResponse{Response: r}
	if data != nil {
		if err := json.Unmarshal(data, errorResponse); err != nil {
			return err
		}
	}
	if errorResponse.Error != nil {
		return fmt.Errorf("error: %s, %s",
			errorResponse.Error.Type,
			errorResponse.Error.Detail)
	}
	return nil
}
