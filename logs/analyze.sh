#!/bin/bash

# Build the analyzer
echo "Building log analyzer..."
mkdir -p bin
go build -o bin/analyze cmd/analyze/main.go

# Determine log prefix based on VM_TYPEironment
# You can set VM_TYPE=tdx or VM_TYPE=legacy as an VM_TYPEironment variable
VM_TYPE=${VM_TYPE:-legacy} # Default to legacy if not set
echo "Using environment: $VM_TYPE"

# Define log file suffixes
LOG_SUFFIXES=("5_10_60.log" "10_50_180.log" "20_50_180.log" "100_10_180.log" "500_10_180.log")

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
