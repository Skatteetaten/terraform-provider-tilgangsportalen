package tilgangsportalapi

////    Entra Group    ////

// EntraGroup represents the API body for creating an Entra group
type EntraGroup struct {
	DisplayName      string `json:"DisplayName"`
	Alias            string `json:"Alias"`
	Tenant           string `json:"Tenant"`
	InheritanceLevel string `json:"InheritanceLevel"`
	Description      string `json:"Description"`
}

// EntraGroups represents the structure of the API response from
// ListAzureADGroups
type EntraGroups struct {
	EntraGroups []EntraGroup `json:"Groups"`
}

// DeleteEntraGroup represents the API body for deleting an Entra group
// The group is identified by its name (id) and the force parameter is
// set to "1" to force deletion even if the group is assigned to a role
type DeleteEntraGroup struct {
	Name  string `json:"Name"`
	Force string `json:"Force"`
	// consider implementing force flag for removing any group assignments
	// before deletion if delete fails due to group still being assigned to
	// a role
}

// RenameEntraGroup represents the API body for renaming an Entra group
type RenameEntraGroup struct {
	OldName string `json:"OldName"`
	NewName string `json:"NewName"`
}

////   Role assigment    ////

// EntraGroupRoleAssignment represents the API body for creating a role
// assignment between a group named EntraGroup and a role named RoleName
type EntraGroupRoleAssignment struct {
	RoleName   string `json:"Name"`
	EntraGroup string `json:"Entitlement"`
}

////    System role     ////

// SystemRole represents the API body for creating a system role
type SystemRole struct {
	Name            string `json:"Name"`
	L2Ident         string `json:"L2Ident"`
	L3Ident         string `json:"L3Ident"`
	ApprovalLevel   string `json:"ApprovalLevel"`
	Description     string `json:"Description"`
	ProductCategory string `json:"ProductCategory"`
	ItShopName      string `json:"ITShopName"`
}

// SystemRoleChange represents the API body for modifying the fields of
// a role
type SystemRoleChange struct {
	RoleName         string `json:"RoleName"`
	L2Ident          string `json:"L2Ident"`
	L3Ident          string `json:"L3Ident"`
	NewApprovalLevel string `json:"NewApprovalLevel"`
	NewDescription   string `json:"NewDescription"`
	ProductCategory  string `json:"ProductCategory"`
}

// RenameSystemRole represents the API body for modifying the name of
// a role
type RenameSystemRole struct {
	OldName string `json:"OldName"`
	NewName string `json:"NewName"`
}

// DeleteSystemRole represents the API body for deleting a role identified
// by its name. If Force is set to "1" any role assignments belonging to the
// role is also cleaned up in the backend
type DeleteSystemRole struct {
	Name  string `json:"Name"`
	Force string `json:"Force"`
}

////       Role         ////

// Role represents a single role with a DisplayName
type Role struct {
	DisplayName string `json:"DisplayName"`
}

// Roles represents the structure of the API response for ListRoles
type Roles struct {
	Roles []Role `json:"Roles"`
}
