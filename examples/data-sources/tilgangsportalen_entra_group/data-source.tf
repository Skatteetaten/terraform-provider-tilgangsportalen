# This will not expose entitlement_uid
data "tilgangsportalen_entra_group" "example" {
  name = "[Ex] group 1"
}

# This will not expose name
data "tilgangsportalen_entra_group" "example_with_entitlement_uid" {
  entitlement_uid = "00000000-0000-0000-0000-000000000000"
}

# This will expose both name and entitlement_uid
data "tilgangsportalen_entra_group" "example_with_name_and_entitlement_uid" {
  name            = "[Ex] group 1"
  entitlement_uid = "00000000-0000-0000-0000-000000000000"
}
