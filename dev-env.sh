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

#=-----------------------------------=#
#= Environment Variables & Overrides =#
#=-----------------------------------=#

export DEFAULTS_PATH=$DIR/data/defaults.yml
export FAVICON_PATH=$DIR/data/favicon.ico
export CONFIG_PATH=$DIR/.dev/config.yml
export TOKENS_DIR=$DIR/.dev/tokens
export DB_PATH=$DIR/.dev/db/db.sqlite3

export CGO_ENABLED=1

export API_URL=http://127.0.0.1:8881

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

#=-----------------------------------=#
#=            Mock server            =#
#=-----------------------------------=#

MOCK_PORT="8881"

MOCK_BIN="/tmp/mockserver-$MOCK_PORT"
MOCK_PID="/tmp/mockserver-$MOCK_PORT.pid"

# Kill mockserver if still running
if [ -f "$MOCK_PID" ]; then
    OLD_PID=$(cat "$MOCK_PID")
    if ps -p "$OLD_PID" > /dev/null 2>&1; then
        cecho "${YELLOW}Stopping previous Mock server (PID $OLD_PID)${END}"
        kill "$OLD_PID"
        sleep 1
    fi
fi

# Build mockserver
go build -o "$MOCK_BIN" utils/mockserver/mockserver.go

set +m

if ! nc -z 127.0.0.1 "$MOCK_PORT" >/dev/null 2>&1; then
    $MOCK_BIN > "$DIR/.dev/mock.log" 2>&1 &
    NEW_PID=$!
    disown

    echo "$NEW_PID" > "$MOCK_PID"

    echo "Mock server started at http://127.0.0.1:$MOCK_PORT"
else
    echo "There is already a Service running at http://127.0.0.1:$MOCK_PORT"
fi