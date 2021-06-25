#!/bin/bash
# vim: ai:ts=8:sw=8:noet
set -efCo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v sloth >/dev/null 2>&1 || { echo 'please install sloth'; exit 1; }

sloth validate -p ./plugins -i ./test/integration/ --debug