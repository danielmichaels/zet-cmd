# Build

When updating this repository to a new release the following steps need to be undertaken.

1. Merge to `main` all changes required for the new release
2. Create a new with `git tag -a v0.0.0 -m "v0.0.0" -s`
3. Push the tag with `git push --tags` which will trigger `goreleaser` to build new binaries
