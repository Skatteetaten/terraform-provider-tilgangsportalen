package tilgangsportalapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

const (
	entraGroupLookupByDisplayNameEndpoint = "GetAzureADGroupV1"
	entraGroupLookupByDisplayNameQueryKey = "DisplayName"
	entraGroupLookupByEntitlementEndpoint = "GetAzureADGroupID"
	entraGroupLookupByEntitlementQueryKey = "EntitlementUID"
)

// GetEntraGroup gets information about a specific Entra ID group if group exists.
// Gets the alias, tenant, inheritance level, and description for a specific
// (named) Entra ID group.
// See https://wiki.sits.no/spaces/OIM/pages/1380418926/21.+Get+Azure+AD+Group+V1
func (client *Client) GetEntraGroup(entraGroupName string, waitForObjectId bool) (*EntraGroup, error) {
	return client.getEntraGroup(entraGroupLookupByDisplayNameEndpoint, entraGroupLookupByDisplayNameQueryKey, entraGroupName, waitForObjectId)
}

// GetEntraGroupByID gets information about a specific Entra ID group by its EntitlementUID.
// See https://wiki.sits.no/spaces/OIM/pages/1460147575/21.+Get+Azure+AD+Group+ID
func (client *Client) GetEntraGroupByID(entitlementUID string, waitForObjectId bool) (*EntraGroup, error) {
	return client.getEntraGroup(entraGroupLookupByEntitlementEndpoint, entraGroupLookupByEntitlementQueryKey, entitlementUID, waitForObjectId)
}

func (client *Client) getEntraGroup(endpointName string, queryKey string, queryValue string, waitForObjectId bool) (*EntraGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var data EntraGroup
	log.Printf("Fetching Entra group with %s %s", queryKey, queryValue)

	// Construct the URL, with query escape to handle special characters.
	getGroupURL := "/SKAT_RoleGovernance/" + endpointName + "?" + queryKey + "=" + url.QueryEscape(queryValue)

	// Get the Entra group from the API, and poll until it returns a value for the group object_id
	for {
		// Perform the GET request
		response, err := client.GetRequest(getGroupURL)
		if err != nil {
			// If waitForObjectId is true, we will retry fetching the group until it has an object_id,
			// as it may take some time for the group to be created and available in the API.
			// If waitForObjectId is false, we will not retry and return the error immediately.
			if waitForObjectId {
				log.Printf("An error was thrown when fetching Entra group with %s \"%s\". Error: %v. Retrying in %v...", queryKey, queryValue, err, 10*time.Second)
				select {
				case <-time.After(10 * time.Second): // Retry after 10 seconds
					continue
				case <-ctx.Done(): // Return error if timeout is reached
					// Wrap ctx.Err() so callers can detect timeout/cancel, while still getting a readable message with lookup context.
					timeoutErr := fmt.Errorf("timeout while retrying failed Entra group lookup for %s %q; object_id was not observed before deadline: %w", queryKey, queryValue, ctx.Err())
					return nil, timeoutErr
				}
			} else {
				log.Printf("Entra ID group with %s \"%s\" was not found.", queryKey, queryValue)
				return nil, err
			}
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
		case <-ctx.Done(): // Return error if timeout is reached
			// Wrap ctx.Err() so callers can detect timeout/cancel, while still getting a readable message with lookup context.
			timeoutErr := fmt.Errorf("timeout reached while waiting for object_id to be created for group with %s %q: %w", queryKey, queryValue, ctx.Err())
			return nil, timeoutErr
		}
	}

	// Manually set the DisplayName for DisplayName-based lookups as it may not be returned by the API.
	if queryKey == entraGroupLookupByDisplayNameQueryKey && data.DisplayName == "" {
		data.DisplayName = queryValue
	}

	// Manually set the EntitlementUID for EntitlementUID-based lookups as it may not be returned by the API.
	if queryKey == entraGroupLookupByEntitlementQueryKey && data.EntitlementUID == "" {
		data.EntitlementUID = queryValue
	}

	// The API returns inheritance level in lowercase, ensure the inheritance level is capitalized correctly
	if data.InheritanceLevel != "" {
		data.InheritanceLevel = strings.ToUpper(data.InheritanceLevel[:1]) + strings.ToLower(data.InheritanceLevel[1:])
	}

	log.Printf("Entra ID group with %s %s was found.", queryKey, queryValue)

	return &data, nil
}
