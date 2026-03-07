package fixture

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// ComputeRequestHash generated a SHA256 hex string from the canonical JSON payload.
func ComputeRequestHash(parsed map[string]interface{}) string {
	canonical := []byte("")

	if contents, ok := parsed["contents"].([]interface{}); ok && len(contents) > 0 {
		idx := len(contents) - 1
		if firstContent, ok := contents[idx].(map[string]interface{}); ok {
			if parts, ok := firstContent["parts"].([]interface{}); ok && len(parts) > 0 {
				if firstPart, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := firstPart["text"].(string); ok {
						canonical = []byte(text)
					}
				}
			}
		}
	}

	if len(canonical) == 0 {
		canonical, _ = json.Marshal(parsed)
	}

	h := sha256.Sum256(canonical)
	return hex.EncodeToString(h[:])
}
