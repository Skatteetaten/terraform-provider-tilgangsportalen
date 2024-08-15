# Changelog

## 0.3.1

BUG FIXES:

- `tilgangsportalen_system_role` - fixed import function to include missing ID field

ENHANCEMENTS:

- Changed data source unit test for `tilgangsportalen_entra_groups_assigned_to_role` to use unique Entra group name to avoid blocking resources

## 0.3.0

CHANGES:

- Rolled back changes to `it_shop_name` in `tilgangsportalen_system_role` that were introduced in v0.2.0. `it_shop_name` is no longer deprecated.
- `it_shop_name` in `tilgangsportalen_system_role` is now optional. Defaults to "General access shop shelf".
- `tilgangsportalen_entra_group` - Added length validation to the group description field. The max length is now set to 1024 characters.

BUG FIXES:

- Fixed a bug introduced in v0.2.0 where system roles were not being published after their creation.

## 0.2.0

WARNING:

This version of the provider contains a bug that causes system roles not to be published after their creation.
This bug has been fixed in v0.3.0 of the provider.

ENHANCEMENTS:

- `tilgangsportalen_system_role` now uses the new API endpoint `CreateRole`
- Minor improvements to docs and tests
- Removed regex validation from description field in `system_role` and `entra_group` resources

DEPRECATIONS:

- `it_shop_name` in `tilgangsportalen_system_role` is now deprecated and will
  be removed in a future release

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
