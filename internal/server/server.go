package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"cheaptrick/internal/fixture"
	"cheaptrick/internal/store"
)

var (
	reqIDCounter int
	reqIDMu      sync.Mutex
)

func nextReqID() string {
	reqIDMu.Lock()
	defer reqIDMu.Unlock()
	reqIDCounter++
	return fmt.Sprintf("req-%02d", reqIDCounter)
}

func StartHTTPServer(port, tlsCert, tlsKey, fixturesDir, logFile string, reqStore *store.Store) {
	mux := http.NewServeMux()

	handler := func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		var modelName string
		for i, p := range pathParts {
			if p == "models" && i+1 < len(pathParts) {
				modelName = pathParts[i+1]
				if idx := strings.Index(modelName, ":"); idx != -1 {
					modelName = modelName[:idx]
				}
				break
			}
		}
		if modelName == "" {
			modelName = "unknown"
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(body, &parsed); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		canonical := []byte("")
		if contents, ok := parsed["contents"].([]interface{}); ok && len(contents) > 0 {
			if firstContent, ok := contents[0].(map[string]interface{}); ok {
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
		hashStr := hex.EncodeToString(h[:])

		reqID := nextReqID()
		timestamp := time.Now()

		req := &store.Request{
			ID:         reqID,
			Model:      modelName,
			Timestamp:  timestamp,
			RawBody:    body,
			ParsedBody: parsed,
			Hash:       hashStr,
			ResponseCh: make(chan string, 1),
			ErrorCh:    make(chan error, 1),
		}

		if fixturesDir != "" {
			fixtureContent, ok := fixture.GetFixture(fixturesDir, hashStr)
			if ok {
				reqStore.NotifyEvent(fmt.Sprintf("%s auto-replied from fixture %s", reqID, hashStr[:8]))
				reqStore.AddRequest(req)
				reqStore.MarkResponded(reqID, "fixture")
				LogRequestResponse(logFile, reqID, timestamp, string(body), fixtureContent, true)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(fixtureContent))
				return
			}
		}

		reqStore.NotifyEvent(fmt.Sprintf("%s received for model %s", reqID, modelName))
		reqStore.AddRequest(req)

		select {
		case respBody := <-req.ResponseCh:
			LogRequestResponse(logFile, reqID, timestamp, string(body), respBody, false)
			reqStore.NotifyEvent(fmt.Sprintf("%s response sent", reqID))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(respBody))
		case err := <-req.ErrorCh:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case <-time.After(5 * time.Minute):
			reqStore.NotifyEvent(fmt.Sprintf("%s timed out", reqID))
			http.Error(w, "Timeout waiting for TUI response", http.StatusGatewayTimeout)
		case <-r.Context().Done():
			reqStore.NotifyEvent(fmt.Sprintf("%s client disconnected", reqID))
		}
	}

	mux.HandleFunc("/v1beta/", handler)

	addr := ":" + port
	reqStore.NotifyEvent("Starting server on " + addr)

	var err error
	if tlsCert != "" && tlsKey != "" {
		reqStore.NotifyEvent("HTTPS Enabled")
		err = http.ListenAndServeTLS(addr, tlsCert, tlsKey, mux)
	} else {
		err = http.ListenAndServe(addr, mux)
	}

	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func LogRequestResponse(logFile, reqID string, timestamp time.Time, req, resp string, auto bool) {
	if logFile == "" {
		return
	}
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	entry := map[string]interface{}{
		"id":           reqID,
		"timestamp":    timestamp.Format(time.RFC3339),
		"request":      json.RawMessage(req),
		"auto_fixture": auto,
	}

	var parsedResp interface{}
	if err := json.Unmarshal([]byte(resp), &parsedResp); err == nil {
		entry["response"] = parsedResp
	} else {
		entry["response"] = resp
	}

	b, _ := json.Marshal(entry)
	_, _ = f.Write(append(b, '\n'))
}
