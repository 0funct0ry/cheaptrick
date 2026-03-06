package shell

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (r *REPL) exportConversationSequence(dir string, name string) error {
	if len(r.history) == 0 {
		return fmt.Errorf("conversation history is empty")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	type ManifestTurn struct {
		Turn     int    `json:"turn"`
		Request  string `json:"request"`
		Response string `json:"response"`
	}

	type Manifest struct {
		Name      string         `json:"name"`
		Model     string         `json:"model"`
		CreatedAt string         `json:"created_at"`
		Turns     int            `json:"turns"`
		Files     []ManifestTurn `json:"files"`
	}

	manifest := Manifest{
		Name:      name,
		Model:     r.cfg.Model,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// We need to pair requests and responses.
	// REPL history is: user -> model -> user (FunctionResponse) -> model
	// This maps exactly to turns.
	// We'll iterate the history and group them.

	var turns []ManifestTurn
	turnNum := 1

	for i := 0; i < len(r.history); i++ {
		reqContent := r.history[i]
		if reqContent.Role != "user" {
			continue // Should start with user
		}

		// Reconstruct the request payload as the mock server receives it
		reqPayload := map[string]any{
			"contents": r.history[:i+1], // The history up to this point
		}
		reqBytes, _ := json.MarshalIndent(reqPayload, "", "  ")

		reqFile := fmt.Sprintf("%03d_request.json", turnNum)
		if err := os.WriteFile(filepath.Join(dir, reqFile), reqBytes, 0644); err != nil {
			return err
		}

		var respBytes []byte
		var respFile string

		if i+1 < len(r.history) && r.history[i+1].Role == "model" {
			// There is a model response
			respContent := r.history[i+1]
			// The mock server responds with the content parts wrapped usually
			respPayload := map[string]any{
				"candidates": []map[string]any{
					{
						"content": respContent,
					},
				},
			}
			respBytes, _ = json.MarshalIndent(respPayload, "", "  ")
			respFile = fmt.Sprintf("%03d_response.json", turnNum)
			if err := os.WriteFile(filepath.Join(dir, respFile), respBytes, 0644); err != nil {
				return err
			}
			i++ // skip the model response in the loop
		}

		turns = append(turns, ManifestTurn{
			Turn:     turnNum,
			Request:  reqFile,
			Response: respFile,
		})
		turnNum++
	}

	manifest.Turns = len(turns)
	manifest.Files = turns

	manifestBytes, _ := json.MarshalIndent(manifest, "", "  ")
	return os.WriteFile(filepath.Join(dir, "manifest.json"), manifestBytes, 0644)
}
