# Gemini Mock Server TUI

A terminal user interface (TUI) application built with Bubble Tea that acts as a mock backend for Google Gemini v1beta APIs. It allows developers to intercept API calls, inspect them, and provide mocked JSON responses (or auto-reply from fixtures) during local development.

## Features

- **Intercept Requests**: Catches Gemini API requests sent to your local server.
- **TUI Interface**: View pending requests and their full JSON payloads in an elegant terminal UI.
- **Custom Responses**: Compose custom JSON responses dynamically and send them back to the client.
- **Pre-built Templates**: Use shortcuts to quickly generate standard response skeletons (Text, Function Call, Error 429, Error 500).
- **Auto-Fixtures**: Save responses as templates for specific requests. Future identical requests will be answered automatically.
- **HTTPS Support**: Spin up the mock server with TLS to mimic secure environments.
- **Request Logging**: Save request-response pairs to a JSONL log file for debugging.

## Installation

Ensure you have Go 1.22+ installed, then simply clone the repository and build:

```bash
git clone <repository_url>
cd cheaptrick
make build
```

This will create a binary in the `bin/` directory.

Alternatively, you can install it globally to your `$GOPATH/bin`:

```bash
make install
```

## Usage & Verification

The application now provides two main subcommands: `start` and `fixtures`.

### `cheaptrick start`

Starts the HTTP server and opens the Bubble Tea TUI.

```bash
# Basic usage
./bin/cheaptrick start

# With specific port and fixture directory
./bin/cheaptrick start --port 9090 --fixtures ./my_fixtures

# With HTTPS enabled
./bin/cheaptrick start --tls-cert cert.pem --tls-key key.pem
```

#### Verification & TUI flow:

1. **Send a Test Request**
   Open a secondary terminal and fire a test `curl` request:
   ```bash
   curl -X POST http://localhost:8080/v1beta/models/gemini-2.0-flash:generateContent \
   -H "Content-Type: application/json" \
   -d '{"contents":[{"parts":[{"text":"Hello Gemini TUI!"}]}]}'
   ```

2. **Provide a Custom Response via TUI**
   - When the request arrives, your `curl` command will hang while waiting for your server to respond. In the TUI window, you'll see a `[PENDING]` request item.
   - Hit `Enter` to focus the Response Composer in the TUI.
   - Press `F1` to insert a ready-to-use template, or write any valid JSON response manually.
   - Press `Ctrl+S` to send the response.
   - Check your secondary terminal — the `curl` operation should finish and print out the JSON you entered.

3. **Test the Auto-Fixture Replay Feature**
   - Resend the exact same `curl` request you sent earlier.
   - Go to the Response Composer again.
   - Press `Ctrl+F` to save the active response as a fixture. The notification bar will display: **"Saved fixture [HASH]"**.
   - Send the `curl` request one more time.
   - The TUI notification bar should rapidly update stating **"[req-id] auto-replied from fixture [HASH]"** and `curl` will receive a response immediately without you typing anything!

### `cheaptrick fixtures`

Generates 30 predefined JSON fixture files and a `MANIFEST.md` file for common text and tool-call prompts.

```bash
# Generate fixtures in the default "fixtures" directory
./bin/cheaptrick fixtures

# Generate fixtures in a specific directory
./bin/cheaptrick fixtures --output-dir ./test_assets
```

## Keybindings

- **`Tab`**: Switch focus between panels
- **`j`/`k` or `↑`/`↓`**: Navigate lists and scroll detail view
- **`Enter`**: Open response composer
- **`F1`**: Insert Text response skeleton
- **`F2`**: Insert FunctionCall response skeleton
- **`F3`**: Insert Error 429 response skeleton
- **`F4`**: Insert Error 500 response skeleton
- **`Ctrl+S`**: Send composed response
- **`Ctrl+F`**: Save current response as an auto-fixture
- **`Esc`**: Exit compose mode or cancel
- **`q`**: Quit application
