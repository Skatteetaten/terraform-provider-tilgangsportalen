resource "tilgangsportalen_entra_group" "example" {
  name              = "group 1"
  alias             = "group_1"
  description       = "Demo of Terraform created Microsoft Entra ID group"
  inheritance_level = "User" # or "Group"
}

resource "tilgangsportalen_system_role" "example" {
  name              = "published-role-name"
  product_category  = "TBD"    # product category of the system role. Must match an avaialable category
  system_role_owner = "m00001" # identity of the user who is the owner of the system role
  approval_level    = "L2"     # approval level of the system role
  description       = "Role for giving access to group_1 assigned resources."
}


resource "tilgangsportalen_entra_group_role_assignment" "example" {
  role_name   = tilgangsportalen_system_role.example.name
  entra_group = tilgangsportalen_entra_group.example.name
}
