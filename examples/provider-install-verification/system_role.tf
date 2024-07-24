# resource "tilgangsportalen_system_role" "test_role" {
#   name                       = "test-rolle-1"
#   product_category           = "TBD"    # Value in test, in production use path to existing product category
#   system_role_owner          = "a00000" # Replace with valid user
#   system_role_security_owner = "b00000" # Replace with valid user
#   approval_level             = "L3"     # Valid choices are L2 and L3
#   description                = "Oppretter en rolle med L3 i tilgangsportalen for test via terraform provider."
#   it_shop_name               = "Access shop shelf"

#   lifecycle {
#     ignore_changes = [description, approval_level, product_category, system_role_owner, it_shop_name]
#   }
# }



# resource "tilgangsportalen_system_role" "test_role_2" {
#   name              = "test-rolle-2"
#   product_category  = "TBD"
#   system_role_owner = "a00000" # Replace with valid user
#   approval_level    = "L2"

# }
