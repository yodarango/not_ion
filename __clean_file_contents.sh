#!/bin/bash

# Function to clean HTML file content
clean_html_file() {
    local html_file="$1"
    echo "🧹 Cleaning HTML file: $html_file"
    
    # Create a backup only for the index.html file
    if [ "$html_file" == "$dir/index.html" ]; then
        cp "$html_file" "${html_file}.backup"
    fi  

    
    # Use perl for more precise pattern matching
    perl -i.bak -pe '
        # Clean id attributes
        s/id::[a-zA-Z0-9-]{32,}/id::cleaned/g;
        
        # Clean href attributes - look for segments that:
        # 1. Are preceded by space or /
        # 2. Are at least 12 chars long
        # 3. Contain at least 3 digits
        # 4. Are followed by / or .html or space or quote
        s/(\s|\/)[a-zA-Z0-9]{12,}(?=\/|\.html|\s|")/\1/g;
        
        # Clean any remaining long alphanumeric strings before .html
        s/\s+[a-zA-Z0-9]{12,}\.html/\.html/g;
        
        # Clean up trailing/multiple spaces in href attributes
        s/href="([^"]*?)"/sub {
            my $url = $1;
            $url =~ s/\s+(?=\/|\.html)//g;  # Remove spaces before / or .html
            $url =~ s/\s+/ /g;              # Convert multiple spaces to single space
            "href=\"$url\""
        }/ge;
        
        # Clean up displayed text between tags
        s/>([^<]*?\.html)</sub {
            my $text = $1;
            $text =~ s/\s+(?=\.html|$)//g;  # Remove spaces before .html or end
            $text =~ s/\s+/ /g;             # Convert multiple spaces to single space
            ">$text<"
        }/ge;
    ' "$html_file"

    echo "✅ Done cleaning HTML file: $html_file"
}

# Process all HTML files in the directory and subdirectories
echo "💉 Processing HTML files..."
find "$dir" -type f -name "*.html" -print0 | while IFS= read -r -d '' html_file; do
    # Skip backup files
    if [[ "$html_file" != *".backup" && "$html_file" != *".bak" ]]; then
        clean_html_file "$html_file"
    fi
done

# Process all files for renaming
echo "🥩 Processing files for renaming..."
find "$dir" -type f -print0 | while IFS= read -r -d '' file; do
    # Skip the backup files
    if [[ "$file" != *".backup" && "$file" != *".bak" ]]; then
        clean_filename "$file"
    fi
done

# Then process all directories (from deepest to shallowest)
echo "📁 Processing directories..."
find "$dir" -type d -print0 | sort -rz | while IFS= read -r -d '' directory; do
    # Skip the root directory
    if [ "$directory" != "$dir" ]; then
        clean_filename "$directory"
    fi
done

# Clean up any remaining .bak files created by sed/perl
find "$dir" -name "*.bak" -type f -delete

echo "🫧 Cleanup completed!"

# Define the yellow color code
YELLOW='\033[1;33m'
# Define the reset color code
RESET='\033[0m'

# Print the message in yellow
echo -e "⭐️🚨🟡 ${YELLOW}Please check all HTML files. You might have to replace \"/\" and \"  \" chars in names manually${RESET}"