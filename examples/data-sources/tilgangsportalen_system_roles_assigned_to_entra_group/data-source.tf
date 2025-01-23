# List all system roles a specific entra id group is assigned to
data "tilgangsportalen_system_roles_assigned_to_entra_group" "this" {
  group_name = "[Ex] group 1"
}
