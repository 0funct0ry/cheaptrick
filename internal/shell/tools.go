package shell

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
	"google.golang.org/genai"
)

var (
	callBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("208")). // Orange
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(1)

	callTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")). // Yellow
			Bold(true).
			Underline(true)

	callArgValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")) // Cyan

	alertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
)

// handleFunctionCalls processes a batch of function calls and returns the corresponding function responses.
func (r *REPL) handleFunctionCalls(ctx context.Context, calls []*genai.FunctionCall) []*genai.Part {
	var responses []*genai.Part

	for _, call := range calls {
		var content strings.Builder
		content.WriteString(callTitleStyle.Render("FUNCTION CALL: "+call.Name) + "\n\n")
		content.WriteString("Arguments:\n")

		var args map[string]any
		if call.Args != nil {
			args = call.Args
		} else {
			args = make(map[string]any)
		}

		if len(args) == 0 {
			content.WriteString("  {}\n")
		} else {
			b, _ := json.MarshalIndent(args, "", "  ")
			lines := strings.Split(string(b), "\n")
			for _, l := range lines {
				content.WriteString("  " + callArgValueStyle.Render(l) + "\n")
			}
		}

		fmt.Println(callBoxStyle.Render(content.String()))

		// Handle active timeouts
		if waitSecs, ok := r.injectedTimes[call.Name]; ok && waitSecs > 0 {
			fmt.Printf("  %s Delaying %ds...\n", alertStyle.Render("[INJECTED TIMEOUT]"), waitSecs)
			time.Sleep(time.Duration(waitSecs) * time.Second)
			delete(r.injectedTimes, call.Name) // One-shot timeout
		}

		// Handle active failures
		if failActivated, ok := r.injectedFails[call.Name]; ok && failActivated {
			fmt.Printf("  %s Returning generic error.\n", alertStyle.Render("[INJECTED FAILURE]"))

			respMap := map[string]any{"error": fmt.Sprintf("Service unavailable: %s failed (injected by /fail)", call.Name)}

			// Always one-shot in our simple struct
			delete(r.injectedFails, call.Name)

			r.currentTrace.AddMockFuncTurn(call.Name, args)
			r.currentTrace.AddToolTurn(respMap, "[injected: error]")
			responses = append(responses, genai.NewPartFromFunctionResponse(call.Name, respMap))
			continue
		}

		r.callCounts[call.Name]++
		callCount := r.callCounts[call.Name]

		respMap, loadedFilePath := r.resolveCannedResponse(call.Name, args, callCount)

		if loadedFilePath != "" {
			fmt.Printf("  Canned response: \033[2m%s\033[0m\n", loadedFilePath)

			// Pretty print canned response
			b, _ := json.MarshalIndent(respMap, "   ", "  ")
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(string(b)) + "\n")

			if r.cfg.AutoMode {
				fmt.Printf("  [Auto-accepted]\n")
			} else {
				// Step mode interactive prompt
				respMap = r.promptCannedResponseAction(call.Name, respMap, loadedFilePath)
				if respMap == nil {
					// User aborted
					return nil
				}
			}
		} else {
			if r.cfg.ToolResDir != "" {
				fmt.Printf("\n  %s No canned response found for \"%s\"\n", alertStyle.Render("⚠ WARNING:"), call.Name)
				fmt.Printf("  Searched in: %s\n", r.cfg.ToolResDir)
			}
			respMap = r.promptManualToolResponse(call.Name)
			if respMap == nil {
				return nil
			}
			loadedFilePath = "manual input"
		}

		r.currentTrace.AddMockFuncTurn(call.Name, args)
		tag := "[via TUI]"
		if loadedFilePath != "manual input" {
			tag = fmt.Sprintf("[canned: %s]", filepath.Base(loadedFilePath))
		}
		r.currentTrace.AddToolTurn(respMap, tag)

		responses = append(responses, genai.NewPartFromFunctionResponse(call.Name, respMap))
	}

	return responses
}

type matchRule struct {
	Field   string            `json:"field"`
	Matches map[string]string `json:"matches"`
	Default string            `json:"default"`
}

