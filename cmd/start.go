package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"cheaptrick/internal/server"
	"cheaptrick/internal/store"
	"cheaptrick/internal/tui"
)

var (
	port        string
	fixturesDir string
	logFile     string
	tlsCert     string
	tlsKey      string
)

type tuiObserver struct {
	reqCh       chan *store.Request
	eventCh     chan string
	respondedCh chan string
}

func (o *tuiObserver) OnNewRequest(req *store.Request) {
	o.reqCh <- req
}
func (o *tuiObserver) OnRequestResponded(id string, via string) {
	o.respondedCh <- id
}
func (o *tuiObserver) OnFixtureSaved(hash string, reqID string) {
}
func (o *tuiObserver) OnEvent(msg string) {
	o.eventCh <- msg
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the TUI and HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		if fixturesDir != "" {
			if err := os.MkdirAll(fixturesDir, 0755); err != nil {
				log.Fatalf("Failed to create fixtures dir: %v", err)
			}
		}

		reqCh := make(chan *store.Request, 100)
		eventCh := make(chan string, 100)
		respondedCh := make(chan string, 100)

		reqStore := store.New()
		reqStore.Register(&tuiObserver{
			reqCh:       reqCh,
			eventCh:     eventCh,
			respondedCh: respondedCh,
		})

		go server.StartHTTPServer(port, tlsCert, tlsKey, fixturesDir, logFile, reqStore)

		p := tea.NewProgram(
			tui.InitialModel(reqStore, reqCh, eventCh, respondedCh, fixturesDir),
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
