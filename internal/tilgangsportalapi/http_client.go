// Package tilgangsportalapi can be used to read, create, modify and delete
// resources in Tilgangsportalen. The resources available are Entra group,
// system role, role assignments, and lists of these.
package tilgangsportalapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const defaultTimeout = 300

// Client struct
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
	urlRequestStr := c.baseURL + urlStr
	_, bodyBytes, err := BuildRequest("GET", urlRequestStr, nil, c.headers, c.cookies, defaultTimeout)

	return bodyBytes, err
}

// PostRequest performs an HTTP POST request
func (c *Client) PostRequest(urlStr string, requestBody io.Reader) (*http.Response, error) {
	urlRequestStr := c.baseURL + urlStr
	response, _, err := BuildRequest("POST", urlRequestStr, requestBody, c.headers, c.cookies, defaultTimeout)

	return response, err
}

// PostRequest performs an HTTP POST request
func (c *Client) PostRequestWithCustomTimeout(urlStr string, requestBody io.Reader, customTimeout int) (*http.Response, error) {
	urlRequestStr := c.baseURL + urlStr
	response, _, err := BuildRequest("POST", urlRequestStr, requestBody, c.headers, c.cookies, customTimeout)

	return response, err
}

// PutRequest performs an HTTP Put request
func (c *Client) PutRequest(urlStr string, requestBody io.Reader) (*http.Response, error) {
	urlRequestStr := c.baseURL + urlStr
	response, _, err := BuildRequest("PUT", urlRequestStr, requestBody, c.headers, c.cookies, defaultTimeout)

	return response, err
}

// BuildRequest builds and performs a new HTTP request to urlRequestStr of type
// requestType with requestBody (use nil for GET). Returns the response and its
// body on 2xx, or an *APIError on any other status.
func BuildRequest(requestType string, urlRequestStr string, requestBody io.Reader, headers map[string]string, cookies []*http.Cookie, timeout int) (*http.Response, []byte, error) {

	log.Printf("Encoding data for %s request and creating request body", requestType)

	req, err := http.NewRequest(requestType, urlRequestStr, requestBody)
	if err != nil {
		return nil, nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Timeout is given in seconds and covers the whole request, including reading the body
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	log.Printf("Performing %s request to %s", requestType, urlRequestStr)

	// Only fails on network/timeout errors; HTTP error statuses are handled below
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	// Close the body on every return path, also when the API responds with an error
	defer func() { _ = resp.Body.Close() }()

	// The body is read here, so callers must use bodyBytes instead of resp.Body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Check for HTTP error response and getting message from response body
	// HTTP 200 is success. Any other status is treated as an API error.
	// The API responds to all errors with HTTP 500 and the real error code (6xx) at the start of Message
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, newAPIError(resp.Status, resp.StatusCode, parseErrorMessage(bodyBytes))
	}

	return resp, bodyBytes, nil
}

// parseErrorMessage extracts Message from an API error body of the form
// [{"Message": "..."}]. If the body is not in that form (e.g. an HTML page
// from a gateway), the raw body is returned so the HTTP status is still
// reported as an APIError.
func parseErrorMessage(body []byte) string {
	var responseBody []struct {
		Message string `json:"Message"`
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return strings.TrimSpace(string(body))
	}
	if len(responseBody) == 0 {
		return ""
	}
	return responseBody[0].Message
}

// CreateRequestBody creates a body of type io.Reader that can be used i an
// API call
func CreateRequestBody(bodyObject interface{}) (io.Reader, error) {
	var body map[string]interface{}
	tempBody, err := json.Marshal(bodyObject)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tempBody, &body)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBytes), nil
}

// GetBaseURL returns the base URL of the client.
func (c *Client) GetBaseURL() string {
	return c.baseURL
}
