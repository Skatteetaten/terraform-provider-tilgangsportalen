# Get system role by name
data "tilgangsportalen_system_role" "by_name" {
  name = "System Role Name"
}

# Get system role by object_id
data "tilgangsportalen_system_role" "by_id" {
  object_id = "12345678-1234-1234-1234-123456789abc"
}
