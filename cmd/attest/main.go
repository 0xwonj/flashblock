package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"flashblock/internal/attest"
)

func main() {
	var (
		userData string
	)

	flag.StringVar(&userData, "data", "", "User data to include in the quote (hex encoded)")
	flag.Parse()

	// Decode user data if provided
	var userDataBytes []byte
	var err error
	if userData != "" {
		userDataBytes, err = hex.DecodeString(userData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding user data: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize TDX provider
	provider, err := attest.NewTDXProvider()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing TDX provider: %v\n", err)
		os.Exit(1)
	}

	// Generate quote using the provider
	quoteBytes, err := provider.GetQuote(userDataBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating TDX quote: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("TDX Quote generated successfully (size: %d bytes)\n", len(quoteBytes))
	fmt.Printf("Quote (hex): %x\n", quoteBytes)
}
