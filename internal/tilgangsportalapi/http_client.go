package tilgangsportalapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client
type Client struct {
	baseURL     string
	HTTPClient  *http.Client
	apiUsername string
	apiPassword string
	cookies     []*http.Cookie
	headers     map[string]string
}

// NewClient creates a new Client
func NewClient(baseURL, apiUsername, apiPassword string) (*Client, error) {
	log.Printf("Creating new client for url %s and with user %s", baseURL, apiUsername)
	c := &Client{
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
		baseURL:     baseURL,
		apiUsername: apiUsername,
		apiPassword: apiPassword,
		headers:     map[string]string{},
	}

	if c.headers == nil {
		log.Println("Setting header for Content-Type")
		c.headers = map[string]string{"Content-Type": "application/json"}
	}

	if c.cookies == nil {
		log.Println("Calling Authenticate function and populating cookie with result")
		cookie, err := c.Authenticate()
		if err != nil {
			return nil, err
		}
		c.cookies = cookie
	}

	return c, nil
}

// GetRequest performs an HTTP get request
func (c *Client) GetRequest(urlStr string) ([]byte, error) {

	urlReqeustStr := c.baseURL + urlStr

	// Create a new HTTP request
	req, err := http.NewRequest("GET", urlReqeustStr, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// Create an HTTP client and send the request
	client := &http.Client{
		Timeout: time.Second * 180,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check for HTTP error response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("response: %s, HTTP request failed with status code: %d", resp.Status, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
	}

	// Important to close the response body
	defer resp.Body.Close()

	return bodyBytes, nil

}

// PostRequest performs an HTTP POST request
func (c *Client) PostRequest(urlStr string, data map[string]interface{}, jsonData map[string]interface{}) (*http.Response, error) {

	var requestBody io.Reader

	urlReqeustStr := c.baseURL + urlStr

	// If jsonData is provided, encode it as JSON
	log.Println("Encoding data for POST request and creating request body")
	if jsonData != nil {
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(jsonBytes)
	} else if data != nil {
		// If data is provided, encode it as form values
		formValues := url.Values{}
		for key, value := range data {
			formValues.Add(key, fmt.Sprintf("%v", value))
		}
		requestBody = strings.NewReader(formValues.Encode())
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", urlReqeustStr, requestBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// Create an HTTP client and send the request
	client := &http.Client{
		Timeout: time.Second * 180,
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Check for HTTP error response and getting message from response body
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var responseBody []map[string]interface{}
		err = json.Unmarshal(bodyBytes, &responseBody)
		if err != nil {
			return nil, err
		}
		message, ok := responseBody[0]["Message"].(string)
		if ok {
			return nil, fmt.Errorf("response: %s, HTTP request failed with status code: %d, message: %s", resp.Status, resp.StatusCode, message)
		}
		return nil, fmt.Errorf("response: %s, HTTP request failed with status code: %d", resp.Status, resp.StatusCode)
	}

	// Important to close the response body
	defer resp.Body.Close()

	return resp, nil
}

// PutRequest performs an HTTP Put request
func (c *Client) PutRequest(urlStr string, data map[string]interface{}, jsonData map[string]interface{}) (*http.Response, error) {

	var requestBody io.Reader

	urlReqeustStr := c.baseURL + urlStr

	// If jsonData is provided, encode it as JSON
	log.Println("Encoding data for PUT request and creating request body")
	if jsonData != nil {
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(jsonBytes)
	} else if data != nil {
		// If data is provided, encode it as form values
		formValues := url.Values{}
		for key, value := range data {
			formValues.Add(key, fmt.Sprintf("%v", value))
		}
		requestBody = strings.NewReader(formValues.Encode())
	}

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", urlReqeustStr, requestBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// Create an HTTP client and send the request
	client := &http.Client{
		Timeout: time.Second * 180,
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Check for HTTP error response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("response: %s, HTTP request failed with status code: %d", resp.Status, resp.StatusCode)
	}

	// Important to close the response body
	defer resp.Body.Close()

	return resp, nil
}
