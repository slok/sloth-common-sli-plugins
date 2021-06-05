#!/bin/bash
# vim: ai:ts=8:sw=8:noet
set -efCo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v sloth >/dev/null 2>&1 || { echo 'please install sloth'; exit 1; }

set +f # Allow asterisk expansion.

# Load all plugins and try generating SLOs without error, for now this is good enough 
# for integration tests along with each plugin unit tests.
for file in ./test/integration/*.yml; do
    fname=$(basename "$file")
    echo -n "[TEST] [${file}] Generating SLOs..."
    sloth generate -p ./plugins -i "${file}" --no-log > /dev/null
    echo ": OK"
done

set -f
