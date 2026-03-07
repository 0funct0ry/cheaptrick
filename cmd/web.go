package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/spf13/cobra"

	"cheaptrick/internal/server"
	"cheaptrick/internal/store"
	"cheaptrick/internal/web"
)

var (
	port        string
	fixturesDir string
	logFile     string
	tlsCert     string
	tlsKey      string
	webPort     string
	openUI      bool
)

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Could not open browser: %v\n", err)
	}
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Launch the mock Gemini server and a browser-based UI",
	Run: func(cmd *cobra.Command, args []string) {
		if fixturesDir != "" {
			if err := os.MkdirAll(fixturesDir, 0755); err != nil {
				log.Fatalf("Failed to create fixtures dir: %v", err)
			}
		}

		reqStore := store.New()

		// The mock server handles the main port
		go server.StartHTTPServer(port, tlsCert, tlsKey, fixturesDir, logFile, reqStore)

		// The Web UI handles the web-port
		r := web.NewRouter(reqStore, fixturesDir)
		addr := ":" + webPort
		go func() {
			fmt.Printf("Web UI listening on http://localhost:%s\n", webPort)
			if err := r.Run(addr); err != nil {
				log.Fatalf("Web UI server failed: %v", err)
			}
		}()

		if openUI {
			openBrowser(fmt.Sprintf("http://localhost:%s", webPort))
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("\nShutting down servers...")
	},
}

func init() {
	webCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on for the Gemini mock server")
	webCmd.Flags().StringVar(&webPort, "web-port", "3000", "Port to listen on for the Web UI server")
	webCmd.Flags().StringVar(&fixturesDir, "fixtures", "", "Directory for fixture files")
	webCmd.Flags().StringVar(&logFile, "log", "mock_log.jsonl", "JSONL log file")
	webCmd.Flags().StringVar(&tlsCert, "tls-cert", "", "Path to TLS cert file")
	webCmd.Flags().StringVar(&tlsKey, "tls-key", "", "Path to TLS key file")
	webCmd.Flags().BoolVar(&openUI, "open", true, "Auto-open browser on startup")

	rootCmd.AddCommand(webCmd)
}
