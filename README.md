<div align="center">
  <h1>🎭 Cheaptrick</h1>
  <p><strong>A locally-hosted, interactive Mock Server and TUI for the Google Gemini API</strong></p>
</div>

---

**Cheaptrick** is a Terminal User Interface (TUI) application built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) that acts as a mock backend for Google Gemini v1beta APIs. It allows developers to intercept API calls, inspect them, and provide mocked JSON responses (or auto-reply from fixtures) during local development, saving time and API costs.

## ✨ Features

- **Intercept Requests**: Catches Gemini API requests sent to your local server.
- **TUI Interface**: View pending requests and their full JSON payloads in an elegant terminal UI.
- **Custom Responses**: Compose custom JSON responses dynamically and send them back to the client.
- **Pre-built Templates**: Use shortcuts to quickly generate standard response skeletons (Text, Function Call, Error 429, Error 500).
- **Auto-Fixtures**: Save responses as templates for specific requests. Future identical requests will be answered automatically.
- **Interactive Shell**: A built-in REPL-like shell allows you to converse with the mock server manually to test your fixtures.
- **HTTPS Support**: Spin up the mock server with TLS to mimic secure environments.
- **Request Logging**: Save request-response pairs to a JSONL log file for debugging.

---

## 🚀 Installation

Ensure you have **Go 1.22+** installed on your system. 

### Option 1: Clone and Build (Recommended)

```bash
git clone https://github.com/yourusername/cheaptrick.git
cd cheaptrick

# Build the binary into the bin/ directory
make build
```

Then you can run `./bin/cheaptrick` to start the application.

### Option 2: Install via Go

Install directly to your `$GOPATH/bin`:

```bash
make install
# OR
go install ./...
```

Verify the installation:

```bash
cheaptrick --help
```

---

## 📖 Usage The App

The application provides three main subcommands: `start`, `fixtures`, and `shell`.

### `cheaptrick start`

Starts the HTTP server and opens the Bubble Tea TUI.

```bash
# Basic usage (defaults to localhost:8080)
cheaptrick start

# With specific port and fixture directory
cheaptrick start --port 9090 --fixtures ./my_fixtures

# With HTTPS enabled
cheaptrick start --tls-cert cert.pem --tls-key key.pem
```

#### TUI Workflow:

1. **Send a Test Request**
   Open a secondary terminal and send a test `curl` request:
   ```bash
   curl -X POST http://localhost:8080/v1beta/models/gemini-2.0-flash:generateContent \
   -H "Content-Type: application/json" \
   -d '{"contents":[{"parts":[{"text":"Hello Gemini!"}]}]}'
   ```

2. **Provide a Custom Response via TUI**
   - When the request arrives, your `curl` command will hang while waiting for your server to respond. In the TUI window, you'll see a `[PENDING]` request item.
   - Hit `Enter` to focus the Response Composer in the TUI.
   - Press `F1` to insert a ready-to-use template, or write any valid JSON response manually.
   - Press `Ctrl+S` to send the response.
   - The `curl` operation finishes and prints out the JSON you entered.

3. **Auto-Fixture Replay Feature**
   - Resend the exact same `curl` request.
   - Focus the Request in the Request List.
   - Press `Ctrl+F` to save the active response as a fixture. The notification bar will display: **"Saved fixture [HASH]"**.
   - Send the `curl` request one more time. The TUI will state **"[req-id] auto-replied from fixture [HASH]"** and `curl` will receive a response immediately!

### `cheaptrick shell`

Starts an interactive shell connecting to the mock server. Use this to quickly test your local mock server's fixtures without needing to construct `curl` payloads manually. It inherently uses the official `google.golang.org/genai` library pointing to the local mock server.

```bash
# Start the interactive shell (defaults to localhost:8080)
cheaptrick shell

# Configure the target server and model
cheaptrick shell --host 127.0.0.1 --port 9090 --model gemini-1.5-pro

# Use a specific history file and API key via environment variable
GEMINI_API_KEY="custom-key" cheaptrick shell --history-file ./my_history.txt
```

**Available Flags:**
- `-H, --host string`: Host address of the mock server (default "localhost")
- `-p, --port int`: Port of the mock server (default 8080)
- `-m, --model string`: Gemini model to use in requests (default "gemini-2.0-flash")
- `--history-file string`: Path to the readline history file (defaults to OS temp directory)

**Environment Variables:**
- `GEMINI_API_KEY`: API key to use (cheaptrick skips validation) (default "mock-key")

### `cheaptrick fixtures`

Generates 30 predefined JSON fixture files and a `MANIFEST.md` file for common text and tool-call prompts. This provides you with an instant set of fixtures to test against, removing the need to build the initial mock state manually.

```bash
# Generate fixtures in the default "fixtures" directory
cheaptrick fixtures

# Generate fixtures in a specific directory
cheaptrick fixtures --output-dir ./test_assets
```

---

## ⌨️ TUI Keybindings

- **`Tab`**: Switch focus between panels (List vs Viewer vs Composer)
- **`j`/`k` or `↑`/`↓`**: Navigate lists and scroll detail view
- **`Enter`**: Open response composer for a pending request
- **`F1`**: Insert Text response skeleton
- **`F2`**: Insert FunctionCall response skeleton
- **`F3`**: Insert Error 429 response skeleton
- **`F4`**: Insert Error 500 response skeleton
- **`Ctrl+S`**: Send composed response
- **`Ctrl+F`**: Save current response as an auto-fixture
- **`Esc`**: Exit compose mode or cancel
- **`q`**: Quit application

---

## 🤝 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.
