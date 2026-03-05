package web

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"cheaptrick/internal/fixture"
	"cheaptrick/internal/store"
)

type apiHandler struct {
	reqStore   *store.Store
	fixtureDir string
}

func extractPreview(parsed map[string]interface{}) string {
	if contents, ok := parsed["contents"].([]interface{}); ok && len(contents) > 0 {
		if contentMap, ok := contents[0].(map[string]interface{}); ok {
			if parts, ok := contentMap["parts"].([]interface{}); ok && len(parts) > 0 {
				if partMap, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := partMap["text"].(string); ok {
						desc := strings.ReplaceAll(text, "\n", " ")
						if len(desc) > 120 {
							return desc[:117] + "..."
						}
						return desc
					}
				}
			}
		}
	}
	return ""
}

func extractTools(parsed map[string]interface{}) []string {
	var names []string
	if tools, ok := parsed["tools"].([]interface{}); ok {
		for _, tool := range tools {
			if tmap, ok := tool.(map[string]interface{}); ok {
				if fds, ok := tmap["functionDeclarations"].([]interface{}); ok {
					for _, fd := range fds {
						if fm, ok := fd.(map[string]interface{}); ok {
							if name, ok := fm["name"].(string); ok {
								names = append(names, name)
							}
						}
					}
				}
			}
		}
	}
	return names
}

func (h *apiHandler) listRequests(c *gin.Context) {
	reqs := h.reqStore.GetRequests()
	sort.Slice(reqs, func(i, j int) bool {
		return reqs[i].Timestamp.After(reqs[j].Timestamp)
	})

	type RequestView struct {
		*store.Request
		Preview string   `json:"preview"`
		Tools   []string `json:"tools"`
	}

	views := make([]RequestView, len(reqs))
	for i, r := range reqs {
		views[i] = RequestView{
			Request: r,
			Preview: extractPreview(r.ParsedBody),
			Tools:   extractTools(r.ParsedBody),
		}
	}

	c.JSON(http.StatusOK, gin.H{"requests": views})
}

func (h *apiHandler) getRequest(c *gin.Context) {
	id := c.Param("id")
	req, ok := h.reqStore.GetRequest(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	sysInst := ""
	if si, ok := req.ParsedBody["systemInstruction"].(map[string]interface{}); ok {
		if parts, ok := si["parts"].([]interface{}); ok && len(parts) > 0 {
			if pm, ok := parts[0].(map[string]interface{}); ok {
				sysInst, _ = pm["text"].(string)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":                 req.ID,
		"timestamp":          req.Timestamp,
		"status":             req.Status,
		"model":              req.Model,
		"body":               req.ParsedBody,
		"tools":              req.ParsedBody["tools"],
		"system_instruction": sysInst,
		"contents":           req.ParsedBody["contents"],
		"generation_config":  req.ParsedBody["generationConfig"],
		"fixture_hash":       req.FixtureHash,
	})
}

func (h *apiHandler) respondToRequest(c *gin.Context) {
	id := c.Param("id")
	req, ok := h.reqStore.GetRequest(id)
	if !ok || req.Status != "pending" {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found or already responded"})
		return
	}

	var reqBody struct {
		Response interface{} `json:"response"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in response body"})
		return
	}

	respBytes, err := json.Marshal(reqBody.Response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	h.reqStore.MarkResponded(id, "manual")
	req.ResponseCh <- string(respBytes)

	c.JSON(http.StatusOK, gin.H{"ok": true, "id": id})
}

func (h *apiHandler) saveFixture(c *gin.Context) {
	id := c.Param("id")
	req, ok := h.reqStore.GetRequest(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	var reqBody struct {
		Response interface{} `json:"response"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON in response body"})
		return
	}

	respBytes, err := json.MarshalIndent(reqBody.Response, "", "  ")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := fixture.SaveFixture(h.fixtureDir, req.Hash, string(respBytes)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save fixture"})
		return
	}

	h.reqStore.NotifyFixtureSaved(req.Hash, id)

	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"hash": req.Hash,
		"path": filepath.Join(h.fixtureDir, req.Hash+".json"),
	})
}

func (h *apiHandler) listTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"templates": []map[string]interface{}{
			{
				"id":       "text",
				"label":    "Text Response",
				"shortcut": "F1",
				"body":     json.RawMessage(fixture.TemplateText()),
			},
			{
				"id":       "function_call",
				"label":    "Function Call",
				"shortcut": "F2",
				"body":     json.RawMessage(fixture.TemplateFunctionCall(map[string]interface{}{})),
			},
			{
				"id":       "error_429",
				"label":    "Rate Limit Error (429)",
				"shortcut": "F3",
				"body":     json.RawMessage(fixture.Template429()),
			},
			{
				"id":       "error_500",
				"label":    "Internal Server Error (500)",
				"shortcut": "F4",
				"body":     json.RawMessage(fixture.Template500()),
			},
		},
	})
}

func (h *apiHandler) listFixtures(c *gin.Context) {
	if h.fixtureDir == "" {
		c.JSON(http.StatusOK, gin.H{"fixtures": []interface{}{}})
		return
	}
	files, err := os.ReadDir(h.fixtureDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"fixtures": []interface{}{}})
		return
	}

	type fix struct {
		Hash string `json:"hash"`
		Size int64  `json:"size"`
	}
	var fixtures []fix
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			info, err := f.Info()
			if err == nil {
				fixtures = append(fixtures, fix{
					Hash: strings.TrimSuffix(f.Name(), ".json"),
					Size: info.Size(),
				})
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"fixtures": fixtures})
}

func (h *apiHandler) deleteFixture(c *gin.Context) {
	hash := c.Param("hash")
	if h.fixtureDir == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "fixtures dir not set"})
		return
	}
	path := filepath.Join(h.fixtureDir, hash+".json")
	if err := os.Remove(path); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "fixture not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *apiHandler) deleteRequest(c *gin.Context) {
	id := c.Param("id")
	if h.reqStore.RemoveRequest(id) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
	}
}

func (h *apiHandler) clearRequests(c *gin.Context) {
	h.reqStore.ClearRespondedRequests()
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
