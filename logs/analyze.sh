#!/bin/bash

# Build the analyzer
echo "Building log analyzer..."
mkdir -p bin
go build -o bin/analyze cmd/analyze/main.go

# Determine log prefix based on environment
VM_TYPE=${VM_TYPE:-legacy} # Default to legacy if not set
echo "Using environment: $VM_TYPE"

# Define log file suffixes
LOG_SUFFIXES=("5_10.log" "10_50.log" "20_50.log" "100_10.log" "500_10.log")

# Create array of log files with appropriate prefix
LOG_FILES=()
for suffix in "${LOG_SUFFIXES[@]}"; do
    LOG_FILES+=("logs/${VM_TYPE}_${suffix}")
done

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
