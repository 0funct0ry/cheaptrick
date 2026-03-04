package main

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

type PendingRequest struct {
	ID         string
	Model      string
	Timestamp  time.Time
	RawBody    []byte
	ParsedBody map[string]interface{}
	Hash       string
	ResponseCh chan string
	ErrorCh    chan error
}

func startHTTPServer(port, tlsCert, tlsKey, fixturesDir, logFile string, reqCh chan<- PendingRequest, eventCh chan<- string) {
	mux := http.NewServeMux()

	handler := func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		var modelName string
		for _, p := range pathParts {
			if strings.HasPrefix(p, "models/") {
				modelName = strings.TrimPrefix(p, "models/")
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

		canonical, _ := json.Marshal(parsed)
		h := sha256.Sum256(canonical)
		hashStr := hex.EncodeToString(h[:])

		reqID := nextReqID()
		timestamp := time.Now()

		if fixturesDir != "" {
			fixture, ok := GetFixture(fixturesDir, hashStr)
			if ok {
				eventCh <- fmt.Sprintf("%s auto-replied from fixture %s", reqID, hashStr[:8])
				logRequestResponse(logFile, reqID, timestamp, string(body), fixture, true)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(fixture))
				return
			}
		}

		respCh := make(chan string, 1)
		errCh := make(chan error, 1)

		req := PendingRequest{
			ID:         reqID,
			Model:      modelName,
			Timestamp:  timestamp,
			RawBody:    body,
			ParsedBody: parsed,
			Hash:       hashStr,
			ResponseCh: respCh,
			ErrorCh:    errCh,
		}

		select {
		case reqCh <- req:
			eventCh <- fmt.Sprintf("%s received for model %s", reqID, modelName)
		default:
			http.Error(w, "TUI queue full", http.StatusServiceUnavailable)
			return
		}

		select {
		case respBody := <-req.ResponseCh:
			logRequestResponse(logFile, reqID, timestamp, string(body), respBody, false)
			eventCh <- fmt.Sprintf("%s response sent", reqID)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(respBody))
		case err := <-req.ErrorCh:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case <-time.After(5 * time.Minute):
			eventCh <- fmt.Sprintf("%s timed out", reqID)
			http.Error(w, "Timeout waiting for TUI response", http.StatusGatewayTimeout)
		case <-r.Context().Done():
			eventCh <- fmt.Sprintf("%s client disconnected", reqID)
		}
	}

	mux.HandleFunc("/v1beta/", handler)

	addr := ":" + port
	eventCh <- "Starting server on " + addr

	var err error
	if tlsCert != "" && tlsKey != "" {
		eventCh <- "HTTPS Enabled"
		err = http.ListenAndServeTLS(addr, tlsCert, tlsKey, mux)
	} else {
		err = http.ListenAndServe(addr, mux)
	}

	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func logRequestResponse(logFile, reqID string, timestamp time.Time, req, resp string, auto bool) {
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
	f.Write(append(b, '\n'))
}
