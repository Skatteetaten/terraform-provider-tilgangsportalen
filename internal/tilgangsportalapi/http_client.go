// Package tilgangsportalapi can be used to read, create, modify and delete 
// resources in Tilgangsportalen. The resources available are Entra group,
// system role, role assignments, and lists of these.
package tilgangsportalapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

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
	_, bodyBytes, err := BuildRequest("GET", urlRequestStr,nil, c.headers, c.cookies)

	return bodyBytes, err
}

// PostRequest performs an HTTP POST request
func (c *Client) PostRequest(urlStr string, requestBody io.Reader) (*http.Response, error) {
	urlRequestStr := c.baseURL + urlStr
	response, _, err := BuildRequest("POST", urlRequestStr,requestBody, c.headers, c.cookies)

	return response,err
}

// PutRequest performs an HTTP Put request
func (c *Client) PutRequest(urlStr string, requestBody io.Reader) (*http.Response, error) {
	urlRequestStr := c.baseURL + urlStr
	response, _, err := BuildRequest("PUT", urlRequestStr,requestBody, c.headers, c.cookies)

	return response,err
}


// BuildRequest builds and performs a new HTTP request to urlRequestStr of type
// requestType with requestBody (use nil for GET). Returns the received
// response and the body
func BuildRequest(requestType string, urlRequestStr string, requestBody io.Reader, headers map[string]string, cookies []*http.Cookie)(*http.Response, []byte, error){

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

	// Create an HTTP client and send the request
	client := &http.Client{
		Timeout: time.Second * 180,
	}

	log.Printf("Performing %s request to %s", requestType, urlRequestStr)
	
	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	// Check for HTTP error response and getting message from response body
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if err != nil {
			return nil, nil, err
		}
		var responseBody []map[string]interface{}
		err = json.Unmarshal(bodyBytes, &responseBody)
		if err != nil {
			return nil, nil, err
		}
		message, ok := responseBody[0]["Message"].(string)
		if ok {
			return nil, nil, fmt.Errorf("response: %s, HTTP request failed with status code: %d, message: %s", resp.Status, resp.StatusCode, message)
		}
		return nil, nil, fmt.Errorf("response: %s, HTTP request failed with status code: %d", resp.Status, resp.StatusCode)
	}

	// Closing the response body
	defer resp.Body.Close()

	return resp, bodyBytes, nil

}


// CreateRequestBody creates a body of type io.Reader that can be used i an 
// API call
func CreateRequestBody(bodyObject interface{}) (io.Reader, error){
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