# Changelog

## 0.5.0

DEPRECATIONS:

- `tilgangsportalen_system_role` - `alias` has been deprecated and the field will be removed in a future release. The field can be safely removed.

ENHANCEMENTS:

- `tilgangsportalen_system_role` - allows DisplayName up to 256 characters.
- New data source `tilgangsportalen_system_role` that gets all details from API.
  NB. The API does not currently seem to return a value for `it_shop_name`

## 0.4.0

Breaking change:

- `tilgangsportalen_system_role` - only allow lowercase for fields `system_role_security_owner` and `system_role_owner`

BUG FIXES:

- `tilgangsportalen_system_role` - fixed import function to include missing ID field

ENHANCEMENTS:

- Add `ImportStateVerify` to import state tests
- On read `tilgangsportalen_system_role` update `system_role_owner` and `system_role_security_owner` in state
- Improved unit tests

tilgangsportalapi:

- `GetSystemRole`always returns lowercase `L2Ident` and `L3Ident`.

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
