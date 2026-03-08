package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const markdownManifestTemplate = `
# Cheaptrick Fixtures Manifest

This file maps predefined user prompts to their corresponding auto-reply fixture files. 
Send the exact payload structure to get the fixture matched automatically.

## {{ .Type }} Responses

| Prompt		| Fixture File		|
|---------------|---------------------|
{{range .Entries}}| ` + "`{{.Prompt}}`" + ` | [` + "`{{.Hash}}`" + `]({{.Hash}}.json) |
{{end}}
`

type FixtureType string

const (
	Text     FixtureType = "text"
	ToolCall             = "tool-call"
)

type Entry struct {
	Prompt string `json:"prompt"`
	Hash   string `json:"hash"`
}

type Manifest struct {
	FixtureType FixtureType `json:"fixture_type"`
	Entries     []Entry     `json:"entries"`
}

func NewTextManifest() *Manifest {
	return &Manifest{
		FixtureType: Text,
		Entries:     make([]Entry, 0),
	}
}

func NewToolCallManifest() *Manifest {
	return &Manifest{
		FixtureType: ToolCall,
		Entries:     make([]Entry, 0),
	}
}

func (m *Manifest) Add(prompt string, hash string) {
	m.Entries = append(m.Entries, Entry{
		Prompt: prompt,
		Hash:   hash,
	})
}

func (m *Manifest) Len() int {
	return len(m.Entries)
}

func (m *Manifest) Type() string {
	switch m.FixtureType {
	case Text:
		return "Text"
	case ToolCall:
		return "Tool Call"
	default:
		return "Unknown"
	}
}

func (m *Manifest) SaveMarkdown(outputDir string) error {
	manifestPath := filepath.Join(outputDir, "MANIFEST.md")
	f, err := os.Create(manifestPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	tmpl, err := template.New("manifest").Parse(markdownManifestTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	if err := tmpl.Execute(f, m); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}

func (m *Manifest) SaveJSON(outputDir string) error {
	manifestPath := filepath.Join(outputDir, "MANIFEST.json")

	f, err := os.Create(manifestPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m); err != nil {
		return fmt.Errorf("failed to encode manifest to JSON: %v", err)
	}

	return nil
}
