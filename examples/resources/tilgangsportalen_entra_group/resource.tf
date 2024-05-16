resource "tilgangsportalen_entra_group" "example" {
  name              = "[Ex] group 1"
  alias             = "group_1"
  description       = "Demo of Terraform created Microsoft Entra ID group"
  inheritance_level = "User" # or "Admin"
}
