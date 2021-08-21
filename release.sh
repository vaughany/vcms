#!/bin/bash

# Test it.
# goreleaser build --skip-validate --rm-dist
# goreleaser release --snapshot --skip-publish --rm-dist

# Release it.
#   --skip-validate ignores git's current state.
#   --rm-dist clears out the ./dist folder first.
goreleaser release --skip-validate --rm-dist