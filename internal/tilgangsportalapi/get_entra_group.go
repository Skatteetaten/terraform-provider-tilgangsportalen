package tilgangsportalapi

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"
)

// GetEntraGroup gets information about a specific Entra ID group if group exists.
// Gets the alias, tenant, inheritance level, and description for a specific
// (named) Entra ID group.
// See https://wiki.sits.no/display/IDABAS/21.+Get+Azure+AD+Group
func (client *Client) GetEntraGroup(entraGroupName string, waitForObjectId bool) (*EntraGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var data EntraGroup
	log.Printf("Fetching Entra group %s", entraGroupName)

	// Construct the URL, with query escape to handle special characters in role name
	getGroupURL := "/SKAT_RoleGovernance/GetAzureADGroupV1?DisplayName=" + url.QueryEscape(entraGroupName)

	// Get the Entra group from the API, and poll until it returns a value for the group object_id
	for {
		// Perform the GET request
		response, err := client.GetRequest(getGroupURL)
		if err != nil {
			log.Printf("Entra ID group with name \"%s\" was not found.", entraGroupName)
			return nil, err
		}

		// Unmarshal the JSON data into the EntraGroup struct
		err = json.Unmarshal(response, &data)
		if err != nil {
			log.Printf("An error was thrown when unmarshaling response from tilgangsportalen.")
			return nil, err
		}

		// Stop polling if specified not to wait for object_id to be set, or if the object_id is set
		if !waitForObjectId || data.EntraIDOID != "" {
			break
		}

		select {
		case <-time.After(10 * time.Second): // Retry after 10 seconds
		case <-ctx.Done():
			log.Printf("Timeout reached while waiting for object_id to be created for group %s", entraGroupName)
			return nil, ctx.Err() // Return error if timeout is reached
		}
	}

	// Manually set the DisplayName as its not returned by the API
	if data.DisplayName == "" {
		data.DisplayName = entraGroupName
	}

	// The API returns inheritance level in lowercase, ensure the inheritance level is capitalized correctly
	data.InheritanceLevel = strings.ToUpper(data.InheritanceLevel[:1]) + strings.ToLower(data.InheritanceLevel[1:])

	log.Printf("Entra ID group with name %s was found.", entraGroupName)

	return &data, nil
}
