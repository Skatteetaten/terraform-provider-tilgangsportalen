package tilgangsportalapi

// creates and publishes system role, if the role fails to publish it will be deleted

import (
	"fmt"
	"log"
	"net/http"
)

// CreateAndPublishSystemRole creates a system role and publishes it to an IT shop. If the publish method fails, 
// the role is deleted again and an error is thrown
func (client *Client) CreateAndPublishSystemRole(role SystemRole) (*http.Response, error) {
	// Create the system role
	response, err := client.CreateSystemRole(SystemRole{ // Creating a new system role instance as we don't need all the fields
		Name:            	role.Name,
		L2Ident: 		 	role.L2Ident,
		L3Ident: 			role.L3Ident,
		ApprovalLevel:   	role.ApprovalLevel,
		Description:     	role.Description,
		ProductCategory: 	role.ProductCategory,
	})
	if err != nil {
		return nil, err
	}

	// Publish the system role
	_, err = client.PublishSystemRole(role)
	if err != nil {
		// If the role fails to publish, delete the system role
		log.Printf("Failed to publish role %s, deleting it...", role.Name)

		deleteObject := DeleteSystemRole{
			Name:  role.Name,
			Force: "1",
		}

		_, err = client.DeleteSystemRole(deleteObject)
		if err != nil {
			return nil, fmt.Errorf("failed to publish role %s and failed to delete it: %v", role.Name, err)
		}

		return nil, fmt.Errorf("failed to publish role %s and deleted it", role.Name)
	}

	log.Printf("System Role %s was successfully created and published.", role.Name)

	return response, nil
}
