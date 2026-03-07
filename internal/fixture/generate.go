package fixture

import (
	"bytes"
	"cheaptrick/internal/fixture/data"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

//var PredefinedToolCallPrompts = []string{
//	"Get the current weather in New York.",
//	"Calculate 25 * 48.",
//	"Set an alarm for 7 AM tomorrow.",
//	"Find cheap flights to Tokyo for next week.",
//	"What is the stock price of Apple?",
//	"Translate this webpage to French.",
//	"Book a table for two at a nearby Italian restaurant.",
//	"Play some relaxing jazz music.",
//	"Turn off the living room lights.",
//	"Convert 100 US dollars to Euros.",
//}

func GenerateFromPrompts(outputDir string, promptType string, count int) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	if promptType == "text" {
		generateTextPrompts(outputDir, count)
	} else {
		generateToolCallPrompts(outputDir, count)
	}

	fmt.Printf("Successfully generated %d %s fixtures and MANIFEST.md in %s\n", count, promptType, outputDir)
}

func textPromptResponseTemplate(response string) string {
	const textPromptResponseTemplate = `
{
  "candidates": [
    {
      "content": {
        "role": "model",
        "parts": [
          {
            "text": "{{.Response}}"
          }
        ]
      },
      "finishReason": "STOP"
    }
  ],
  "usageMetadata": {
    "promptTokenCount": 0,
    "candidatesTokenCount": 0,
    "totalTokenCount": 0
  }
}`
	t := template.Must(template.New("t").Parse(textPromptResponseTemplate))
	buffer := bytes.NewBufferString("")
	err := t.Execute(buffer, struct {
		Response string
	}{
		Response: response,
	})
	if err != nil {
		return ""
	}
	return buffer.String()
}

func generateTextPrompts(outputDir string, count int) {
	manifestPath := filepath.Join(outputDir, "MANIFEST.md")
	f, err := os.Create(manifestPath)
	if err != nil {
		log.Fatalf("Failed to create manifest file: %v", err)
	}
	defer f.Close()

	_, _ = f.WriteString("# Cheaptrick Fixtures Manifest\n\n")
	_, _ = f.WriteString("This file maps predefined user prompts to their corresponding auto-reply fixture files. Send the exact payload structure to get the fixture matched automatically.\n\n")

	_, _ = f.WriteString("## Text Responses\n\n")
	_, _ = f.WriteString("| Prompt | Fixture File |\n")
	_, _ = f.WriteString("|---|---|\n")

	generatedTextPrompts := data.GenerateTextPromptDataset(count)
	for prompt, response := range generatedTextPrompts {
		h := sha256.Sum256([]byte(prompt))
		hashStr := hex.EncodeToString(h[:])

		if err := SaveFixture(outputDir, hashStr, textPromptResponseTemplate(response)); err != nil {
			log.Printf("Failed to save text fixture %s: %v", hashStr, err)
		} else {
			_, _ = f.WriteString(fmt.Sprintf("| `%s` | [`%s`](%s.json) |\n", prompt, hashStr, hashStr))
		}
	}

}

func generateToolCallPrompts(outputDir string, count int) {
	manifestPath := filepath.Join(outputDir, "MANIFEST.md")
	f, err := os.Create(manifestPath)
	if err != nil {
		log.Fatalf("Failed to create manifest file: %v", err)
	}
	defer f.Close()

	_, _ = f.WriteString("# Cheaptrick Fixtures Manifest\n\n")
	_, _ = f.WriteString("This file maps predefined user prompts to their corresponding auto-reply fixture files. Send the exact payload structure to get the fixture matched automatically.\n\n")

	_, _ = f.WriteString("\n## Tool Call Responses\n\n")
	_, _ = f.WriteString("| Prompt | Fixture File |\n")
	_, _ = f.WriteString("|---|---|\n")

	generatedToolCallPrompts := data.GenerateToolCallDataset(count)
	for prompt, response := range generatedToolCallPrompts {
		h := sha256.Sum256([]byte(prompt))
		hashStr := hex.EncodeToString(h[:])

		if err := SaveFixture(outputDir, hashStr, response); err != nil {
			log.Printf("Failed to save tool call fixture %s: %v", hashStr, err)
		} else {
			_, _ = f.WriteString(fmt.Sprintf("| `%s` | `%s.json` |\n", prompt, hashStr))
		}
	}
}
