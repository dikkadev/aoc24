#!/bin/bash

# Exit on any error
set -euo pipefail

LOG_FILE="create_day.log"
TEMPLATE_FILE="days/template.go"
BUILD_EXCLUDE="// +build exclude"

log() {
    # echo "$(date +'%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
    echo "$(date +'%Y-%m-%d %H:%M:%S') - $1"
}

usage() {
    echo "Usage: $0 [-gold] <day>"
    exit 1
}

if [[ $# -lt 1 || $# -gt 2 ]]; then
    usage
fi

GOLD_MODE=false
if [[ "$1" == "-gold" ]]; then
    GOLD_MODE=true
    shift
fi

INPUT="$1"
DAY_DIR="days/$INPUT"
INPUT_FILE="$DAY_DIR/$INPUT.go"
GOLD_FILE="$DAY_DIR/${INPUT}_gold.go"
INPUT_DIR="input"
INPUT_INPUT_FILE="$INPUT_DIR/$INPUT.input"
INPUT_SMALL_FILE="$INPUT_DIR/${INPUT}_small.input"

if $GOLD_MODE; then
    # Ensure the base script has already been run for the given input
    if [[ ! -f "$INPUT_FILE" ]]; then
        log "ERROR: Base setup for day '$INPUT' is not complete. Run the script without -gold first."
        exit 1
    fi

    # Copy the INPUT.go to INPUT_gold.go
    if [[ ! -f "$GOLD_FILE" ]]; then
        cp "$INPUT_FILE" "$GOLD_FILE"
        log "Created gold file: $GOLD_FILE"
    else
        log "Gold file already exists: $GOLD_FILE"
    fi

    # Add build exclude line to the top of INPUT.go if not already present
    if ! grep -q "$BUILD_EXCLUDE" "$INPUT_FILE"; then
        sed -i "1s|^|$BUILD_EXCLUDE\n\n|" "$INPUT_FILE"
        log "Added build exclude line to: $INPUT_FILE"
    else
        log "Build exclude line already present in: $INPUT_FILE"
    fi

    log "Gold mode completed successfully for day: $INPUT"

    # *** Begin Git Commit Section ***
    # Stage the original .go file and both input files
    git add "$INPUT_FILE" "$INPUT_INPUT_FILE" "$INPUT_SMALL_FILE"

    # Format the day number to ensure it's two digits (e.g., '1' becomes '01')
    DAY_NUM=$(printf "%02d" "$INPUT")

    # Create the commit message
    COMMIT_MSG="AUTO: day $DAY_NUM regular"

    # Commit the changes with the formatted message
    git commit -m "$COMMIT_MSG"

    log "Committed changes with message: '$COMMIT_MSG'"
    # *** End Git Commit Section ***

    exit 0
fi

log "Script started for day: $INPUT"

# Create day directory
if [[ ! -d "$DAY_DIR" ]]; then
    mkdir -p "$DAY_DIR"
    log "Created directory: $DAY_DIR"
else
    log "Directory already exists: $DAY_DIR"
fi

# Create INPUT.go file from template
if [[ -f "$TEMPLATE_FILE" ]]; then
    if [[ ! -f "$INPUT_FILE" ]]; then
        tail -n +2 "$TEMPLATE_FILE" > "$INPUT_FILE"
        log "Created file: $INPUT_FILE from template (without first line)"
    else
        log "File already exists: $INPUT_FILE"
    fi
else
    log "ERROR: Template file not found: $TEMPLATE_FILE"
    exit 1
fi

# Create input directory
if [[ ! -d "$INPUT_DIR" ]]; then
    mkdir -p "$INPUT_DIR"
    log "Created directory: $INPUT_DIR"
else
    log "Directory already exists: $INPUT_DIR"
fi

# Create empty input files
for FILE in "$INPUT_INPUT_FILE" "$INPUT_SMALL_FILE"; do
    if [[ ! -f "$FILE" ]]; then
        touch "$FILE"
        log "Created empty file: $FILE"
    else
        log "File already exists: $FILE"
    fi
done

log "Script completed successfully for day: $INPUT"
