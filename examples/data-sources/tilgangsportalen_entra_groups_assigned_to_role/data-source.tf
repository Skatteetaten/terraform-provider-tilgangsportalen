# List all entra id groups assigned to a specific role
data "tilgangsportalen_entra_groups_assigned_to_role" "this" {
  role_name = "Role 123"
}
