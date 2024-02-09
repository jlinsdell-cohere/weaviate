#!/bin/bash

# Function to display help message
display_help() {
    echo "Usage: $0 [VOTERS]"
    echo
    echo "Generate a docker-compose file for a Weaviate raft cluster."
    echo
    echo "VOTERS is the number of voters in the raft cluster. Default is 2."
    echo
    echo "Options:"
    echo "  --help    Display this help message."
}

# Check command-line arguments
case "$1" in
    --help)
        display_help
        exit 0
        ;;
esac

# Set default value for VOTERS
VOTERS=${1:-2}
FILE_NAME="docker-compose-raft.yml"
CURRENT_DIR=$(dirname "$0")

# Check if jinja2 is installed
if ! command -v jinja2 &> /dev/null; then
    # Check if running on Mac
    echo "Installing jinja2-cli..."
    if [[ $(uname) == "Darwin" ]]; then
        # Install jinja2-cli via Homebrew
        brew install jinja2-cli
    else
        pip install jinja2-cli
    fi
fi

# Generate docker-compose-raft.yml using jinja2
jinja2 ${CURRENT_DIR}/docker-compose-raft.yml.j2 -D NUMBER_VOTERS=${VOTERS} -o ${FILE_NAME}

echo -e "You can now start your multinode Weaviate compose! To do so, run the following command:\n\
    docker-compose -f ${FILE_NAME} up -d\n\
This command will start $((VOTERS + 1)) nodes.\n\
If you want to start more nodes, for example 10 nodes. You can use the following command:\n\
    docker-compose -f ${FILE_NAME} up -d --scale weaviate=$((10 - VOTERS))\n\
To stop the nodes, you can use the following command:\n\
    docker-compose -f ${FILE_NAME} down\n"
