package tilgangsportalapi

import (
	"log"
	"net/http"
)

// Authenticate performs authentication with the given username and password
func (c *Client) Authenticate() ([]*http.Cookie, error) {
	log.Println("Authenticating using username and password...")

	// Prepare the request body
	authBody := map[string]interface{}{
		"Module":   "DialogUser",
		"Password": c.apiPassword,
		"User":     c.apiUsername,
	}

	// Construct the URL
	authURL := "/imx/login/SKAT_RoleGovernance"

	// Perform the POST request
	response, err := c.PostRequest(authURL, nil, authBody)
	if err != nil {
		log.Println("Received error when attempting to authenticate to the server.")
		return nil, err
	}

	// Extract cookies from the response
	cookies := response.Cookies()

	log.Println("Authentication successful.")
	return cookies, nil
}
