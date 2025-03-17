#!/bin/bash

# Build the analyzer
echo "Building log analyzer..."
mkdir -p bin
go build -o bin/analyze cmd/analyze/main.go

# Define an array of log files
LOG_FILES=("logs/legacy_5_10_60.log" "logs/legacy_10_50_180.log" "logs/legacy_20_50_180.log" "logs/legacy_100_10_180.log" "logs/legacy_500_10_180.log")

# Process each log file in the array
for LOG_FILE in "${LOG_FILES[@]}"; do
    # Create output filename based on log file name and timestamp
    OUTPUT_FILE="${LOG_FILE%.*}_analysis.log"

    # Run the analyzer
    echo "Analyzing log file: $LOG_FILE..."
    ./bin/analyze -log "$LOG_FILE" -output "$OUTPUT_FILE"

    echo "Analysis complete and saved to $OUTPUT_FILE"
    echo "Displaying results..."
    cat "$OUTPUT_FILE"
    echo "----------------------------------------"
done 