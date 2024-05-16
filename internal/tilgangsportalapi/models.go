package tilgangsportalapi


////    Entra Group    ////

type EntraGroup struct {
	DisplayName      string `json:"DisplayName"`
	Alias            string `json:"Alias"`
	Tenant           string `json:"Tenant"`
	InheritanceLevel string `json:"InheritanceLevel"`
	Description      string `json:"Description"`
}

// EntraGroups represents the structure of the API response from ListAzureADGroups
type EntraGroups struct {
	EntraGroups []EntraGroup `json:"Groups"`
}

type DeleteEntraGroup struct {
	Name  string `json:"Name"`
	Force string `json:"Force"`
	// consider implementing force flag for removing any group assignments before deletion
	// if delete fails due to group still being assigned to a role
}

type RenameEntraGroup struct {
	OldName string `json:"OldName"`
	NewName string `json:"NewName"`
}

////   Role assigment    ////

type EntraGroupRoleAssignment struct {
	RoleName   string `json:"Name"`
	EntraGroup string `json:"Entitlement"`
}

////    System role     ////

type SystemRole struct {
	Name            		 string `json:"Name"`
	L2Ident 				 string `json:"L2Ident"`
	L3Ident 				 string `json:"L3Ident"`
	ApprovalLevel   		 string `json:"ApprovalLevel"`
	Description     		 string `json:"Description"`
	ProductCategory 		 string `json:"ProductCategory"`
	ItShopName      		 string `json:"ITShopName"`
}

type PublishSystemRole struct {
	Name   string `json:"Name"`
	ITShop string `json:"ITShop"`
}

type SystemRoleChange struct {
	RoleName         		string `json:"RoleName"`
	SystemRoleOwner  		string `json:"SystemRoleOwner"`
	SystemRoleSecurityOwner string `json:"SystemRoleSecurityOwner"`
	NewApprovalLevel 		string `json:"NewApprovalLevel"`
	NewDescription   		string `json:"NewDescription"`
	ProductCategory  		string `json:"ProductCategory"`
}

type RenameSystemRole struct {
	OldName string `json:"OldName"`
	NewName string `json:"NewName"`
}

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
