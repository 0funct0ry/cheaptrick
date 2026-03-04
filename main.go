package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func main() {
	var port string
	var fixturesDir string
	var logFile string
	var tlsCert string
	var tlsKey string

	rootCmd := &cobra.Command{
		Use:   "cheaptrick",
		Short: "A mock server for the Gemini API with a built-in TUI",
		Run: func(cmd *cobra.Command, args []string) {
			if fixturesDir != "" {
				if err := os.MkdirAll(fixturesDir, 0755); err != nil {
					log.Fatalf("Failed to create fixtures dir: %v", err)
				}
			}

			reqCh := make(chan PendingRequest, 100)
			eventCh := make(chan string, 100) // For log messages in notification bar

			go startHTTPServer(port, tlsCert, tlsKey, fixturesDir, logFile, reqCh, eventCh)

			p := tea.NewProgram(
				initialModel(reqCh, eventCh, fixturesDir),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)

			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	rootCmd.Flags().StringVar(&fixturesDir, "fixtures", "", "Directory for fixture files")
	rootCmd.Flags().StringVar(&logFile, "log", "mock_log.jsonl", "JSONL log file")
	rootCmd.Flags().StringVar(&tlsCert, "tls-cert", "", "Path to TLS cert file")
	rootCmd.Flags().StringVar(&tlsKey, "tls-key", "", "Path to TLS key file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
