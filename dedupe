#!/bin/bash

# Function to show help message
show_help() {
    echo "Usage: $0 [options] [<input_file>]"
    echo "Options:"
    echo "  -h, --help                    Usage/help info for dedupe"
    echo "  -u, --urls <filename>         Filename containing URLs (use this if you don't pipe URLs via stdin)"
    echo "  -V, --version                 Get current version for dedupe"
    echo "  -qs, --query-strings-only     Only include URLs if they have query strings"
    echo "  -ne, --no-extensions <ext1,ext2,...>"
    echo "                                Do not include URLs with specific extensions"
}

# Function to show version
show_version() {
    echo "dedupe_v0.3"
}

# Function to show banner
show_banner() {
    echo "dedupe"
    echo "______  _______ ______  _     _  _____  _______"
    echo "|     \ |______ |     \ |     | |_____] |______"
    echo "|_____/ |______ |_____/ |_____| |       |______ v0.3"
    echo "                                      @0xPugal"
}

# Function to normalize a URL by hostname and sorted query parameter names
normalize_url() {
    local url="$1"
    local hostname=$(echo "$url" | awk -F '/' '{print $3}')
    local query_string=$(echo "$url" | awk -F '?' '{print $2}')

    if [ -n "$query_string" ]; then
        sorted_param_names=$(echo "$query_string" | tr '&' '\n' | cut -d '=' -f1 | sort | tr '\n' '&')
        sorted_param_names="${sorted_param_names%&}"
        echo "${hostname}?${sorted_param_names}"
    else
        echo "${hostname}"
    fi
}

# Function to deduplicate URLs based on the normalized representation
deduplicate_urls() {
    local input="$1"
    local query_strings_only="$2"
    local no_extensions="$3"
    local seen_urls=()

    while IFS= read -r url; do
        if [ -z "$url" ]; then
            continue
        fi

        # Check if the URL has any of the specified extensions
        local has_extension=false
        for ext in ${no_extensions//,/ }; do
            if [[ "$url" == *".$ext"* ]]; then
                has_extension=true
                break
            fi
        done

        if [ "$has_extension" == true ]; then
            continue
        fi

        if [[ "$query_strings_only" == true && "$url" != *\?* ]]; then
            continue
        fi

        normalized_url=$(normalize_url "$url")

        if ! [[ " ${seen_urls[@]} " =~ " $normalized_url " ]]; then
            seen_urls+=("$normalized_url")
            echo "$url"
        fi
    done < "$input"
}

# Parse command-line arguments
input_file=""
query_strings_only=false
no_extensions=""
show_help_only=false

# Show banner if no arguments are provided
if [[ "$#" -eq 0 ]]; then
    show_banner
    exit 0
fi

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            show_help_only=true
            shift
            ;;
        -u|--urls)
            input_file="$2"
            shift 2
            ;;
        -V|--version)
            show_version
            exit 0
            ;;
        -qs|--query-strings-only)
            query_strings_only=true
            shift
            ;;
        -ne|--no-extensions)
            no_extensions="$2"
            shift 2
            ;;
        *)
            if [ -z "$input_file" ]; then
                input_file="$1"
            else
                echo "Unknown option: $1"
                show_help
                exit 1
            fi
            shift
            ;;
    esac
done

# Show help if -h or --help is specified
if [[ "$show_help_only" == true ]]; then
    show_help
    exit 0
fi

# If no input file is provided, use stdin
if [ -z "$input_file" ]; then
    input_file="/dev/stdin"
fi

# Execute deduplication based on provided options
deduplicate_urls "$input_file" "$query_strings_only" "$no_extensions"
