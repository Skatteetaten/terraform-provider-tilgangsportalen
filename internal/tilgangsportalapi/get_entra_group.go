package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

// GetEntraGroup gets information about a specific Entra ID group if group exists.
// Gets the alias, tenant, inheritance level, and description for a specific
// (named) Entra ID group.
// See https://wiki.sits.no/display/IDABAS/21.+Get+Azure+AD+Group
func (client *Client) GetEntraGroup(entraGroupName string) (*EntraGroup, error) {
	var data EntraGroup
	log.Printf("Fetching Entra group %s", entraGroupName)

	// Construct the URL, with query escape to handle special characters in role name
	getRoleURL := "/SKAT_RoleGovernance/GetAzureADGroup?DisplayName=" + url.QueryEscape(entraGroupName)
	// Perform the POST request
	response, err := client.GetRequest(getRoleURL)
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

	// Manually set the DisplayName as its not returned by the API
	if data.DisplayName == "" {
		data.DisplayName = entraGroupName
	}

	// The API returns inheritance level in lowercase, ensure the inheritance level is capitalized correctly
	data.InheritanceLevel = strings.ToUpper(data.InheritanceLevel[:1]) + strings.ToLower(data.InheritanceLevel[1:])

	log.Printf("Entra ID group with name %s was found.", entraGroupName)

	return &data, nil
}
