memo .envはvpsに保管

Docker image (GHCR):
- `ghcr.io/oudentabetai/pterodactyl-go:latest`
- `ghcr.io/oudentabetai/pterodactyl-go:<release-tag>` (例: `v1.0.5`)

## Repository merger notes

- Origin repositories:
  - `oudentabetai/pterodactyl-go` (base)
  - `oudentabetai/TwitterLinkfixer-go` (merged source for the requested linkfixer repo)
- Merge approach: unrelated-histories merge with `TwitterLinkfixer-go` imported under `twitterlinkfixer-go/` to avoid root-level collisions and preserve history.
- Module structure: multi-module layout.
  - Root module: `github.com/oudentabetai/pterodactyl-go`
  - Imported module: `twitterlinkfixer-go/go.mod` (`github.com/oudentabetai/twitterlinkfixer-go`)
- Breaking path/module changes:
  - Files from the imported repository now live under `twitterlinkfixer-go/`.
  - Any previous root-relative paths for the imported project must be updated to the new subdirectory path.
