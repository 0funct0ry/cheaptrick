package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

var (
	shellHost    string
	shellPort    int
	shellAPIKey  string
	shellModel   string
	shellHistory string
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start an interactive shell connecting to the mock server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		baseURL := fmt.Sprintf("http://%s:%d", shellHost, shellPort)

		// Initialize the client configured to hit the mock server
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			HTTPOptions: genai.HTTPOptions{
				BaseURL: baseURL,
			},
			APIKey: shellAPIKey,
		})
		if err != nil {
			log.Fatalf("Failed to create genai client: %v", err)
		}

		rl, err := readline.NewEx(&readline.Config{
			Prompt:          "\033[36mPrompt>\033[0m ",
			HistoryFile:     shellHistory,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
		})
		if err != nil {
			log.Fatalf("Failed to initialize readline: %v", err)
		}
		defer rl.Close()

		fmt.Printf("Connected to local Gemini Mock Server (%s)\n", baseURL)
		fmt.Printf("Using model: %s\n", shellModel)
		fmt.Println("Enter your prompts below. Press Ctrl+D (EOF) to exit.")
		fmt.Println("---------------------------------------------------------")

		for {
			line, err := rl.Readline()
			if err != nil { // EOF (Ctrl+D) or Ctrl+C
				break
			}
			if line == "" {
				continue
			}

			// Generate content
			resp, err := client.Models.GenerateContent(ctx, shellModel, genai.Text(line), nil)
			if err != nil {
				fmt.Printf("\033[31mError generating content: %v\033[0m\n", err)
				continue
			}

			// Print response text
			fmt.Println("\033[32mMock Server Response:\033[0m")
			if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
				fmt.Printf("%s\n\n", resp.Text())
			} else {
				fmt.Printf("%+v\n\n", resp)
			}
		}
	},
}

func init() {
	defaultHistory := filepath.Join(os.TempDir(), "go-genai-client-history.tmp")

	shellCmd.Flags().StringVarP(&shellHost, "host", "H", "localhost", "Host address of the mock server")
	shellCmd.Flags().IntVarP(&shellPort, "port", "p", 8080, "Port of the mock server")
	shellCmd.Flags().StringVar(&shellAPIKey, "api-key", "mock-key", "API key to use (cheaptrick skips validation)")
	shellCmd.Flags().StringVarP(&shellModel, "model", "m", "gemini-2.0-flash", "Gemini model to use in requests")
	shellCmd.Flags().StringVar(&shellHistory, "history-file", defaultHistory, "Path to the readline history file")

	rootCmd.AddCommand(shellCmd)
}
