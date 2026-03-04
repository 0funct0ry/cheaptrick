package main

import (
	"crypto/rand"
	"encoding/hex"
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
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the TUI and HTTP server",
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

	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	startCmd.Flags().StringVar(&fixturesDir, "fixtures", "", "Directory for fixture files")
	startCmd.Flags().StringVar(&logFile, "log", "mock_log.jsonl", "JSONL log file")
	startCmd.Flags().StringVar(&tlsCert, "tls-cert", "", "Path to TLS cert file")
	startCmd.Flags().StringVar(&tlsKey, "tls-key", "", "Path to TLS key file")
	rootCmd.AddCommand(startCmd)

	var outputDir string
	fixturesCmd := &cobra.Command{
		Use:   "fixtures",
		Short: "Generate 20 text fixtures and 10 tool call fixtures",
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < 20; i++ {
				b := make([]byte, 16)
				if _, err := rand.Read(b); err != nil {
					log.Fatalf("Failed to generate random bytes: %v", err)
				}
				hash := hex.EncodeToString(b)
				if err := SaveFixture(outputDir, hash, getTemplateText()); err != nil {
					log.Printf("Failed to save text fixture %s: %v", hash, err)
				} else {
					fmt.Printf("Generated text fixture: %s.json\n", hash)
				}
			}
			for i := 0; i < 10; i++ {
				b := make([]byte, 16)
				if _, err := rand.Read(b); err != nil {
					log.Fatalf("Failed to generate random bytes: %v", err)
				}
				hash := hex.EncodeToString(b)
				if err := SaveFixture(outputDir, hash, getTemplateFunctionCall(PendingRequest{})); err != nil {
					log.Printf("Failed to save tool call fixture %s: %v", hash, err)
				} else {
					fmt.Printf("Generated tool call fixture: %s.json\n", hash)
				}
			}
		},
	}
	fixturesCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "fixtures", "Directory to output fixtures")
	rootCmd.AddCommand(fixturesCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
