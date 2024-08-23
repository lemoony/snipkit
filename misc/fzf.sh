#!/bin/bash

# Function to check if a command is installed
check_command() {
    if ! command -v "$1" &> /dev/null; then
        echo "Error: $1 is not installed."
        exit 1
    fi
}

# Check if either batcat or bat is installed (required for syntax highlighting)
if command -v batcat &> /dev/null; then
    BAT_CMD="batcat"
elif command -v bat &> /dev/null; then
    BAT_CMD="bat"
else
    echo "Error: bat or batcat is not installed."
    exit 1
fi

# Check if other required commands are installed
check_command "awk"
check_command "jq"
check_command "fzf"
check_command "snipkit"

# Directly call snipkit export to get JSON data
json_data=$(snipkit export -f=id,title,content)

# Extract titles, base64 encoded scripts, and IDs using jq
combined=$(echo "$json_data" | jq -r '.snippets[] | [.id, .title, (.content | @base64)] | @tsv')

# Use fzf to select a title and show the script in a preview window with syntax highlighting
selected=$(echo "$combined" | fzf --prompt="Select a script: " \
                                  --delimiter=$'\t' \
                                  --with-nth=2 \
                                  --preview="echo {} | awk -F'\t' '{print \$3}' | base64 --decode | $BAT_CMD --language=sh --style=numbers --color=always" \
                                  --preview-window=right:60%:wrap)

# Extract the ID and script part from the selected line
selected_id=$(echo "$selected" | awk -F'\t' '{print $1}')
selected_script=$(echo "$selected" | awk -F'\t' '{print $3}' | base64 --decode)

# If a script was selected, execute it using snipkit
if [ -n "$selected_id" ]; then
    snipkit exec --id "$selected_id"
else
    echo "No script selected."
fi
