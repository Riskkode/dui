#!/bin/bash

# Check if the current directory is a git repository
if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Error: Current directory is not a git repository."
    exit 1
fi

# Get the remote URL (defaults to origin)
REMOTE_NAME=${1:-origin}
REMOTE_URL=$(git remote get-url "$REMOTE_NAME" 2>/dev/null)

if [ -z "$REMOTE_URL" ]; then
    echo "Error: No remote named '$REMOTE_NAME' found in this repository."
    exit 1
fi

# Convert SSH URL to HTTPS if needed
if [[ $REMOTE_URL == git@* ]]; then
    # Convert SSH format (git@github.com:username/repo.git) to HTTPS format
    REMOTE_URL=$(echo "$REMOTE_URL" | sed -E 's|git@([^:]+):|https://\1/|g')
fi

# Remove .git extension if present
REMOTE_URL=$(echo "$REMOTE_URL" | sed 's/\.git$//')

echo "Opening $REMOTE_URL in your browser..."

# Open the URL in the default browser based on the operating system
case "$(uname -s)" in
    Linux*)     xdg-open "$REMOTE_URL";;
    Darwin*)    open "$REMOTE_URL";;
    CYGWIN*|MINGW*|MSYS*)    start "$REMOTE_URL";;
    *)          echo "Unsupported operating system. Please open this URL manually: $REMOTE_URL";;
esac