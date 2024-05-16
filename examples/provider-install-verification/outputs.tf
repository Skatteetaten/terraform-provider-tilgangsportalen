# output "all_roles" {
#   value       = data.tilgangsportalen_system_roles.all_roles
#   description = "Skriver ut navnet på alle roller eid av brukeren som brukes til identifisering mot Tilgangsportalen"
# }

# output "all_groups" {
#   value       = data.tilgangsportalen_entra_groups.all_groups
#   description = "Skriver ut navnet på alle Entra grupper eid av brukeren som brukes til identifisering mot Tilgangsportalen"
# }

# output "entra_role_assignment" {
#   value       = data.tilgangsportalen_entra_groups_assigned_to_role.role_assigned_to_group
#   description = "Skriver ut rollenavn og navn på Entra grupper knyttet til rollen i en role assignment"
# }

# output "entra_group_created" {
#   value       = tilgangsportalen_entra_group.example.name
#   description = "Skriver ut navnet på den opprettede gruppen"
# }

# output "system_role_created" {
#   value       = tilgangsportalen_system_role.test_role.name
#   description = "Skriver ut navnet på den opprettede gruppen"
# }
