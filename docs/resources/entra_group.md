---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tilgangsportalen_entra_group Resource - tilgangsportalen"
subcategory: ""
description: |-
  This resource is used to create a new Entra Group using Tilgangsportalen
---

# tilgangsportalen_entra_group (Resource)

This resource is used to create a new Entra Group using Tilgangsportalen

## Example Usage

```terraform
resource "tilgangsportalen_entra_group" "example" {
  name              = "[Ex] group 1"
  description       = "Demo of Terraform created Microsoft Entra ID group"
  inheritance_level = "User" # or "Admin"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `inheritance_level` (String) The inheritance level of the Entra Group (User or Admin). Determines what type of AD account the group can be assigned to.
- `name` (String) The display name of the Entra Group. Must be unique. Please follow the standardized naming conventions for Entra ID groups.

### Optional

- `alias` (String, Deprecated) Alias for the Entra Group. Deprecated and no longer in use.
- `description` (String) A description of the Entra Group

### Read-Only

- `id` (String) Identifier for the Entra Group. Currently, as we do not get a unique ID we can use from the API, ID is set equal to DisplayName

## Import

Import is supported using the following syntax:

```shell
terraform import tilgangsportalen_entra_group.example "[Ex] group 1"
```

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import the resource. For example:

```terraform
import {
  id = "[Ex] group 1"
  to = tilgangsportalen_entra_group.example
}
```