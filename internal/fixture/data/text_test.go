package data

import (
	"fmt"
	"testing"
)

func TestGenerateDataset(t *testing.T) {
	dataset := GenerateTextPromptDataset(20)

	for p, r := range dataset {
		fmt.Printf("%q: %q,\n\n", p, r)
	}

	if len(dataset) != 20 {
		t.Errorf("Expected 20 prompts, got %d", len(dataset))
	}
}
