package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var toolsOutputDir string

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Generate 20 sample canned tool response files for use with `cheaptrick shell --tool-responses`",
	Long: `Generates a directory of canned tool response files that simulate common
external tool integrations. These files are used by the shell's tool-call
loop to automatically respond to FunctionCall parts without manual input.

The generated tools cover weather, search, email, databases, calendars,
file I/O, translation, ticketing, messaging, and more — providing a
realistic starting point for agent development.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateTools(toolsOutputDir)
	},
}

func init() {
	toolsCmd.Flags().StringVarP(&toolsOutputDir, "output-dir", "o", "./mock_tools", "Output directory for generated tool files")
	rootCmd.AddCommand(toolsCmd)
}

// tool represents a single canned tool response definition.
type tool struct {
	// dir is the subdirectory name under the output root (empty = root level file).
	dir string
	// file is the filename relative to dir (or root if dir is empty).
	file string
	// data is the JSON-serializable content to write.
	data any
}

func generateTools(outDir string) error {
	tools := allTools()

	for _, t := range tools {
		var path string
		if t.dir != "" {
			dir := filepath.Join(outDir, t.dir)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("creating directory %s: %w", dir, err)
			}
			path = filepath.Join(dir, t.file)
		} else {
			if err := os.MkdirAll(outDir, 0o755); err != nil {
				return fmt.Errorf("creating directory %s: %w", outDir, err)
			}
			path = filepath.Join(outDir, t.file)
		}

		b, err := json.MarshalIndent(t.data, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling %s: %w", path, err)
		}

		if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	if err := writeManifest(outDir); err != nil {
		return err
	}

	fmt.Printf("Generated 20 tool response sets in %s/\n", outDir)
	fmt.Println("Use with: cheaptrick shell --tool-responses", outDir)
	return nil
}

func allTools() []tool {
	var t []tool
	t = append(t, getWeather()...)
	t = append(t, webSearch()...)
	t = append(t, sendEmail()...)
	t = append(t, executeSQL()...)
	t = append(t, readFile()...)
	t = append(t, writeFile()...)
	t = append(t, createCalendarEvent()...)
	t = append(t, getStockPrice()...)
	t = append(t, translateText()...)
	t = append(t, getDirections()...)
	t = append(t, createTicket()...)
	t = append(t, sendSlackMessage()...)
	t = append(t, getExchangeRate()...)
	t = append(t, dnsLookup()...)
	t = append(t, scrapeURL()...)
	t = append(t, getGitHubRepo()...)
	t = append(t, getUserProfile()...)
	t = append(t, httpRequest()...)
	t = append(t, getNews()...)
	t = append(t, runShellCommand()...)
	return t
}

// --------------------------------------------------------------------------
// 1. get_weather — argument-matched by city
// --------------------------------------------------------------------------

func getWeather() []tool {
	return []tool{
		{dir: "get_weather", file: "_match.json", data: map[string]any{
			"field": "city",
			"matches": map[string]string{
				"Paris":     "paris",
				"Tokyo":     "tokyo",
				"New York":  "new_york",
				"London":    "london",
				"São Paulo": "sao_paulo",
			},
			"default": "_default",
		}},
		{dir: "get_weather", file: "paris.json", data: map[string]any{
			"city": "Paris", "country": "FR", "temp_celsius": 18, "temp_fahrenheit": 64,
			"condition": "Partly cloudy", "humidity": 65, "wind_kph": 12,
			"forecast": []map[string]any{
				{"day": "Tomorrow", "high": 21, "low": 14, "condition": "Sunny"},
				{"day": "Day after", "high": 19, "low": 13, "condition": "Rain"},
			},
		}},
		{dir: "get_weather", file: "tokyo.json", data: map[string]any{
			"city": "Tokyo", "country": "JP", "temp_celsius": 28, "temp_fahrenheit": 82,
			"condition": "Sunny", "humidity": 40, "wind_kph": 8,
			"forecast": []map[string]any{
				{"day": "Tomorrow", "high": 30, "low": 24, "condition": "Humid"},
				{"day": "Day after", "high": 27, "low": 22, "condition": "Thunderstorms"},
			},
		}},
		{dir: "get_weather", file: "new_york.json", data: map[string]any{
			"city": "New York", "country": "US", "temp_celsius": 22, "temp_fahrenheit": 72,
			"condition": "Clear", "humidity": 55, "wind_kph": 15,
			"forecast": []map[string]any{
				{"day": "Tomorrow", "high": 25, "low": 18, "condition": "Partly cloudy"},
				{"day": "Day after", "high": 20, "low": 15, "condition": "Overcast"},
			},
		}},
		{dir: "get_weather", file: "london.json", data: map[string]any{
			"city": "London", "country": "GB", "temp_celsius": 14, "temp_fahrenheit": 57,
			"condition": "Overcast", "humidity": 78, "wind_kph": 20,
			"forecast": []map[string]any{
				{"day": "Tomorrow", "high": 16, "low": 10, "condition": "Drizzle"},
				{"day": "Day after", "high": 15, "low": 9, "condition": "Rain"},
			},
		}},
		{dir: "get_weather", file: "sao_paulo.json", data: map[string]any{
			"city": "São Paulo", "country": "BR", "temp_celsius": 25, "temp_fahrenheit": 77,
			"condition": "Warm", "humidity": 70, "wind_kph": 10,
			"forecast": []map[string]any{
				{"day": "Tomorrow", "high": 28, "low": 20, "condition": "Thunderstorms"},
				{"day": "Day after", "high": 26, "low": 19, "condition": "Partly cloudy"},
			},
		}},
		{dir: "get_weather", file: "_default.json", data: map[string]any{
			"city": "{{args.city}}", "country": "XX", "temp_celsius": 20, "temp_fahrenheit": 68,
			"condition": "Unknown", "humidity": 50, "wind_kph": 10,
			"forecast": []map[string]any{},
		}},
	}
}

func webSearch() []tool {
	return []tool{
		{file: "web_search.1.json", data: map[string]any{
			"query":         "{{args.query}}",
			"page":          1,
			"total_results": 42,
			"results": []map[string]any{
				{"title": "Introduction to LLM Agents", "url": "https://example.com/llm-agents", "snippet": "A comprehensive guide to building autonomous LLM-powered agents with tool use."},
				{"title": "Tool Calling in Gemini API", "url": "https://ai.google.dev/docs/tool-calling", "snippet": "Learn how to define and use function calling with the Gemini API."},
				{"title": "Agent Frameworks Compared", "url": "https://example.com/agent-frameworks", "snippet": "A comparison of popular agent frameworks: LangChain, CrewAI, Rig, and AutoGen."},
			},
		}},
		{file: "web_search.2.json", data: map[string]any{
			"query":         "{{args.query}}",
			"page":          2,
			"total_results": 42,
			"results": []map[string]any{
				{"title": "Building Reliable Agents", "url": "https://example.com/reliable-agents", "snippet": "Patterns for error handling, retries, and graceful degradation in LLM agents."},
				{"title": "Mock Testing for AI Applications", "url": "https://example.com/mock-testing-ai", "snippet": "How to build deterministic test suites for non-deterministic AI systems."},
			},
		}},
		{file: "web_search.3.json", data: map[string]any{
			"query":         "{{args.query}}",
			"page":          3,
			"total_results": 42,
			"results":       []map[string]any{},
		}},
	}
}

func sendEmail() []tool {
	return []tool{
		{file: "send_email.json", data: map[string]any{
			"status":     "sent",
			"message_id": "msg_a1b2c3d4e5f6",
			"to":         "{{args.to}}",
			"subject":    "{{args.subject}}",
			"timestamp":  "2025-03-05T10:30:00Z",
			"size_bytes": 2048,
		}},
	}
}

func executeSQL() []tool {
	return []tool{
		{dir: "execute_sql", file: "_match.json", data: map[string]any{
			"field": "query_type",
			"matches": map[string]string{
				"select": "select",
				"insert": "insert",
				"update": "update",
				"delete": "delete",
			},
			"default": "_default",
		}},
		{dir: "execute_sql", file: "select.json", data: map[string]any{
			"status":        "ok",
			"rows_returned": 3,
			"columns":       []string{"id", "name", "email", "created_at"},
			"rows": []map[string]any{
				{"id": 1, "name": "Alice Chen", "email": "alice@example.com", "created_at": "2024-11-01T09:00:00Z"},
				{"id": 2, "name": "Bob Kumar", "email": "bob@example.com", "created_at": "2024-11-15T14:30:00Z"},
				{"id": 3, "name": "Carol Santos", "email": "carol@example.com", "created_at": "2025-01-03T11:15:00Z"},
			},
			"execution_time_ms": 12,
		}},
		{dir: "execute_sql", file: "insert.json", data: map[string]any{
			"status":            "ok",
			"rows_affected":     1,
			"last_insert_id":    47,
			"execution_time_ms": 8,
		}},
		{dir: "execute_sql", file: "update.json", data: map[string]any{
			"status":            "ok",
			"rows_affected":     1,
			"execution_time_ms": 5,
		}},
		{dir: "execute_sql", file: "delete.json", data: map[string]any{
			"status":            "ok",
			"rows_affected":     1,
			"execution_time_ms": 4,
		}},
		{dir: "execute_sql", file: "_default.json", data: map[string]any{
			"status":            "ok",
			"rows_affected":     0,
			"execution_time_ms": 2,
		}},
	}
}

func readFile() []tool {
	return []tool{
		{file: "read_file.json", data: map[string]any{
			"path":          "{{args.path}}",
			"exists":        true,
			"size_bytes":    4096,
			"mime_type":     "text/plain",
			"last_modified": "2025-02-28T16:45:00Z",
			"content":       "# Project README\n\nThis is a sample project for testing LLM agent tool calls.\n\n## Setup\n\n1. Install dependencies\n2. Configure environment\n3. Run tests\n",
		}},
	}
}

func writeFile() []tool {
	return []tool{
		{file: "write_file.json", data: map[string]any{
			"status":        "written",
			"path":          "{{args.path}}",
			"bytes_written": 1024,
			"timestamp":     "2025-03-05T10:35:00Z",
		}},
	}
}

func createCalendarEvent() []tool {
	return []tool{
		{file: "create_calendar_event.json", data: map[string]any{
			"status":   "created",
			"event_id": "evt_x7k9m2p4",
			"title":    "{{args.title}}",
			"start":    "{{args.start_time}}",
			"end":      "{{args.end_time}}",
			"calendar": "primary",
			"link":     "https://calendar.google.com/event?eid=evt_x7k9m2p4",
		}},
	}
}

func getStockPrice() []tool {
	return []tool{
		{dir: "get_stock_price", file: "_match.json", data: map[string]any{
			"field": "symbol",
			"matches": map[string]string{
				"AAPL":  "aapl",
				"GOOGL": "googl",
				"MSFT":  "msft",
			},
			"default": "_default",
		}},
		{dir: "get_stock_price", file: "aapl.json", data: map[string]any{
			"symbol": "AAPL", "name": "Apple Inc.", "price": 187.44, "currency": "USD",
			"change": 2.31, "change_percent": 1.25, "volume": 54_200_000,
			"market_cap": "2.89T", "pe_ratio": 29.8, "day_high": 188.10, "day_low": 184.92,
			"timestamp": "2025-03-05T16:00:00Z",
		}},
		{dir: "get_stock_price", file: "googl.json", data: map[string]any{
			"symbol": "GOOGL", "name": "Alphabet Inc.", "price": 174.82, "currency": "USD",
			"change": -1.15, "change_percent": -0.65, "volume": 28_400_000,
			"market_cap": "2.15T", "pe_ratio": 24.1, "day_high": 176.50, "day_low": 173.20,
			"timestamp": "2025-03-05T16:00:00Z",
		}},
		{dir: "get_stock_price", file: "msft.json", data: map[string]any{
			"symbol": "MSFT", "name": "Microsoft Corp.", "price": 415.60, "currency": "USD",
			"change": 3.80, "change_percent": 0.92, "volume": 22_100_000,
			"market_cap": "3.09T", "pe_ratio": 35.2, "day_high": 417.00, "day_low": 411.30,
			"timestamp": "2025-03-05T16:00:00Z",
		}},
		{dir: "get_stock_price", file: "_default.json", data: map[string]any{
			"symbol": "{{args.symbol}}", "name": "Unknown Corp.", "price": 100.00, "currency": "USD",
			"change": 0.0, "change_percent": 0.0, "volume": 0,
			"error": "Symbol not found in mock data. Add a canned response file for this symbol.",
		}},
	}
}

func translateText() []tool {
	return []tool{
		{dir: "translate_text", file: "_match.json", data: map[string]any{
			"field": "target_language",
			"matches": map[string]string{
				"es": "spanish",
				"fr": "french",
				"de": "german",
				"ja": "japanese",
			},
			"default": "_default",
		}},
		{dir: "translate_text", file: "spanish.json", data: map[string]any{
			"source_language": "en", "target_language": "es",
			"original":   "{{args.text}}",
			"translated": "Este es un texto de ejemplo traducido al español por el servidor simulado.",
			"confidence": 0.95,
		}},
		{dir: "translate_text", file: "french.json", data: map[string]any{
			"source_language": "en", "target_language": "fr",
			"original":   "{{args.text}}",
			"translated": "Ceci est un exemple de texte traduit en français par le serveur simulé.",
			"confidence": 0.94,
		}},
		{dir: "translate_text", file: "german.json", data: map[string]any{
			"source_language": "en", "target_language": "de",
			"original":   "{{args.text}}",
			"translated": "Dies ist ein Beispieltext, der vom Mock-Server ins Deutsche übersetzt wurde.",
			"confidence": 0.93,
		}},
		{dir: "translate_text", file: "japanese.json", data: map[string]any{
			"source_language": "en", "target_language": "ja",
			"original":   "{{args.text}}",
			"translated": "これはモックサーバーによって日本語に翻訳されたサンプルテキストです。",
			"confidence": 0.91,
		}},
		{dir: "translate_text", file: "_default.json", data: map[string]any{
			"source_language": "en", "target_language": "{{args.target_language}}",
			"original":   "{{args.text}}",
			"translated": "[mock translation to {{args.target_language}}]",
			"confidence": 0.50,
		}},
	}
}

// --------------------------------------------------------------------------
// 10. get_directions — static
// --------------------------------------------------------------------------

func getDirections() []tool {
	return []tool{
		{file: "get_directions.json", data: map[string]any{
			"origin":        "{{args.origin}}",
			"destination":   "{{args.destination}}",
			"distance_km":   12.4,
			"duration_min":  28,
			"mode":          "driving",
			"route_summary": "Via A1 Highway",
			"steps": []map[string]any{
				{"instruction": "Head north on Main St", "distance_km": 0.5, "duration_min": 2},
				{"instruction": "Turn right onto Highway A1 ramp", "distance_km": 0.3, "duration_min": 1},
				{"instruction": "Continue on A1 Highway for 10 km", "distance_km": 10.0, "duration_min": 18},
				{"instruction": "Take exit 14 toward Downtown", "distance_km": 0.8, "duration_min": 3},
				{"instruction": "Arrive at destination on your right", "distance_km": 0.8, "duration_min": 4},
			},
			"tolls":          false,
			"traffic_status": "moderate",
		}},
	}
}

func createTicket() []tool {
	return []tool{
		{file: "create_ticket.json", data: map[string]any{
			"status":     "created",
			"ticket_id":  "PROJ-1042",
			"title":      "{{args.title}}",
			"priority":   "{{args.priority}}",
			"assignee":   "{{args.assignee}}",
			"project":    "{{args.project}}",
			"url":        "https://jira.example.com/browse/PROJ-1042",
			"created_at": "2025-03-05T10:40:00Z",
		}},
	}
}

func sendSlackMessage() []tool {
	return []tool{
		{file: "send_slack_message.json", data: map[string]any{
			"status":     "sent",
			"channel":    "{{args.channel}}",
			"message_ts": "1709632800.001200",
			"permalink":  "https://workspace.slack.com/archives/C01ABC/p1709632800001200",
			"timestamp":  "2025-03-05T10:40:00Z",
		}},
	}
}

func getExchangeRate() []tool {
	return []tool{
		{dir: "get_exchange_rate", file: "_match.json", data: map[string]any{
			"field": "to",
			"matches": map[string]string{
				"EUR": "eur",
				"GBP": "gbp",
				"JPY": "jpy",
			},
			"default": "_default",
		}},
		{dir: "get_exchange_rate", file: "eur.json", data: map[string]any{
			"from": "USD", "to": "EUR", "rate": 0.9215, "inverse_rate": 1.0852,
			"timestamp": "2025-03-05T12:00:00Z", "source": "ECB",
		}},
		{dir: "get_exchange_rate", file: "gbp.json", data: map[string]any{
			"from": "USD", "to": "GBP", "rate": 0.7890, "inverse_rate": 1.2674,
			"timestamp": "2025-03-05T12:00:00Z", "source": "ECB",
		}},
		{dir: "get_exchange_rate", file: "jpy.json", data: map[string]any{
			"from": "USD", "to": "JPY", "rate": 149.82, "inverse_rate": 0.00668,
			"timestamp": "2025-03-05T12:00:00Z", "source": "ECB",
		}},
		{dir: "get_exchange_rate", file: "_default.json", data: map[string]any{
			"from": "{{args.from}}", "to": "{{args.to}}", "rate": 1.0, "inverse_rate": 1.0,
			"error": "Currency pair not found in mock data.",
		}},
	}
}

func dnsLookup() []tool {
	return []tool{
		{file: "dns_lookup.json", data: map[string]any{
			"domain": "{{args.domain}}",
			"records": map[string]any{
				"A":     []string{"93.184.216.34"},
				"AAAA":  []string{"2606:2800:220:1:248:1893:25c8:1946"},
				"MX":    []map[string]any{{"priority": 10, "host": "mail.example.com"}},
				"NS":    []string{"ns1.example.com", "ns2.example.com"},
				"TXT":   []string{"v=spf1 include:_spf.example.com ~all"},
				"CNAME": nil,
			},
			"ttl":             3600,
			"query_time_ms":   24,
			"nameserver_used": "8.8.8.8",
		}},
	}
}

func scrapeURL() []tool {
	return []tool{
		{file: "scrape_url.json", data: map[string]any{
			"url":              "{{args.url}}",
			"status_code":      200,
			"title":            "Example Domain",
			"meta_description": "This domain is for use in illustrative examples in documents.",
			"content_type":     "text/html",
			"content_length":   1256,
			"text_content":     "Example Domain\n\nThis domain is for use in illustrative examples in documents. You may use this domain in literature without prior coordination or asking for permission.\n\nMore information...",
			"links": []map[string]any{
				{"text": "More information...", "href": "https://www.iana.org/domains/example"},
			},
			"headers": map[string]string{
				"content-type":  "text/html; charset=UTF-8",
				"cache-control": "max-age=604800",
				"server":        "ECS (dcb/7F83)",
			},
			"fetch_time_ms": 145,
		}},
	}
}

func getGitHubRepo() []tool {
	return []tool{
		{dir: "get_github_repo", file: "_match.json", data: map[string]any{
			"field": "repo",
			"matches": map[string]string{
				"charmbracelet/bubbletea": "bubbletea",
				"0xPlaygrounds/rig":       "rig",
			},
			"default": "_default",
		}},
		{dir: "get_github_repo", file: "bubbletea.json", data: map[string]any{
			"full_name": "charmbracelet/bubbletea", "description": "A powerful little TUI framework",
			"language": "Go", "stars": 28400, "forks": 820, "open_issues": 65,
			"license": "MIT", "default_branch": "master",
			"created_at": "2020-01-15T00:00:00Z", "updated_at": "2025-03-04T18:00:00Z",
			"topics":    []string{"tui", "terminal", "go", "elm-architecture", "cli"},
			"clone_url": "https://github.com/charmbracelet/bubbletea.git",
		}},
		{dir: "get_github_repo", file: "rig.json", data: map[string]any{
			"full_name": "0xPlaygrounds/rig", "description": "Rust library for building LLM-powered applications",
			"language": "Rust", "stars": 4200, "forks": 310, "open_issues": 42,
			"license": "MIT", "default_branch": "main",
			"created_at": "2024-03-20T00:00:00Z", "updated_at": "2025-03-03T12:00:00Z",
			"topics":    []string{"llm", "rust", "ai", "agents", "rag"},
			"clone_url": "https://github.com/0xPlaygrounds/rig.git",
		}},
		{dir: "get_github_repo", file: "_default.json", data: map[string]any{
			"full_name": "{{args.repo}}", "description": "Mock repository",
			"language": "Unknown", "stars": 0, "forks": 0, "open_issues": 0,
			"error": "Repository not found in mock data.",
		}},
	}
}

func getUserProfile() []tool {
	return []tool{
		{dir: "get_user_profile", file: "_match.json", data: map[string]any{
			"field": "user_id",
			"matches": map[string]string{
				"usr_001": "alice",
				"usr_002": "bob",
			},
			"default": "_default",
		}},
		{dir: "get_user_profile", file: "alice.json", data: map[string]any{
			"user_id": "usr_001", "name": "Alice Chen", "email": "alice@example.com",
			"role": "admin", "department": "Engineering",
			"created_at": "2023-06-15T09:00:00Z", "last_login": "2025-03-05T08:30:00Z",
			"preferences": map[string]any{"timezone": "America/New_York", "locale": "en-US", "theme": "dark"},
		}},
		{dir: "get_user_profile", file: "bob.json", data: map[string]any{
			"user_id": "usr_002", "name": "Bob Kumar", "email": "bob@example.com",
			"role": "member", "department": "Product",
			"created_at": "2024-01-10T14:00:00Z", "last_login": "2025-03-04T17:45:00Z",
			"preferences": map[string]any{"timezone": "Asia/Kolkata", "locale": "en-IN", "theme": "light"},
		}},
		{dir: "get_user_profile", file: "_default.json", data: map[string]any{
			"user_id": "{{args.user_id}}", "name": "Unknown User",
			"error": "User not found in mock data.",
		}},
	}
}

func httpRequest() []tool {
	return []tool{
		{dir: "http_request", file: "_match.json", data: map[string]any{
			"field": "method",
			"matches": map[string]string{
				"GET":    "get",
				"POST":   "post",
				"PUT":    "put",
				"DELETE": "delete",
			},
			"default": "get",
		}},
		{dir: "http_request", file: "get.json", data: map[string]any{
			"status_code": 200,
			"headers": map[string]string{
				"content-type": "application/json",
				"x-request-id": "req_mock_abc123",
			},
			"body":             map[string]any{"message": "OK", "data": map[string]any{"id": 1, "status": "active"}},
			"response_time_ms": 85,
		}},
		{dir: "http_request", file: "post.json", data: map[string]any{
			"status_code": 201,
			"headers": map[string]string{
				"content-type": "application/json",
				"location":     "/api/resources/42",
			},
			"body":             map[string]any{"message": "Created", "id": 42},
			"response_time_ms": 120,
		}},
		{dir: "http_request", file: "put.json", data: map[string]any{
			"status_code":      200,
			"headers":          map[string]string{"content-type": "application/json"},
			"body":             map[string]any{"message": "Updated"},
			"response_time_ms": 95,
		}},
		{dir: "http_request", file: "delete.json", data: map[string]any{
			"status_code":      204,
			"headers":          map[string]string{},
			"body":             nil,
			"response_time_ms": 60,
		}},
	}
}

func getNews() []tool {
	return []tool{
		{file: "get_news.1.json", data: map[string]any{
			"query":    "{{args.query}}",
			"page":     1,
			"has_more": true,
			"articles": []map[string]any{
				{"title": "AI Agents Are Reshaping Software Development", "source": "TechCrunch", "published": "2025-03-05T08:00:00Z", "url": "https://techcrunch.com/example-1"},
				{"title": "Open Source LLM Tooling Reaches Inflection Point", "source": "The Verge", "published": "2025-03-04T14:00:00Z", "url": "https://theverge.com/example-2"},
				{"title": "New Benchmarks Show Rapid Progress in Tool Use", "source": "ArXiv Blog", "published": "2025-03-03T10:00:00Z", "url": "https://arxiv.org/example-3"},
			},
		}},
		{file: "get_news.2.json", data: map[string]any{
			"query":    "{{args.query}}",
			"page":     2,
			"has_more": false,
			"articles": []map[string]any{
				{"title": "Mock Servers: The Unsung Hero of Agent Testing", "source": "Dev.to", "published": "2025-03-01T12:00:00Z", "url": "https://dev.to/example-4"},
			},
		}},
		{file: "get_news.3.json", data: map[string]any{
			"query":    "{{args.query}}",
			"page":     3,
			"has_more": false,
			"articles": []map[string]any{},
		}},
	}
}

func runShellCommand() []tool {
	return []tool{
		{dir: "run_shell_command", file: "_match.json", data: map[string]any{
			"field": "command",
			"matches": map[string]string{
				"ls":     "ls",
				"pwd":    "pwd",
				"whoami": "whoami",
			},
			"default": "_default",
		}},
		{dir: "run_shell_command", file: "ls.json", data: map[string]any{
			"exit_code": 0,
			"stdout":    "README.md\nMakefile\ncmd/\ninternal/\ngo.mod\ngo.sum\nmain.go\n",
			"stderr":    "",
		}},
		{dir: "run_shell_command", file: "pwd.json", data: map[string]any{
			"exit_code": 0,
			"stdout":    "/home/developer/projects/my-agent\n",
			"stderr":    "",
		}},
		{dir: "run_shell_command", file: "whoami.json", data: map[string]any{
			"exit_code": 0,
			"stdout":    "developer\n",
			"stderr":    "",
		}},
		{dir: "run_shell_command", file: "_default.json", data: map[string]any{
			"exit_code": 0,
			"stdout":    "[mock output for: {{args.command}}]\n",
			"stderr":    "",
		}},
	}
}

func writeManifest(outDir string) error {
	content := `# Mock Tools Manifest

