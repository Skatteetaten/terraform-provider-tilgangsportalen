package tilgangsportalapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
)

// Gets a list of all the roles that are assigned to an Azure AD Group, along with the owner of each role.
// This resource shows all the roles that are assigned to a Entra Group owned by the current user (also roles not owned by current user).
// See https://wiki.sits.no/display/IDABAS/22.+List+Roles+for+Azure+AD+Groups
func (client *Client) ListAllSystemRolesForEntraGroup(GroupName string) (*RolesWithOwner, error) {
	var data RolesWithOwner
	log.Printf("Listing roles for Entra Group %s ...", GroupName)
	// Construct the URL
	listEntraGroupsForRoleURL := fmt.Sprintf("/SKAT_RoleGovernance/ListRolesForAzureADGroup?GroupName=%s", url.QueryEscape(GroupName))

	// Perform the POST request
	response, err := client.GetRequest(listEntraGroupsForRoleURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the Roles struct
	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}

	// lowercase the L2Ident
	for i := range data.Roles {
		data.Roles[i].L2Ident = strings.ToLower(data.Roles[i].L2Ident)
	}

	log.Printf("Listing roles for Entra Group %s successful. Found %d role(s).", GroupName, len(data.Roles))

	return &data, nil
}
