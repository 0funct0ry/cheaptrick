package fixture

import (
	"bytes"
	"cheaptrick/internal/fixture/data"
	"cheaptrick/internal/fixture/manifest"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"text/template"
)

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
	mf := manifest.NewTextManifest()

	generatedTextPrompts := data.GenerateTextPromptDataset(count)
	for prompt, response := range generatedTextPrompts {
		h := sha256.Sum256([]byte(prompt))
		hashStr := hex.EncodeToString(h[:])

		if err := SaveFixture(outputDir, hashStr, textPromptResponseTemplate(response)); err != nil {
			log.Printf("Failed to save text fixture %s: %v", hashStr, err)
		} else {
			mf.Add(prompt, hashStr)
		}
	}
	if err := mf.SaveMarkdown(outputDir); err != nil {
		log.Fatalf("Failed to save markdown manifest: %v", err)
	}
	if err := mf.SaveJSON(outputDir); err != nil {
		log.Fatalf("Failed to save JSON manifest: %v", err)
	}

}

func generateToolCallPrompts(outputDir string, count int) {
	mf := manifest.NewToolCallManifest()

	generatedToolCallPrompts := data.GenerateToolCallDataset(count)
	for prompt, response := range generatedToolCallPrompts {
		h := sha256.Sum256([]byte(prompt))
		hashStr := hex.EncodeToString(h[:])

		if err := SaveFixture(outputDir, hashStr, response); err != nil {
			log.Printf("Failed to save tool call fixture %s: %v", hashStr, err)
		} else {
			mf.Add(prompt, hashStr)
		}
	}

	if err := mf.SaveMarkdown(outputDir); err != nil {
		log.Fatalf("Failed to save markdown manifest: %v", err)
	}
	if err := mf.SaveJSON(outputDir); err != nil {
		log.Fatalf("Failed to save JSON manifest: %v", err)
	}
}
