package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"cheaptrick/internal/server"
	"cheaptrick/internal/tui"
)

var (
	port        string
	fixturesDir string
	logFile     string
	tlsCert     string
	tlsKey      string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the TUI and HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		if fixturesDir != "" {
			if err := os.MkdirAll(fixturesDir, 0755); err != nil {
				log.Fatalf("Failed to create fixtures dir: %v", err)
			}
		}

		reqCh := make(chan server.PendingRequest, 100)
		eventCh := make(chan string, 100) // For log messages in notification bar

		go server.StartHTTPServer(port, tlsCert, tlsKey, fixturesDir, logFile, reqCh, eventCh)

		p := tea.NewProgram(
			tui.InitialModel(reqCh, eventCh, fixturesDir),
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	startCmd.Flags().StringVar(&fixturesDir, "fixtures", "", "Directory for fixture files")
	startCmd.Flags().StringVar(&logFile, "log", "mock_log.jsonl", "JSONL log file")
	startCmd.Flags().StringVar(&tlsCert, "tls-cert", "", "Path to TLS cert file")
	startCmd.Flags().StringVar(&tlsKey, "tls-key", "", "Path to TLS key file")

	rootCmd.AddCommand(startCmd)
}
