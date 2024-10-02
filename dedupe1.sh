#!/bin/bash

# Function to normalize a URL by hostname and sorted query parameter names
normalize_url() {
    local url="$1"
    # Extract the hostname
    local hostname=$(echo "$url" | awk -F '/' '{print $3}')
    # Extract the query string
    local query_string=$(echo "$url" | awk -F '?' '{print $2}')

    if [ -n "$query_string" ]; then
        # Extract and sort the parameter names
        sorted_param_names=$(echo "$query_string" | tr '&' '\n' | cut -d '=' -f1 | sort | tr '\n' '&')
        sorted_param_names="${sorted_param_names%&}"  # Remove trailing '&'
        # Return the unique representation with hostname and sorted parameter names
        echo "${hostname}?${sorted_param_names}"
    else
        # If there's no query string, return just the hostname
        echo "${hostname}"
    fi
}

# Function to deduplicate URLs based on the normalized representation
deduplicate_urls() {
    local input_file="$1"
    declare -A unique_urls
    declare -A original_urls

    # Read each URL from the input file
    while IFS= read -r url; do
        # Skip empty or invalid lines
        if [ -z "$url" ]; then
            continue
        fi

        # Normalize the URL to get a unique representation based on hostname and sorted parameter names
        normalized_url=$(normalize_url "$url")

        # Deduplicate based on normalized representation
        if [ -z "${unique_urls["$normalized_url"]}" ]; then
            unique_urls["$normalized_url"]=1
            original_urls["$normalized_url"]="$url"
        fi
    done < "$input_file"

    # Output the unique original URLs
    for normalized in "${!original_urls[@]}"; do
        echo "${original_urls["$normalized"]}"
    done
}

# Ensure a valid input file is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <input_file>"
    exit 1
fi

# Deduplicate the URLs from the provided input file
deduplicate_urls "$1"
