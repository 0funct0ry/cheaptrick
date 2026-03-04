package fixture

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var PredefinedTextPrompts = []string{
	"Hello, how are you?",
	"What is the meaning of life?",
	"Tell me a joke.",
	"Write a poem about the ocean.",
	"Explain quantum physics in simple terms.",
	"How do I bake a chocolate cake?",
	"What is the capital of France?",
	"Translate 'good morning' to Spanish.",
	"Summarize the plot of Romeo and Juliet.",
	"Who won the World Series in 2020?",
	"Give me a workout routine for beginners.",
	"What are the benefits of meditation?",
	"Suggest some good sci-fi books to read.",
	"How does a combustion engine work?",
	"Write a short story about a time traveler.",
	"What is the difference between an API and an SDK?",
	"How do I start a vegetable garden?",
	"Explain the greenhouse effect.",
	"What are the symptoms of a common cold?",
	"Give me a recipe for vegan lasagna.",
}

var PredefinedToolCallPrompts = []string{
	"Get the current weather in New York.",
	"Calculate 25 * 48.",
	"Set an alarm for 7 AM tomorrow.",
	"Find cheap flights to Tokyo for next week.",
	"What is the stock price of Apple?",
	"Translate this webpage to French.",
	"Book a table for two at a nearby Italian restaurant.",
	"Play some relaxing jazz music.",
	"Turn off the living room lights.",
	"Convert 100 US dollars to Euros.",
}

func GeneratePromptPayload(promptText string) []byte {
	// The minimal payload structure often seen in testing curls:
	// {"contents":[{"parts":[{"text":"prompt text"}]}]}
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": promptText,
					},
				},
			},
		},
	}
	b, _ := json.Marshal(payload)
	return b
}

func GenerateFromPrompts(outputDir string) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	manifestPath := filepath.Join(outputDir, "MANIFEST.md")
	f, err := os.Create(manifestPath)
	if err != nil {
		log.Fatalf("Failed to create manifest file: %v", err)
	}
	defer f.Close()

	_, _ = f.WriteString("# cheaptrick Fixtures Manifest\n\n")
	_, _ = f.WriteString("This file maps predefined user prompts to their corresponding auto-reply fixture files. Send the exact payload structure to get the fixture matched automatically.\n\n")

	_, _ = f.WriteString("## Text Responses\n\n")
	_, _ = f.WriteString("| Prompt | Fixture File |\n")
	_, _ = f.WriteString("|---|---|\n")

	for _, prompt := range PredefinedTextPrompts {
		payload := GeneratePromptPayload(prompt)
		h := sha256.Sum256(payload)
		hashStr := hex.EncodeToString(h[:])

		if err := SaveFixture(outputDir, hashStr, TemplateText()); err != nil {
			log.Printf("Failed to save text fixture %s: %v", hashStr, err)
		} else {
			_, _ = f.WriteString(fmt.Sprintf("| `%s` | `%s.json` |\n", prompt, hashStr))
		}
	}

	_, _ = f.WriteString("\n## Tool Call Responses\n\n")
	_, _ = f.WriteString("| Prompt | Fixture File |\n")
	_, _ = f.WriteString("|---|---|\n")

	for _, prompt := range PredefinedToolCallPrompts {
		payload := GeneratePromptPayload(prompt)
		h := sha256.Sum256(payload)
		hashStr := hex.EncodeToString(h[:])

		if err := SaveFixture(outputDir, hashStr, TemplateFunctionCall(map[string]interface{}{})); err != nil {
			log.Printf("Failed to save tool call fixture %s: %v", hashStr, err)
		} else {
			_, _ = f.WriteString(fmt.Sprintf("| `%s` | `%s.json` |\n", prompt, hashStr))
		}
	}

	fmt.Printf("Successfully generated 30 fixtures and MANIFEST.md in %s\n", outputDir)
}
