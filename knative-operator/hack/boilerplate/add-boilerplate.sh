#!/bin/bash

USAGE=$(cat <<EOF
Add boilerplate.<ext>.txt to all .<ext> files missing it in a directory.

Usage: (from repository root)
       ./hack/boilerplate/add-boilerplate.sh <ext> <DIR>

Example: (from repository root)
         ./hack/boilerplate/add-boilerplate.sh go cmd
EOF
)

set -e

if [[ -z $1 || -z $2 ]]; then
  echo "${USAGE}"
  exit 1
fi

grep -r -L -P "Copyright \d+ The Knative Authors" $2  \
  | grep -P "\.$1\$" \
  | xargs -I {} sh -c \
  "cat hack/boilerplate/boilerplate.$1.txt {} > /tmp/boilerplate && mv /tmp/boilerplate {}"
