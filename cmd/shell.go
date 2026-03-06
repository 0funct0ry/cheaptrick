package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"cheaptrick/internal/shell"
)

var (
	shellHost    string
	shellPort    int
	shellModel   string
	shellHistory string
	shellTools   string
	shellAuto    bool
	shellMax     int
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Start an interactive shell connecting to the mock server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		baseURL := "http://" + shellHost + ":" + cmd.Flag("port").Value.String()

		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			apiKey = "mock-key"
		}

		cfg := shell.Config{
			BaseURL:     baseURL,
			Model:       shellModel,
			APIKey:      apiKey,
			HistoryPath: shellHistory,
			ToolResDir:  shellTools,
			AutoMode:    shellAuto,
			MaxTurns:    shellMax,
		}

		repl, err := shell.NewREPL(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize shell REPL: %v", err)
		}
		defer repl.Close()

		if err := repl.Run(ctx); err != nil {
			log.Fatalf("Shell REPL error: %v", err)
		}
	},
}

func init() {
	defaultHistory := filepath.Join(os.TempDir(), "go-genai-client-history.tmp")

	shellCmd.Flags().StringVarP(&shellHost, "host", "H", "localhost", "Host address of the mock server")
	shellCmd.Flags().IntVarP(&shellPort, "port", "p", 8080, "Port of the mock server")
	shellCmd.Flags().StringVarP(&shellModel, "model", "m", "gemini-2.0-flash", "Gemini model to use in requests")
	shellCmd.Flags().StringVar(&shellHistory, "history-file", defaultHistory, "Path to the readline history file")

	shellCmd.Flags().StringVar(&shellTools, "tool-responses", "", "Directory containing canned tool response files")
	shellCmd.Flags().BoolVar(&shellAuto, "auto", false, "Enable auto-mode for tool responses")
	shellCmd.Flags().IntVar(&shellMax, "max-turns", 20, "Maximum turns in a tool call chain before aborting")

	rootCmd.AddCommand(shellCmd)
}
