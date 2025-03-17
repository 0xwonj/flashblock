package attest

import (
	"fmt"

	"github.com/google/go-tdx-guest/client"
)

// TDXProvider encapsulates the TDX quote provider
type TDXProvider struct {
	provider client.QuoteProvider
}

// NewTDXProvider creates a new TDX provider
func NewTDXProvider() (*TDXProvider, error) {
	// Get the quote provider once at initialization
	quoteProvider, err := client.GetQuoteProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to get quote provider: %v", err)
	}

	return &TDXProvider{
		provider: quoteProvider,
	}, nil
}

// GetQuote generates a TDX quote using the existing provider
func (p *TDXProvider) GetQuote(userData []byte) ([]byte, error) {
	// Prepare the report data (64 bytes)
	var reportData [64]byte
	if userData != nil {
		copy(reportData[:], userData)
	}

	// Get the raw quote using the cached provider
	rawQuote, err := client.GetRawQuote(p.provider, reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to get raw TDX quote: %v", err)
	}

	return rawQuote, nil
}
