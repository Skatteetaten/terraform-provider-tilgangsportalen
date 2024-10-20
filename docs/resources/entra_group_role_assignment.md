---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tilgangsportalen_entra_group_role_assignment Resource - tilgangsportalen"
subcategory: ""
description: |-
  This resource is used to create assignments between Entra Groups and System Roles in Tilgangsportalen
---

# tilgangsportalen_entra_group_role_assignment (Resource)

This resource is used to create assignments between Entra Groups and System Roles in Tilgangsportalen

## Example Usage

```terraform
resource "tilgangsportalen_entra_group" "example" {
  name              = "group 1"
  description       = "Demo of Terraform created Microsoft Entra ID group"
  inheritance_level = "User" # or "Group"
}

resource "tilgangsportalen_system_role" "example" {
  name              = "published-role-name"
  product_category  = "TBD"    # product category of the system role. Must match an avaialable category
  system_role_owner = "m00001" # identity of the user who is the owner of the system role
  approval_level    = "L2"     # approval level of the system role
  description       = "Role for giving access to group_1 assigned resources."
  it_shop_name      = "General access shop shelf" # Optional. Defaults to "General access shop shelf"
}


resource "tilgangsportalen_entra_group_role_assignment" "example" {
  role_name   = tilgangsportalen_system_role.example.name
  entra_group = tilgangsportalen_entra_group.example.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entra_group` (String) The name of the Entra Group to assign to the Role
- `role_name` (String) The name of the Role to assign the Entra Group to

### Optional

- `force` (Boolean) Force the assignment even if it already exists

### Read-Only

- `id` (String) Identifier for the Entra Group System Role assignment. Currently, as we do not get a unique ID we can use from the API, ID is set by combining the role name and the Entra group name, with a pipe symbol as separator: RoleName|EntraGroupName

## Import

Import is supported using the following syntax:

```shell
terraform import tilgangsportalen_entra_group_role_assignment.example "published-role-name|group 1"
```

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import the resource. For example:

```terraform
import {
  id = "published-role-name|group 1"
  to = tilgangsportalen_entra_group_role_assignment.example
}
```