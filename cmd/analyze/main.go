package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// Parse command line arguments
	logFilePath := flag.String("log", "", "Path to the log file")
	outputFilePath := flag.String("output", "", "Path to save results (if empty, results are printed to stdout)")
	flag.Parse()

	if *logFilePath == "" {
		log.Fatal("Please provide a log file path using the -log flag")
	}

	// Setup output - either file or stdout
	var output io.Writer = os.Stdout
	if *outputFilePath != "" {
		file, err := os.Create(*outputFilePath)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer file.Close()
		output = file
		log.Printf("Results will be saved to %s", *outputFilePath)
	}

	// Read the log file
	file, err := os.Open(*logFilePath)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Regular expression to extract creation time - updated to match format "28.081µs"
	creationTimeRegex := regexp.MustCompile(`Creation Time=(\d+\.?\d*)µs`)
	
	var creationTimes []float64
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Block created") {
			matches := creationTimeRegex.FindStringSubmatch(line)
			if len(matches) == 2 {
				timeValue, err := strconv.ParseFloat(matches[1], 64)
				if err != nil {
					log.Printf("Failed to parse time value: %v", err)
					continue
				}
				
				creationTimes = append(creationTimes, timeValue)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading log file: %v", err)
	}

	if len(creationTimes) == 0 {
		log.Fatal("No creation times found in the log file")
	}

	// Calculate statistics
	min, max := minMax(creationTimes)
	mean := calculateMean(creationTimes)
	median := calculateMedian(creationTimes)
	stdDev := calculateStdDev(creationTimes, mean)
	p95 := calculatePercentile(creationTimes, 95)
	p99 := calculatePercentile(creationTimes, 99)

	// Print results
	fmt.Fprintln(output, "Block Creation Time Statistics (in microseconds):")
	fmt.Fprintf(output, "Total blocks analyzed: %d\n", len(creationTimes))
	fmt.Fprintf(output, "Min: %.3f µs\n", min)
	fmt.Fprintf(output, "Max: %.3f µs\n", max)
	fmt.Fprintf(output, "Mean: %.3f µs\n", mean)
	fmt.Fprintf(output, "Median: %.3f µs\n", median)
	fmt.Fprintf(output, "Standard Deviation: %.3f µs\n", stdDev)
	fmt.Fprintf(output, "95th Percentile: %.3f µs\n", p95)
	fmt.Fprintf(output, "99th Percentile: %.3f µs\n", p99)
	
	// Print histogram
	fmt.Fprintln(output, "\nCreation Time Distribution (µs):")
	printHistogram(output, creationTimes, 10)

	// Group by transaction count if available
	analyzeByTransactionCount(output, *logFilePath)
}

func minMax(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	
	min := values[0]
	max := values[0]
	
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	
	return min, max
}

func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	var sum float64
	for _, v := range values {
		sum += v
	}
	
	return sum / float64(len(values))
}

func calculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	// Create a copy to avoid modifying the original slice
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	middle := len(sorted) / 2
	
	if len(sorted)%2 == 0 {
		return (sorted[middle-1] + sorted[middle]) / 2
	}
	
	return sorted[middle]
}

func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}
	
	var sumSquaredDiff float64
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	
	variance := sumSquaredDiff / float64(len(values))
	return math.Sqrt(variance)
}

func calculatePercentile(values []float64, percentile int) float64 {
	if len(values) == 0 {
		return 0
	}
	
	// Create a copy to avoid modifying the original slice
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	index := int(math.Ceil(float64(percentile)/100.0*float64(len(sorted)))) - 1
	// Ensure index is within bounds
	if index < 0 {
		index = 0
	} else if index >= len(sorted) {
		index = len(sorted) - 1
	}
	
	return sorted[index]
}

func printHistogram(w io.Writer, values []float64, bins int) {
	if len(values) == 0 || bins <= 0 {
		return
	}
	
	min, max := minMax(values)
	
	// Add a small buffer to max to ensure the highest value falls within a bin
	max += 0.001
	
	binWidth := (max - min) / float64(bins)
	histogram := make([]int, bins)
	
	// Count values in each bin
	for _, v := range values {
		binIndex := int((v - min) / binWidth)
		// Handle edge case for the max value
		if binIndex >= bins {
			binIndex = bins - 1
		}
		histogram[binIndex]++
	}
	
	// Find the maximum count for scaling
	maxCount := 0
	for _, count := range histogram {
		if count > maxCount {
			maxCount = count
		}
	}
	
	// Print the histogram
	maxBarWidth := 50
	for i := 0; i < bins; i++ {
		lowerBound := min + float64(i)*binWidth
		upperBound := min + float64(i+1)*binWidth
		count := histogram[i]
		
		// Calculate bar width
		var barWidth int
		if maxCount > 0 {
			barWidth = count * maxBarWidth / maxCount
		}
		
		bar := strings.Repeat("█", barWidth)
		fmt.Fprintf(w, "%7.1f - %7.1f µs | %4d | %s\n", lowerBound, upperBound, count, bar)
	}
}

func analyzeByTransactionCount(w io.Writer, logFilePath string) {
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Printf("Failed to reopen log file for transaction analysis: %v", err)
		return
	}
	defer file.Close()

	transactionRegex := regexp.MustCompile(`Transactions=(\d+)`)
	creationTimeRegex := regexp.MustCompile(`Creation Time=(\d+\.?\d*)µs`)
	
	transactionGroups := make(map[int][]float64)
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Block created") {
			// Extract transaction count
			transMatches := transactionRegex.FindStringSubmatch(line)
			timeMatches := creationTimeRegex.FindStringSubmatch(line)
			
			if len(transMatches) == 2 && len(timeMatches) == 2 {
				transCount, err := strconv.Atoi(transMatches[1])
				if err != nil {
					continue
				}
				
				timeValue, err := strconv.ParseFloat(timeMatches[1], 64)
				if err != nil {
					continue
				}
				
				transactionGroups[transCount] = append(transactionGroups[transCount], timeValue)
			}
		}
	}

	if len(transactionGroups) > 0 {
		fmt.Fprintln(w, "\nStatistics Grouped by Transaction Count:")
		
		// Process each transaction group
		for transCount, times := range transactionGroups {
			if len(times) <= 1 {
				continue // Skip transaction counts with only one sample
			}
			
			mean := calculateMean(times)
			stdDev := calculateStdDev(times, mean)
			min, max := minMax(times)
			
			fmt.Fprintf(w, "\nTransaction Count: %d (Blocks: %d)\n", transCount, len(times))
			fmt.Fprintf(w, "  Min Creation Time: %.3f µs\n", min)
			fmt.Fprintf(w, "  Max Creation Time: %.3f µs\n", max)
			fmt.Fprintf(w, "  Mean Creation Time: %.3f µs\n", mean)
			fmt.Fprintf(w, "  Std Deviation: %.3f µs\n", stdDev)
			
			if len(times) >= 20 { // Only show histogram for transaction counts with sufficient samples
				fmt.Fprintf(w, "\n  Creation Time Distribution:\n")
				printHistogram(w, times, 8)
			}
		}
	}
} 