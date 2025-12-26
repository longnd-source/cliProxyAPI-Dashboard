package pricing

import "strings"

// ModelPrice defines the cost per 1M tokens (USD).
type ModelPrice struct {
	Input  float64
	Output float64
}

// KnownPrices stores the pricing for supported models.
// Prices are in USD per 1 Million tokens.
var KnownPrices = map[string]ModelPrice{
	"gemini-1.5-flash":        {Input: 0.075, Output: 0.30},
	"gemini-1.5-pro":          {Input: 1.25, Output: 5.00},
	"gemini-1.0-pro":          {Input: 0.50, Output: 1.50},
	"gpt-4o":                  {Input: 2.50, Output: 10.00},
	"gpt-4o-mini":             {Input: 0.15, Output: 0.60},
	"claude-3-5-sonnet":       {Input: 3.00, Output: 15.00},
	"claude-3-haiku":          {Input: 0.25, Output: 1.25},

	// New models ($5 Input / $20 Output)
	"gemini-3-pro-preview":                    {Input: 5.00, Output: 20.00},
	"gpt-oss-120b-medium":                     {Input: 5.00, Output: 20.00},
	"gemini-3-pro-image-preview":              {Input: 5.00, Output: 20.00},
	"gemini-claude-sonnet-4-5":                {Input: 5.00, Output: 20.00},
	"gemini-2.5-flash":                        {Input: 5.00, Output: 20.00},
	"gemini-2.5-flash-lite":                   {Input: 5.00, Output: 20.00},
	"gemini-claude-sonnet-4-5-thinking":       {Input: 5.00, Output: 20.00},
	"gemini-claude-opus-4-5-thinking":         {Input: 5.00, Output: 20.00},
	"gemini-2.5-computer-use-preview-10-2025": {Input: 5.00, Output: 20.00},
}

// CalculateCost returns the estimated cost in USD for the given usage.
func CalculateCost(model string, inputTokens, outputTokens int64) float64 {
	// Normalize model name for partial matching if needed, 
	// but for now strict or simple prefix matching
	model = strings.ToLower(model)
	
	// Try exact match first
	price, ok := KnownPrices[model]
	if !ok {
		// Try fuzzy matching for versions/suffixes
		for k, v := range KnownPrices {
			if strings.Contains(model, k) {
				price = v
				break
			}
		}
	}

	// Calculate: (tokens / 1,000,000) * price
	inCost := (float64(inputTokens) / 1000000.0) * price.Input
	outCost := (float64(outputTokens) / 1000000.0) * price.Output

	return inCost + outCost
}
