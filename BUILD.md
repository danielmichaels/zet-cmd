# Build

When updating this repository to a new release the following steps need to be undertaken.

1. Update `Cmd.Version` in `zet-cmd/cmd/zet/main.go`
2. Merge to `main` all changes required for the new release
3. Create a new tag to match the `Cmd.Version` with `git tag -a v0.0.0 -m "v0.0.0" -s
4. Push the tag with `git push --tags` which will trigger `goreleaser` to build new binaries
