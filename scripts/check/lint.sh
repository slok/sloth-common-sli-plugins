#!/usr/bin/env sh

set -o errexit
set -o nounset

# Lint.
golangci-lint run