Generated by ` + "`cheaptrick tools`" + `. Use with:

` + "```" + `
cheaptrick shell --tool-responses ` + outDir + `
` + "```" + `

## Tools

| # | Function | Type | Description |
|---|----------|------|-------------|
| 1 | ` + "`get_weather`" + ` | Argument-matched (city) | Weather data for Paris, Tokyo, New York, London, São Paulo |
| 2 | ` + "`web_search`" + ` | Sequenced (3 pages) | Paginated search results, third page returns empty |
| 3 | ` + "`send_email`" + ` | Static | Confirms email sent with message ID |
| 4 | ` + "`execute_sql`" + ` | Argument-matched (query_type) | SELECT returns rows, INSERT/UPDATE/DELETE return affected counts |
| 5 | ` + "`read_file`" + ` | Static | Returns sample file content with metadata |
| 6 | ` + "`write_file`" + ` | Static | Confirms bytes written with timestamp |
| 7 | ` + "`create_calendar_event`" + ` | Static | Confirms event creation with link |
| 8 | ` + "`get_stock_price`" + ` | Argument-matched (symbol) | Live-style quotes for AAPL, GOOGL, MSFT |
| 9 | ` + "`translate_text`" + ` | Argument-matched (target_language) | Translations to Spanish, French, German, Japanese |
| 10 | ` + "`get_directions`" + ` | Static | Turn-by-turn driving directions |
| 11 | ` + "`create_ticket`" + ` | Static | Jira-style ticket creation with URL |
| 12 | ` + "`send_slack_message`" + ` | Static | Confirms message posted with permalink |
| 13 | ` + "`get_exchange_rate`" + ` | Argument-matched (to) | USD→EUR, GBP, JPY exchange rates |
| 14 | ` + "`dns_lookup`" + ` | Static | A, AAAA, MX, NS, TXT records |
| 15 | ` + "`scrape_url`" + ` | Static | Page title, text content, links, headers |
| 16 | ` + "`get_github_repo`" + ` | Argument-matched (repo) | Repo metadata for bubbletea and rig |
| 17 | ` + "`get_user_profile`" + ` | Argument-matched (user_id) | User profiles for Alice and Bob |
| 18 | ` + "`http_request`" + ` | Argument-matched (method) | GET/POST/PUT/DELETE with status codes |
| 19 | ` + "`get_news`" + ` | Sequenced (3 pages) | Paginated news articles, third page returns empty |
| 20 | ` + "`run_shell_command`" + ` | Argument-matched (command) | Output for ls, pwd, whoami |

## Response Types

- **Static**: Single ` + "`<function>.json`" + ` file, same response every time.
- **Argument-matched**: ` + "`<function>/`" + ` directory with ` + "`_match.json`" + ` routing file.
  Selects response based on a top-level argument value (case-insensitive).
  Falls back to ` + "`_default.json`" + ` when no rule matches.
- **Sequenced**: ` + "`<function>.N.json`" + ` files (1-indexed). Returns different
  responses on successive calls. Clamps to highest N on overflow.

## Argument Substitution

Files marked with ` + "`{{args.field}}`" + ` placeholders will have those values
replaced with actual function call arguments at resolution time. This
lets a single canned file produce dynamic responses without needing
separate files per argument value.

## Adding Your Own

Create a new JSON file following the naming conventions above. The shell
picks up files at startup — no restart needed if you add files between
conversations (use ` + "`/clear`" + ` to reset call counters).
`

	path := filepath.Join(outDir, "MANIFEST.md")
	return os.WriteFile(path, []byte(content), 0o644)
}
