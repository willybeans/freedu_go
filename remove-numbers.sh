#!/bin/bash

# Check if a file name is provided as an argument
if [ $# -eq 0 ]; then
  echo "Usage: bash script.sh <filename>"
  exit 1
fi

# Check if the file exists
if [ ! -f "$1" ]; then
  echo "File '$1' not found."
  exit 1
fi

# Read the file content and remove numbers and periods,
# and replace all white spaces with a single comma
file_content=$(cat "$1" | sed -e 's/[0-9.]//g' -e 's/[[:space:]]\+/,/g')

# Output the modified content
echo "$file_content"

