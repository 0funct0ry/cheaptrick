package shell

import (
	"fmt"
	"time"
)

type TraceTurn struct {
	Index   int
	Role    string // user, mock, tool
	Content string
	Tag     string // e.g. [via TUI], [canned: paris.json]
}

type TraceInfo struct {
	StartTime time.Time
	Turns     []*TraceTurn
}

func NewTraceInfo() *TraceInfo {
	return &TraceInfo{
		StartTime: time.Now(),
		Turns:     make([]*TraceTurn, 0),
	}
}

func (t *TraceInfo) AddUserTurn(text string) {
	t.Turns = append(t.Turns, &TraceTurn{
		Index:   len(t.Turns) + 1,
		Role:    "user",
		Content: truncateString(text, 60),
	})
}

func (t *TraceInfo) AddMockTextTurn(text string) {
	t.Turns = append(t.Turns, &TraceTurn{
		Index:   len(t.Turns) + 1,
		Role:    "mock",
		Content: truncateString(text, 60),
		Tag:     "[via TUI]",
	})
}

func (t *TraceInfo) AddMockFuncTurn(callName string, args map[string]any) {
	// Simple summary of args
	argsStr := ""
	for k, v := range args {
		if argsStr != "" {
			argsStr += ", "
		}
		argsStr += fmt.Sprintf("%s=\"%v\"", k, v)
	}
	t.Turns = append(t.Turns, &TraceTurn{
		Index:   len(t.Turns) + 1,
		Role:    "mock",
		Content: fmt.Sprintf("ƒ %s(%s)", callName, truncateString(argsStr, 30)),
		Tag:     "[via TUI]",
	})
}

func (t *TraceInfo) AddToolTurn(respMap map[string]any, fileTag string) {
	// Summarize resp
	respStr := fmt.Sprintf("%v", respMap)
	t.Turns = append(t.Turns, &TraceTurn{
		Index:   len(t.Turns) + 1,
		Role:    "tool",
		Content: fmt.Sprintf("→ %s", truncateString(respStr, 40)),
		Tag:     fileTag,
	})
}

func (t *TraceInfo) PrintTrace() {
	if len(t.Turns) == 0 {
		return
	}

	duration := time.Since(t.StartTime)
	fmt.Printf("\n\033[2m\u2500\u2500\u2500 Chain Complete (%d turns, %.1fs) \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\033[0m\n", len(t.Turns), duration.Seconds())
	for _, turn := range t.Turns {
		roleColor := "\033[0m"
		if turn.Role == "user" {
			roleColor = "\033[34m" // blue
		} else if turn.Role == "mock" {
			roleColor = "\033[33m" // yellow
		} else if turn.Role == "tool" {
			roleColor = "\033[32m" // green
		}

		fmt.Printf(" \033[1;37m%2d\033[0m  %s%-6s\033[0m \u2502 %-60s \033[2m%s\033[0m\n", turn.Index, roleColor, turn.Role, turn.Content, turn.Tag)
	}
	fmt.Printf("\033[2m\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\033[0m\n\n")
}

func truncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length-3] + "..."
}
