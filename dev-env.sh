#!/usr/bin/env bash

###############################################################################
# This file is version-controlled and shared across developers.
#
# It is used for setting up a developer environment for local testing.
#
# Additionally overrides can be specified in dev-env.local.sh or dev.local.env.
# Local environment variables may be added in dev.local.env.
#
# To successfully load the environment variables this file needs to be sourced.
#
# Do not commit secrets or machine-specific values.
###############################################################################

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
END='\033[0m'

cecho() {
    echo -e "$*"
}

# Exit if not sourced
(return 0 2>/dev/null) || {
    cecho "${RED}ERROR${END} dev-env.sh should be sourced, not executed."
    cecho "${BLUE}Usage:${END} source dev-env.sh"
    exit 1
}

# Get directory path
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export DEFAULTS_PATH=$DIR/data/defaults.yml
export FAVICON_PATH=$DIR/data/favicon.ico
export CONFIG_PATH=$DIR/.dev/config.yml
export TOKENS_PATH=$DIR/.dev/tokens

export LOG_LEVEL=dev

# Load overrides if present
if [ -f "$DIR/dev-env.local.sh" ]; then
    cecho "${BLUE}Found overrides:${END} $DIR/dev-env.local.sh"
    source "$DIR/dev-env.local.sh"
fi
if [ -f "$DIR/dev.local.env" ]; then
    cecho "${BLUE}Found overrides:${END} $DIR/dev.local.env"
    source "$DIR/dev.local.env"
fi

cecho "${GREEN}Successfully loaded development environment!${END}"