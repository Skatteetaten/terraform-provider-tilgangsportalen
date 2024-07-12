# Changelog

## 0.1.2

BUG FIXES:

- `tilgangsportalen_system_role` - Updated the user intent validation to
  require excactly 6 alphanumeric characters, as not all user identifiers follow
  the format `x00000` or `x00xxx`.
- `tilgangsportalen_entra_group` - Allow colons (:) and semicolons (;) in the description
- `tilgangsportalen_system_role` - Allow colons (:) and semicolons (;) in the description

## 0.1.1

ENHANCEMENTS:

- dependencies: Updated `google.golang.org/grpc` to version `v1.64.1` due to a
  vulnerability in the previous version

BUG FIXES:

- `tilgangsportalen_entra_group` - Allow underscores in the group name
- `tilgangsportalen_system_role` - Allow underscores in the role name

## 0.1.0

FEATURES:
