# Changelog

## v0.2.0

### Breaking changes

- Removed `Server.Somewhat`; use `Server.Passworded` (`bool`) instead.
- Added `Server.PlayerList []Player` with `Name`, `Score`, and `Ping`.
- `ReadRow` now returns `error` instead of `*Server`.
- `UpdateOnce` and `ReadRow` return parse errors instead of silently ignoring them.

### Added

- Explicit ASE (`EYE1`) parser with bounds checking and wrapped errors.
- Godoc package documentation and `Example*` tests for pkg.go.dev.

## v0.1.0

- Initial release with reflect-based ASE parsing.
