package data

import (
	"fmt"
	"testing"
)

func TestGenerateToolCallDataset(t *testing.T) {

	dataset := GenerateToolCallDataset(10)

	for prompt, resp := range dataset {
		fmt.Println(prompt)
		fmt.Println(resp)
	}
	if len(dataset) != 10 {
		t.Errorf("Expected 10 prompts, got %d", len(dataset))
	}
}
