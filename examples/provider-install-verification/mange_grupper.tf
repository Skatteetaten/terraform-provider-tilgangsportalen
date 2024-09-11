# locals {
#   group_names = toset(["1", "2", "3", "4", "5"])
# }

# resource "tilgangsportalen_entra_group" "eksempel" {
#   for_each          = local.group_names
#   name              = "[Test] terraform-group-nummer-${each.value}"
#   inheritance_level = "User"
#   description       = "Group opprettet i fleng nummer ${each.value}"
# }
