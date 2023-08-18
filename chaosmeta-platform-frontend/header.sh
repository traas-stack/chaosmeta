#!/bin/bash
# Get the directory where the Go files are located
DIR="$1"
# Set the file header
HEADER="/*\n"
HEADER="$HEADER * Copyright 2022-2023 Chaos Meta Authors.\n"
HEADER="$HEADER *\n"
HEADER="$HEADER * Licensed under the Apache License, Version 2.0 (the \"License\");\n"
HEADER="$HEADER * you may not use this file except in compliance with the License.\n"
HEADER="$HEADER * You may obtain a copy of the License at\n"
HEADER="$HEADER *\n"
HEADER="$HEADER *  http://www.apache.org/licenses/LICENSE-2.0\n"
HEADER="$HEADER *\n"
HEADER="$HEADER * Unless required by applicable law or agreed to in writing, software\n"
HEADER="$HEADER * distributed under the License is distributed on an \"AS IS\" BASIS,\n"
HEADER="$HEADER * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n"
HEADER="$HEADER * See the License for the specific language governing permissions and\n"
HEADER="$HEADER * limitations under the License.\n"
HEADER="$HEADER */\n\n"
# Define a function to add the header to a file
function add_header {
  # Check if the file already has the header
  if ! grep -q "$HEADER" "$1"; then
    # Add the header to the file
    sed -i "" "1i $HEADER" "$1"
  fi
}
# Loop through all the Go files in the directory
find "$DIR" -type f -name "*.ts" | while read file; do
# for file in $(find "$DIR" -type f \( -name "*.js" -o -name "*.ts" \)); do    # Check if the file has a package declaration

  # Check if the file has a package declaration
  if ! grep -q "^package " "$file"; then
    # Add a blank line and the package declaration to the file
    sed -i "1i\npackage main" "$file"
  fi
  # Add the file header to the file
  add_header "$file"
done