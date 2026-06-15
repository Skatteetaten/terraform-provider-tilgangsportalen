package tilgangsportalapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"
)

// GetSystemRole gets the name of the Terms of Use information assigned to a specific role.
// Returns Terms of use identifier (name), description and UID.
// See https://wiki.sits.no/spaces/OIM/pages/1333270032/23.+Get+Terms+Of+Use+for+Role
func (client *Client) GetTermsOfUseForRole(roleName string) (*TermsOfUse, error) {
	var data TermsOfUse
	log.Printf("Fetching terms of use object attached to system role %s", roleName)
	// Construct the URL, with query escape to handle special characters in role name
	getTermsOfUseURL := "/SKAT_RoleGovernance/GetTermsOfUseForRole?RoleName=" + url.QueryEscape(roleName)
	// Perform the GET request
	response, err := client.GetRequest(getTermsOfUseURL)
	// Error handling
	if err != nil {
		// Check for specific error codes in the error message to provide more context-specific logs and errors
		errMsg := err.Error()

		if IsErrorRoleDoesNotExistOrUnauthorized(errMsg) {
			log.Printf("Role with name \"%s\" doesn't exist or you don't have permissions against it.", roleName)
			return nil, fmt.Errorf("role '%s' doesn't exist or you don't have permissions against it", roleName)
		}
		if IsErrorNoTermsOfUseAssignedToRole(errMsg) {
			log.Printf("No terms of use is assigned to role \"%s\".", roleName)
			return nil, fmt.Errorf("no terms of use is assigned to role '%s'", roleName)
		}
		log.Printf("Error fetching terms of use for role \"%s\": %v", roleName, err)
		return nil, err
	}

	// Unmarshal the JSON data into the TermsOfUse struct
	err = json.Unmarshal(response, &data)
	if err != nil {
		log.Printf("An error was thrown when unmarshaling response from tilgangsportalen.")
		return nil, err
	}

	log.Printf("Terms of use \"%s\" attached to role with name \"%s\" was found.", data.TermsOfUseIdentifier, roleName)

	return &data, nil
}

// WaitForTermsOfUseAssignment waits for a specific terms of use to be assigned to a role
func (client *Client) WaitForTermsOfUseAssignment(termsOfUse string, roleName string) (*TermsOfUse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	log.Printf("Waiting for terms of use %s to be assigned to role %s...", termsOfUse, roleName)

	for {
		termsOfUseData, err := client.GetTermsOfUseForRole(roleName)

		if err != nil {
			// Check for specific error code 606 (no terms of use assigned)
			if IsErrorNoTermsOfUseAssignedToRole(err.Error()) {
				// Continue waiting - no terms of use assigned yet
			} else {
				// For any other error, return it
				return nil, err
			}
		} else {
			// No error means terms of use is assigned, check if it's the one we want
			if termsOfUseData != nil && termsOfUseData.TermsOfUseIdentifier == termsOfUse {
				log.Printf("Terms of use %s is assigned to role %s", termsOfUse, roleName)
				return termsOfUseData, nil // Assignment successful - return the data
			}
			// Different terms of use is assigned, continue waiting
		}

		select {
		case <-time.After(5 * time.Second): // Retry after 5 seconds
		case <-ctx.Done(): // Return error if timeout is reached
			// Wrap ctx.Err() so callers can detect timeout/cancel, while still getting a readable message with lookup context.
			timeoutErr := fmt.Errorf("timeout reached while waiting for terms of use %s to be assigned to role %s: %w", termsOfUse, roleName, ctx.Err())
			return nil, timeoutErr
		}
	}
}

// WaitForTermsOfUseRemoval waits for a specific terms of use to be removed from a role
func (client *Client) WaitForTermsOfUseRemoval(roleName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	log.Printf("Waiting for terms of use to be removed from role %s...", roleName)

	for {
		_, err := client.GetTermsOfUseForRole(roleName)

		if err != nil {
			// Check for specific error code 606 (no terms of use assigned)
			if IsErrorNoTermsOfUseAssignedToRole(err.Error()) {
				log.Printf("A terms of use object is no longer assigned to role %s", roleName)
				return nil // Removal successful - no terms of use assigned
			} else {
				// For any other error, return it
				return err
			}
		}

		select {
		case <-time.After(5 * time.Second): // Retry after 5 seconds
		case <-ctx.Done():
			timeoutErr := fmt.Errorf("timeout reached while waiting for terms of use to be removed from role %s: %w", roleName, ctx.Err())
			log.Printf("Timeout reached while waiting for terms of use to be removed from role %s", roleName)
			return timeoutErr // Return error if timeout is reached
		}
	}
}
