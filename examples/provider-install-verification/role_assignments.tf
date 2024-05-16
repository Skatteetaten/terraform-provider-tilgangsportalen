# resource "tilgangsportalen_entra_group_role_assignment" "example" {
#   role_name   = tilgangsportalen_system_role.test_role.name
#   entra_group = tilgangsportalen_entra_group.example.name

#   lifecycle {
#     replace_triggered_by = [tilgangsportalen_entra_group.example, tilgangsportalen_system_role.test_role]
#   }

# }
