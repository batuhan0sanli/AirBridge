# Changelog

All notable changes to this project will be documented in this file.

## [v0.2.0]

### Added
- **CI/CD**: Added GitHub Actions workflows for testing and linting.
- **CLI**: Added support for file input (`-i`) and output (`-o`) flags.
- **CLI**: Added headless mode (`--headless`, `-H`) for non-interactive use.
- **CLI**: Added key generation command (`keygen`) and functionality.
- **CLI**: Added support for public key file input (`--pubkey`).
- **Security**: Switched encryption from AES-CTR to **AES-GCM** for authenticated encryption.
- **Security**: Enabled private key loading (`--privkey`) for decryption.
- **Documentation**: Added documentation for key generation and usage.
- **Testing**: Added headless mode tests.

### Changed
- Refactors variable declaration and defers file close.
- Updates golangci-lint version.

## [v0.1.3]

### Fixed
- **Security**: Sanitize filename to prevent path traversal vulnerabilities (#1).

## [v0.1.2]

### Added
- Adds workflow status badge to README.
- Adds animation and updates images in README.

### Changed
- Updates Scoop bucket URL.
- Sets up release configurations.

## [v0.1.1]

### Fixed
- Add `GH_PAT` env var for goreleaser.

## [v0.1.0]

### Added
- Initial release.
- Added goreleaser config and workflow.
- Basic send/receive functionality with TUI.
