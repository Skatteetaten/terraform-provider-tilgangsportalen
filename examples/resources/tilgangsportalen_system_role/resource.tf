resource "tilgangsportalen_system_role" "example" {
  name              = "Published role name"
  product_category  = "TBD"    # product category of the system role. Must match an avaialable category
  system_role_owner = "a00000" # identity of the user who is the owner of the system role
  approval_level    = "L2"     # approval level of the system role
  description       = "Role for giving access to xyz."
  it_shop_name      = "Access shop shelf"
}
