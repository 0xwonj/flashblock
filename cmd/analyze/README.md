# Block Creation Time Analyzer

This is a utility to analyze block creation time statistics from FlashBlock log files. It calculates various statistical measures like mean, median, standard deviation, and percentiles, and also generates histograms to visualize the distribution of creation times.

## Features

- Calculate statistical metrics:
  - Minimum and maximum creation times
  - Mean and median creation times
  - Standard deviation
  - 95th and 99th percentiles
- Generate visual histograms of creation time distributions
- Group statistics by transaction count
- Save analysis results to file

## Usage

### Building

```bash
go build -o analyze main.go
```

### Running

```bash
# Basic usage (output to stdout)
./analyze -log path/to/log/file.log

# Save results to file
./analyze -log path/to/log/file.log -output analysis_results.txt
```

### Using the analyze_logs.sh Script

For convenience, you can use the provided shell script:

```bash
./analyze_logs.sh
```

This will:
1. Build the analyzer
2. Run it on the log file
3. Save the results to a timestamped file
4. Display the results

## Output Example

```
Block Creation Time Statistics (in microseconds):
Total blocks analyzed: 721
Min: 28.081 µs
Max: 507.286 µs
Mean: 72.942 µs
Median: 57.651 µs
Standard Deviation: 56.253 µs
95th Percentile: 206.324 µs
99th Percentile: 334.955 µs

Creation Time Distribution (µs):
   28.1 -    76.0 µs |  628 | ██████████████████████████████████████████████████
   76.0 -   123.9 µs |   41 | ███
  123.9 -   171.8 µs |   10 | 
  171.8 -   219.8 µs |   10 | 
...
``` 