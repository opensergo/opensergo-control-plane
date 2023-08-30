#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

#
# This script builds the required plugins.
set -e

GOOS=$(go env GOOS)
echo "${GOOS}"
BINARY_SUFFIX=""
if [ "${GOOS}" = "windows" ]; then
    BINARY_SUFFIX=".exe"
fi

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
export DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"
echo "DIR: $DIR"
# Create boundary plugins
echo "==> Building opensergo plugins..."
rm -f $DIR/assets/opensergo-plugin-*
for CURR_PLUGIN in $(ls $DIR/server); do
    cd $DIR/server/$CURR_PLUGIN;
    go build -o $DIR/assets/opensergo-plugin-${CURR_PLUGIN}${BINARY_SUFFIX} .;
    cd $DIR;
done;
#cd $DIR/assets;
#for CURR_PLUGIN in $(ls opensergo-plugin*); do
#    gzip -f -9 $CURR_PLUGIN;
#done;
echo "==> Building opensergo plugins... DONE"