func (r *REPL) resolveCannedResponse(name string, args map[string]any, count int) (map[string]any, string) {
	if r.cfg.ToolResDir == "" {
		return nil, ""
	}

	loadAndSubstitute := func(path string) (map[string]any, string) {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, ""
		}

		// Substitute {{args.field}}
		content := string(b)
		for k, v := range args {
			placeholder := fmt.Sprintf("{{args.%s}}", k)
			content = strings.ReplaceAll(content, placeholder, fmt.Sprintf("%v", v))
		}

		var resp map[string]any
		if err := json.Unmarshal([]byte(content), &resp); err != nil {
			return nil, ""
		}
		return resp, path
	}

	baseDir := r.cfg.ToolResDir

	// 1. Argument-matched file
	matchPath := filepath.Join(baseDir, name, "_match.json")
	if b, err := os.ReadFile(matchPath); err == nil {
		var rule matchRule
		if err := json.Unmarshal(b, &rule); err == nil {
			if argVal, ok := args[rule.Field]; ok {
				argStr := strings.ToLower(fmt.Sprintf("%v", argVal))
				matchedBasename := ""
				for k, v := range rule.Matches {
					if strings.ToLower(k) == argStr {
						matchedBasename = v
						break
					}
				}

				if matchedBasename != "" {
					targetPath := filepath.Join(baseDir, name, matchedBasename+".json")
					if resp, p := loadAndSubstitute(targetPath); resp != nil {
						return resp, p
					}
				}
			}

			// Try default if no match
			if rule.Default != "" {
				targetPath := filepath.Join(baseDir, name, rule.Default+".json")
				if resp, p := loadAndSubstitute(targetPath); resp != nil {
					return resp, p
				}
			}
		}
	}

	// 2. Sequenced
	for i := count; i >= 1; i-- {
		seqPath := filepath.Join(baseDir, fmt.Sprintf("%s.%d.json", name, i))
		if resp, p := loadAndSubstitute(seqPath); resp != nil {
			return resp, p
		}
		// If clamping, we break if we didn't find the highest available file?
		// Actually, standard says "If N exceeds available files, clamp to the highest numbered file."
		// A simple way to do this is to check N, N-1, N-2 ... down to 1 if the exact file doesn't exist.
	}

	// 3. Simple static
	staticPath := filepath.Join(baseDir, fmt.Sprintf("%s.json", name))
	if resp, p := loadAndSubstitute(staticPath); resp != nil {
		return resp, p
	}

	// 4. Subdirectory default
	subDefPath := filepath.Join(baseDir, name, "_default.json")
	if resp, p := loadAndSubstitute(subDefPath); resp != nil {
		return resp, p
	}

	return nil, ""
}

func (r *REPL) promptManualToolResponse(name string) map[string]any {
	fmt.Println()
	fmt.Println("  Type the tool's return value as JSON (brace-balanced):")

	var input strings.Builder
	braceCount := 0
	inString := false
	var escapeNext bool

	for {
		// Use a different prompt for continuation if we are inside a json block
		prompt := "Tool Response> "
		if input.Len() > 0 {
			prompt := "... "
			r.rl.SetPrompt(prompt)
		} else {
			r.rl.SetPrompt(prompt)
		}

		line, err := r.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Abort manual entry? Let's just return empty map for now
				return map[string]any{"error": "aborted"}
			}
			return map[string]any{"error": "read_error"}
		}

		input.WriteString(line)
		input.WriteString("\n")

		// Calculate braces
		for _, char := range line {
			if escapeNext {
				escapeNext = false
				continue
			}
			if char == '\\' {
				escapeNext = true
				continue
			}
			if char == '"' {
				inString = !inString
				continue
			}
			if !inString {
				if char == '{' {
					braceCount++
				} else if char == '}' {
					braceCount--
				}
			}
		}

		if braceCount <= 0 {
			break
		}
	}

	// Reset prompt
	r.rl.SetPrompt("\033[36mPrompt>\033[0m ")

	raw := strings.TrimSpace(input.String())
	if raw == "" {
		return map[string]any{}
	}

	var respMap map[string]any
	if err := json.Unmarshal([]byte(raw), &respMap); err != nil {
		fmt.Printf("\033[31mInvalid JSON: %v. Sending empty response.\033[0m\n", err)
		return map[string]any{}
	}

	return respMap
}

func (r *REPL) promptCannedResponseAction(name string, defaultMap map[string]any, filePath string) map[string]any {
	for {
		fmt.Printf("  [Enter] Accept  [e] Edit  [s] Accept & Save\n")
		fmt.Printf("  [t] Type new    [x] Abort chain\n")

		r.rl.SetPrompt("\033[36mAction>\033[0m ")
		line, err := r.rl.Readline()
		if err != nil {
			return nil
		}
		line = strings.TrimSpace(line)

		switch line {
		case "":
			r.rl.SetPrompt("\033[36mPrompt>\033[0m ")
			return defaultMap
		case "x":
			fmt.Println("Aborting chain.")
			r.rl.SetPrompt("\033[36mPrompt>\033[0m ")
			return nil
		case "t":
			r.rl.SetPrompt("\033[36mPrompt>\033[0m ")
			return r.promptManualToolResponse(name)
		case "e", "s":
			// We can pre-fill readline with the JSON so the user can edit it.
			// Readline doesn't cleanly support multi-line initial strings via standard API,
			// but we can just use promptManualToolResponse for this since user wants to edit/type.
			// Actually, Edit means we let them type a new one, but ideally pre-filled.
			// Let's just prompt them to type a new JSON for now.

			// A simpler way: we just read a single line or multi-line JSON replacing the old
			fmt.Printf("\n  Type new JSON (multi-line supported):\n")
			r.rl.SetPrompt("Tool Response> ")
			newMap := r.promptManualToolResponse(name)
			if newMap == nil {
				// Aborted editing
				continue
			}

			if line == "s" {
				newB, _ := json.MarshalIndent(newMap, "", "  ")
				if err := os.WriteFile(filePath, newB, 0644); err == nil {
					fmt.Printf("  Saved to %s\n", filePath)
				} else {
					fmt.Printf("\033[31mFailed to save: %v\033[0m\n", err)
				}
			}

			r.rl.SetPrompt("\033[36mPrompt>\033[0m ")
			return newMap
		default:
			fmt.Printf("\033[31mInvalid option.\033[0m\n")
		}
	}
}
