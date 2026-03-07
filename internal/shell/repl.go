package shell

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/chzyer/readline"
	"google.golang.org/genai"
)

// Config holds the configuration for the interactive shell.
type Config struct {
	BaseURL     string
	Model       string
	APIKey      string
	HistoryPath string
	ToolResDir  string
	AutoMode    bool
	MaxTurns    int
}

// REPL manages the read-eval-print loop and conversation state.
type REPL struct {
	cfg           Config
	client        *genai.Client
	history       []*genai.Content
	rl            *readline.Instance
	callCounts    map[string]int
	currentTrace  *TraceInfo
	injectedFails map[string]bool
	injectedTimes map[string]int
}

// NewREPL creates and initializes a new REPL.
func NewREPL(cfg Config) (*REPL, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{
			BaseURL: cfg.BaseURL,
		},
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	var completer = readline.NewPrefixCompleter(
		readline.PcItem("/clear"),
		readline.PcItem("/auto"),
		readline.PcItem("/step"),
		readline.PcItem("/history"),
		readline.PcItem("/trace"),
		readline.PcItem("/export"),
		readline.PcItem("/fail"),
		readline.PcItem("/timeout"),
		readline.PcItem("/tools"),
		readline.PcItem("/help"),
		readline.PcItem("/quit"),
		readline.PcItem("/exit"),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[36mPrompt>\033[0m ",
		HistoryFile:     cfg.HistoryPath,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize readline: %w", err)
	}

	return &REPL{
		cfg:           cfg,
		client:        client,
		history:       make([]*genai.Content, 0),
		rl:            rl,
		callCounts:    make(map[string]int),
		injectedFails: make(map[string]bool),
		injectedTimes: make(map[string]int),
	}, nil
}

// Close cleans up resources.
func (r *REPL) Close() error {
	if r.rl != nil {
		return r.rl.Close()
	}
	return nil
}

// Run starts the REPL loop.
func (r *REPL) Run(ctx context.Context) error {
	fmt.Printf("Connected to local Gemini Mock Server (%s)\n", r.cfg.BaseURL)
	fmt.Printf("Using model: %s\n", r.cfg.Model)
	if r.cfg.ToolResDir != "" {
		fmt.Printf("Tool responses dir: %s (Auto mode: %v)\n", r.cfg.ToolResDir, r.cfg.AutoMode)
	}
	fmt.Println("Enter your prompts below. Press Ctrl+D (EOF) to exit.")
	fmt.Println("---------------------------------------------------------")

	for {
		line, err := r.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				continue // Handle Ctrl+C gracefully
			}
			if err == io.EOF {
				break // Handle Ctrl+D
			}
			return fmt.Errorf("read error: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle slash commands here (to be implemented)
		if strings.HasPrefix(line, "/") {
			r.handleCommand(line, ctx)
			continue
		}

		r.history = append(r.history, &genai.Content{
			Role:  "user",
			Parts: []*genai.Part{genai.NewPartFromText(line)},
		})

		r.currentTrace = NewTraceInfo()
		r.currentTrace.AddUserTurn(line)

		r.processTurn(ctx)

		r.currentTrace.PrintTrace()
	}
	return nil
}

func (r *REPL) processTurn(ctx context.Context) {
	turnCount := 0

	for {
		if r.cfg.MaxTurns > 0 && turnCount >= r.cfg.MaxTurns {
			fmt.Printf("\033[1;31m⚠ Hit maximum tool-call turns (%d). Aborting loop.\033[0m\n", r.cfg.MaxTurns)
			break
		}
		turnCount++

		resp, err := r.client.Models.GenerateContent(ctx, r.cfg.Model, r.history, nil)
		if err != nil {
			fmt.Printf("\033[31mError generating content: %v. Is the mock server running at %s?\033[0m\n", err, r.cfg.BaseURL)
			return // Cannot proceed with this turn
		}

		// Check response parts
		var parts []*genai.Part
		if len(resp.Candidates) > 0 {
			parts = resp.Candidates[0].Content.Parts
		}

		if len(parts) == 0 {
			fmt.Println("\033[31mEmpty response from mock server.\033[0m")
			return
		}

		// Append the model's response to history
		r.history = append(r.history, &genai.Content{
			Role:  "model",
			Parts: parts,
		})

		fmt.Println("\033[32mMock Server Response:\033[0m")

		var calls []*genai.FunctionCall
		hasFunctionCall := false

		// If there is text, record the final text turn or partial text turn
		for _, part := range parts {
			if part.Text != "" {
				out, err := glamour.Render(part.Text, "dark")
				if err != nil {
					out = part.Text + "\n"
				}
				fmt.Print(out)
				r.currentTrace.AddMockTextTurn(part.Text)
			}
			if part.FunctionCall != nil {
				hasFunctionCall = true
				calls = append(calls, part.FunctionCall)
			}
		}
		fmt.Println()

		if !hasFunctionCall {
			// Done with this turn
			break
		}

		responseParts := r.handleFunctionCalls(ctx, calls)
		r.history = append(r.history, &genai.Content{
			Role:  "user",
			Parts: responseParts,
		})
	}
}
