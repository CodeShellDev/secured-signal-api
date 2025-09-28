#!/bin/bash

# Usage: ./replace_placeholders.sh source_file destination_file

template() {
    local content="$1"
    
    # Extract unique placeholders: match { { file.xxx } } with arbitrary spaces
    local placeholders
    placeholders=$(printf '%s' "$content" | grep -oP '\{\s*{\s*file\.(.*)\s*}\s*}' | sed -E 's/\{\s*\{\s*file\.//;s/\s*\}\s*\}//' | sort -u)

    for placeholder in $placeholders; do
        local file_content

        if [[ -f "$placeholder" ]]; then
            file_content=$(<"$placeholder")
        
            if [[ "$placeholder" == *.template.* ]]; then
                # Template further
                file_content=$(template "$file_content")
            fi
        else
            file_content="File not found: '$placeholder'."
        fi

        # Escape special characters
        local escaped_content
        local escaped_placeholder
        escaped_content=$(printf '%s' "$file_content" | perl -pe 's/([\\\/\$])/\\$1/g; s/\n/\\n/g;')
        escaped_placeholder=$(printf '%s' "$placeholder" | perl -pe 's/([\\\/])/\\$1/g; s/\n/\\n/g;')
        
        content=$(printf '%s' "$content" | perl -pe "s/{\s*{\s*file\.${escaped_placeholder}\s*}\s*}/$escaped_content/g")
    done

    printf '%s' "$content"
}

write_file() {
    local file_content="$1"
    local dest_file="$2"

    # Write to destination file
    printf '%s\n' "$file_content" > "$dest_file"
}

template_file() {
    local source_file="$1"
    local dest_file="$2"

    # Read the entire source file into a variable
    local file_content
    file_content=$(<"$source_file")

    local templated_content
    templated_content=$(template "$file_content")

    write_file "$templated_content" "$dest_file"

    echo "'$source_file' complete. Output written to '$dest_file'."
}

SOURCE="$1"
DEST_FILE="$2"

if [[ -f "$SOURCE" ]]; then
    template_file "$SOURCE_FILE" "$DEST_FILE"
elif [[ -d "$SOURCE" ]]; then
    i=0
    while IFS= read -r source_file; do
        file_content=$(<"$source_file")

        first_line=$(printf '%s\n' "$file_content" | head -n 1)

        dest_file=$(echo "$first_line" | sed -n 's/^[[:space:]]*>>\([[:graph:]]\+\)[[:space:]]*$/\1/p')

        if [[ -n "$dest_file" ]]; then
            file_content=$(printf '%s' "$file_content" | tail -n +2)
        
            templated_content=$(template "$file_content")

            write_file "$templated_content" "$dest_file"

            echo "'$source_file' complete. Output written to '$dest_file'."
        
            (( i++ ))
        fi
    done < <(find "$SOURCE" -type f -name '*.template.md')

    if [[ $i -eq 0 ]]; then
        echo "Source '$SOURCE' does not contain any '*.template.md'!"
        exit 1
    fi
else
    echo "Source '$SOURCE' does not exist!"
    exit 1
fi
