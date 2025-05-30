#!/bin/bash

# Usage: ./replace_placeholders.sh source_file destination_file

SOURCE_FILE="$1"
DEST_FILE="$2"

# Check if source file exists
if [[ ! -f "$SOURCE_FILE" ]]; then
    echo "Source file '$SOURCE_FILE' does not exist!"
    exit 1
fi

# Read the entire source file into a variable
content=$(<"$SOURCE_FILE")

# Extract unique placeholders: match { { file.xxx } } with arbitrary spaces
placeholders=$(grep -oP '\{\s*{\s*file\.(.*)\s*}\s*}' "$SOURCE_FILE" | sed -E 's/\{\s*\{\s*file\.//;s/\s*\}\s*\}//' | sort -u)

for placeholder in $placeholders; do
    if [[ -f "$placeholder" ]]; then
        file_content=$(<"$placeholder")
        
        # Escape special characters
        escaped_content=$(printf '%s' "$file_content" | perl -pe 's/([\\\/])/\\$1/g; s/\n/\\n/g;')
        escaped_placeholder=$(printf '%s' "$placeholder" | perl -pe 's/([\\\/])/\\$1/g; s/\n/\\n/g;')
        
        content=$(printf '%s' "$content" | perl -pe "s/{\s*{\s*file\.${escaped_placeholder}\s*}\s*}/$escaped_content/g")
    else
        echo "Warning: File '$placeholder' not found, skipping."
    fi
done

# Write to destination file
printf '%b\n' "$content" > "$DEST_FILE"

echo "Replacements complete. Output written to '$DEST_FILE'."
