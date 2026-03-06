<div align="center">

# 🎭 Cheaptrick

**A human-in-the-loop mock server for the Google Gemini API**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](CONTRIBUTING.md)
[![Built with Bubble Tea](https://img.shields.io/badge/Built_with-Bubble_Tea-ff69b4?style=flat-square)](https://github.com/charmbracelet/bubbletea)

Cheaptrick intercepts Google Gemini API requests, lets you craft responses
by hand, and replays them from fixtures — so you can develop and debug
LLM-powered agents locally without spending tokens.

[Getting Started](#getting-started) · [Usage](#usage) · [Tool-Call Debugging](#tool-call-debugging-with-the-shell) · [SDK Examples](#connecting-your-app) · [Keybindings](#keybindings) · [Contributing](#contributing)

</div>

---

## The Problem

Developing an LLM agent means hundreds of API round-trips: testing parsing
logic, iterating on tool-call schemas, handling edge cases in multi-turn
flows. Each call costs tokens and hits rate limits, but the responses you
need during development are often predictable.

Cheaptrick replaces the Gemini API with a local endpoint where you control
every response. Requests appear in a TUI or web dashboard. You compose a
JSON response (or pick a template), send it, and your client receives it
as if it came from Gemini. Save a response as a fixture, and identical
future requests are replayed automatically.

The result: deterministic, reproducible agent development at zero API cost.

---

## Features

- **TUI and Web UI** — Full Bubble Tea terminal interface or a React-based
  web dashboard (embedded in the binary, no Node.js required). Both
  provide request inspection, response composing, and fixture management.
- **Response templates** — Pre-built skeletons for text responses, function
  calls, 429 rate-limit errors, and 500 server errors.
- **Fixture replay** — Save any response as a fixture keyed by request
  hash. Matching requests are auto-replied without human intervention.
- **Interactive shell** — A REPL backed by the official `google.golang.org/genai`
  client that sends prompts to your mock server and handles tool-call
  loops with canned response files.
- **TLS support** — Serve over HTTPS when your client requires it.
- **JSONL logging** — Every request/response pair is logged for post-hoc
  analysis.
- **Drop-in SDK compatibility** — Works with the official Gemini SDKs for
  Go, Python, TypeScript, and Rust. Just change the base URL.

---

## Getting Started

### Prerequisites

- Go 1.22 or later

### Build from source

```bash
git clone https://github.com/yourusername/cheaptrick.git
cd cheaptrick
make build
./bin/cheaptrick --help
```

Or install to `$GOPATH/bin`:

```bash
make install
cheaptrick --help
```

---

## Usage

Cheaptrick provides four subcommands: `start`, `web`, `shell`, and
`fixtures`.

### `start` — Mock server with TUI

```bash
cheaptrick start
cheaptrick start --port 9090 --fixtures ./my_fixtures
cheaptrick start --tls-cert cert.pem --tls-key key.pem
```

Send a request from another terminal:

```bash
curl -s -X POST \
  http://localhost:8080/v1beta/models/gemini-2.0-flash:generateContent \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"Hello Gemini!"}]}]}'
```

The request blocks. In the TUI, select the pending request, press
**Enter**, compose a response (or press **F1** for a text template), and
**Ctrl+S** to send. The curl command returns with your response.

Press **Ctrl+F** to save the response as a fixture. Subsequent identical
requests are auto-replied.

### `web` — Mock server with browser UI

```bash
cheaptrick web
cheaptrick web --port 9090 --web-port 4000 --open
```

Starts the mock API on `:8080` and serves the React frontend on `:3000`.
The frontend is embedded in the binary — no separate build step or
Node.js installation is needed. Functionality is identical to the TUI:
request list, detail view, response composer, templates, and fixtures.

### `shell` — Interactive REPL

```bash
cheaptrick shell
cheaptrick shell --host 127.0.0.1 --port 9090 --model gemini-1.5-pro
```

The shell sends `GenerateContent` requests to your running mock server
using the official genai client. It maintains conversation history across
turns, detects function-call responses, and supports canned tool
responses for automated tool-call loop testing.

| Flag | Default | Description |
|------|---------|-------------|
| `-H, --host` | `localhost` | Mock server host |
| `-p, --port` | `8080` | Mock server port |
| `-m, --model` | `gemini-2.0-flash` | Model name in requests |
| `--tool-responses` | | Directory for canned tool response files |
| `--auto` | `false` | Auto-send canned responses without prompting |
| `--max-turns` | `20` | Maximum tool-call loop iterations |
| `--history-file` | OS temp dir | Readline history path |

`GEMINI_API_KEY` can be set but is not validated by the mock server.

### `fixtures` — Generate starter fixtures

```bash
cheaptrick fixtures
cheaptrick fixtures --output-dir ./test_assets
```

Generates 30 predefined fixture files for common text and tool-call
prompts, plus a `MANIFEST.md` index.

---

## Tool-Call Debugging with the Shell

The shell is designed for debugging multi-turn tool-calling flows. In this
setup, the developer operates both sides of the conversation:

- **The shell** acts as the tool executor — it resolves canned responses
  from disk and sends `FunctionResponse` parts automatically.
- **The developer** (via the TUI or web UI) acts as the model — crafting
  `FunctionCall` or text responses for each turn.

This two-terminal workflow lets you simulate complete agentic loops,
control every decision boundary, and build fixture libraries that capture
multi-step flows.

### How it works

```
 Terminal 1 (shell)                    Terminal 2 (TUI or browser)
 ──────────────────                    ────────────────────────────
 Prompt> What's the weather?
   → GenerateContent sent to :8080
   → blocks, waiting...
                                       [PENDING] req-01 appears
                                       Press F2, compose:
                                         get_weather(city="Paris")
                                       Ctrl+S to send

 Receives FunctionCall:
   get_weather(city="Paris")
 Resolves: mock_tools/get_weather/paris.json
   → {"temp": 18, "condition": "cloudy"}
 [Enter] to accept

 FunctionResponse sent,
   next GenerateContent to :8080
   → blocks again...
                                       [PENDING] req-02 appears
                                       (request shows FunctionResponse
                                        in the conversation history)
                                       Press F1, type:
                                         "It's 18°C and cloudy."
                                       Ctrl+S to send

 Receives text response:
   "It's 18°C and cloudy."
 Prints to terminal. Chain complete.

 ─── Chain Complete (3 turns, 12.4s) ──────────
  1  user  │ What's the weather?
  2  mock  │ ƒ get_weather(city="Paris")     [via TUI]
  3  tool  │ → {"temp":18,"condition":"cloudy"} [canned: paris.json]
  4  mock  │ It's 18°C and cloudy.           [via TUI]
 ─────────────────────────────────────────────

 Prompt> _
```

### Setting up canned tool responses

Create a directory for your tool response files and pass it to the shell
with `--tool-responses`:

```bash
cheaptrick shell --tool-responses ./mock_tools
```

The shell resolves canned responses using the following priority:

1. **Argument-matched** — `<dir>/<function>/_match.json` routes to
   different files based on argument values.
2. **Sequenced** — `<dir>/<function>.N.json` returns different responses
   for the Nth call to that function.
3. **Static** — `<dir>/<function>.json` returns the same response every
   time.
4. **Subdirectory default** — `<dir>/<function>/_default.json` as a
   fallback.
5. **Manual** — If no file matches, the shell displays the function call
   as structured JSON with a warning listing every path that was checked,
   and prompts for manual input.

#### Static response

A single file per function, returned regardless of arguments.

```
mock_tools/
  web_search.json
```

```json
{
  "results": [
    {"title": "Example", "url": "https://example.com", "snippet": "..."}
  ]
}
```

#### Argument-matched responses

Route to different files based on a top-level argument field.

```
mock_tools/
  get_weather/
    _match.json
    paris.json
    tokyo.json
    _default.json
```

**`_match.json`** — routing rule:

```json
{
  "field": "city",
  "matches": {
    "Paris": "paris",
    "Tokyo": "tokyo"
  },
  "default": "_default"
}
```

Matching is case-insensitive. The `default` key names the fallback file
when no rule matches.

**`paris.json`**:

```json
{"temp": 18, "condition": "cloudy", "humidity": 65}
```

**`tokyo.json`**:

```json
{"temp": 28, "condition": "sunny", "humidity": 40}
```

**`_default.json`**:

```json
{"temp": 20, "condition": "unknown", "humidity": 50}
```

When the developer crafts a `get_weather(city="Paris")` function call in
the TUI, the shell matches on the `city` field and loads `paris.json`.

#### Sequenced responses

Different response per call count, useful for pagination or polling flows.

```
mock_tools/
  fetch_page.1.json    # returned on 1st call
  fetch_page.2.json    # returned on 2nd call
  fetch_page.3.json    # returned on 3rd+ calls (clamped to highest)
```

**`fetch_page.1.json`**:

```json
{"items": [{"id": 1}, {"id": 2}], "has_more": true}
```

**`fetch_page.2.json`**:

```json
{"items": [{"id": 3}, {"id": 4}], "has_more": true}
```

**`fetch_page.3.json`**:

```json
{"items": [], "has_more": false}
```

Call counters reset on `/clear`.

#### Argument substitution

Canned files support `{{args.fieldname}}` placeholders that are replaced
with actual argument values at resolution time:

**`create_record.json`**:

```json
{
  "id": "{{args.id}}",
  "name": "{{args.name}}",
  "status": "created"
}
```

If the function call is `create_record(id="abc-123", name="Test")`, the
shell sends:

```json
{
  "id": "abc-123",
  "name": "Test",
  "status": "created"
}
```

Missing fields leave the placeholder as-is and print a warning.

### Step mode and auto mode

In **step mode** (default), the shell pauses at every function call and
shows the resolved canned response for review:

```
┌─ FUNCTION CALL ──────────────────────────────────┐
│  get_weather                                     │
│                                                  │
│  Arguments:                                      │
│    city: "Paris"                                 │
│    units: "celsius"                              │
│                                                  │
│  Canned: mock_tools/get_weather/paris.json       │
│  {"temp": 18, "condition": "cloudy"}             │
│                                                  │
│  [Enter] Accept  [e] Edit  [s] Accept & Save    │
│  [t] Type new    [x] Abort chain                 │
└──────────────────────────────────────────────────┘
```

When no canned response is found, the shell shows the full function call
with a warning listing every path it checked:

```
┌─ FUNCTION CALL ──────────────────────────────────┐
│  analyze_sentiment                               │
│                                                  │
│  Arguments:                                      │
│  {                                               │
│    "text": "I love this product",                │
│    "language": "en"                              │
│  }                                               │
│                                                  │
│  ⚠ No canned response found.                    │
│  Checked:                                        │
│    ✗ mock_tools/analyze_sentiment/_match.json    │
│    ✗ mock_tools/analyze_sentiment.1.json         │
│    ✗ mock_tools/analyze_sentiment.json           │
│    ✗ mock_tools/analyze_sentiment/_default.json  │
│                                                  │
│  Type the tool's return value as JSON:           │
└──────────────────────────────────────────────────┘
Tool Response> {"sentiment": "positive", "score": 0.95}
Save as mock_tools/analyze_sentiment.json? [y/N]
```

In **auto mode** (`--auto` flag or `/auto` command), canned responses are
sent immediately without pausing. The shell falls back to step mode for
any function call without a matching file.

```
  → get_weather(city="Paris") ← mock_tools/get_weather/paris.json [auto]
  → get_weather(city="Tokyo") ← mock_tools/get_weather/tokyo.json [auto]
```

### Failure injection

Test error handling by injecting failures before the mock server sends
a function call:

```
/fail get_weather           # next call returns an error (one-shot)
/fail get_weather persist   # every call errors until /unfail get_weather
/timeout fetch_page 10      # next call delays 10s before responding
```

The injected error `FunctionResponse`:

```json
{"error": "Service unavailable: get_weather failed (injected by /fail)"}
```

### Exporting conversations as fixture sequences

After completing a multi-turn flow, export the full conversation for
deterministic replay:

```
/export weather_flow
```

Creates:

```
fixtures/weather_flow/
  001_request.json      # 1st GenerateContent request body
  001_response.json     # FunctionCall response from mock
  002_request.json      # 2nd GenerateContent (includes FunctionResponse)
  002_response.json     # 2nd FunctionCall or text response
  003_request.json      # 3rd GenerateContent
  003_response.json     # Final text response
  manifest.json         # model, timestamp, turn count, file index
```

These sequences can be loaded by `cheaptrick start` or `cheaptrick web`
to replay entire multi-turn flows without human input.

### Shell commands

| Command | Description |
|---------|-------------|
| `/clear` | Reset conversation history and call counters |
| `/history` | Print the full conversation as formatted turns |
| `/trace` | Reprint the trace of the last completed chain |
| `/auto` | Switch to auto mode |
| `/step` | Switch to step mode |
| `/tools` | List canned tool response files found in `--tool-responses` dir |
| `/fail <fn>` | Inject an error for the next call to `<fn>` |
| `/fail <fn> persist` | Inject errors for all calls to `<fn>` |
| `/unfail <fn>` | Remove persistent failure injection |
| `/timeout <fn> <sec>` | Delay the next call to `<fn>` by N seconds |
| `/export <name>` | Export conversation as a numbered fixture sequence |
| `/help` | List all commands |
| `/quit` | Exit the shell |

---

## Connecting Your App

Point any Gemini SDK at `http://localhost:8080` with any API key string.

### Go

```go
client, _ := genai.NewClient(ctx, &genai.ClientConfig{
    HTTPOptions: genai.HTTPOptions{
        BaseURL: "http://localhost:8080",
    },
    APIKey: "mock-key",
})
```

### Python

```python
from google import genai
from google.genai.types import HttpOptions

client = genai.Client(
    api_key="mock-key",
    http_options=HttpOptions(base_url="http://localhost:8080"),
)
```

### TypeScript

```typescript
import { GoogleGenAI } from "@google/genai";

const ai = new GoogleGenAI({
    apiKey: "mock-key",
    baseURL: "http://localhost:8080",
});
```

### Rust (rig)

```rust
use rig::providers::gemini;

let client = gemini::Client::builder()
    .api_key("mock-key")
    .base_url("http://localhost:8080")
    .build()
    .expect("mock client");

let agent = client
    .agent("gemini-2.0-flash")
    .preamble("You are a helpful assistant.")
    .build();

let response = agent.prompt("Hello!").await?;
```

### Environment-based switching (Rust)

```rust
fn gemini_client() -> gemini::Client {
    match std::env::var("GEMINI_MOCK_URL").ok() {
        Some(url) => gemini::Client::builder()
            .api_key(std::env::var("GEMINI_API_KEY").unwrap_or("mock".into()))
            .base_url(url)
            .build()
            .expect("mock client"),
        None => gemini::Client::from_env(),
    }
}
```

```bash
GEMINI_MOCK_URL=http://localhost:8080 cargo run   # development
GEMINI_API_KEY=your-real-key cargo run             # production
```

---

## Keybindings

These apply to both the TUI and Web UI.

| Key | Action |
|-----|--------|
| `Tab` | Cycle focus between panels |
| `j` / `k` or `↑` / `↓` | Navigate list, scroll detail view |
| `Enter` | Open response composer for selected request |
| `F1` | Insert text response template |
| `F2` | Insert function call template |
| `F3` | Insert `429 Too Many Requests` error |
| `F4` | Insert `500 Internal Server Error` |
| `Ctrl+S` | Send composed response |
| `Ctrl+F` | Save response as fixture for auto-replay |
| `Esc` | Exit composer / cancel |
| `q` | Quit (confirms if requests are pending) |

---

## Architecture

```
┌─────────────┐         ┌──────────────────────────────────────────┐
│  Your App   │  HTTP   │              Cheaptrick                  │
│  (any SDK)  │────────▶│                                          │
│             │         │  ┌────────┐  observer     ┌───────────┐  │
│             │◀────────│  │ Server │──────────────▶│ TUI / Web │  │
│             │         │  │  :8080 │◀──────────────│   UI      │  │
└─────────────┘         │  └────────┘  response ch  └───────────┘  │
                        │       │                                  │
                        │       ▼                                  │
                        │  ┌────────┐                              │
                        │  │Fixtures│  (auto-reply if hash match)  │
                        │  └────────┘                              │
                        └──────────────────────────────────────────┘
```

The mock server and the UI (TUI or web) share a `RequestStore` through
a `RequestObserver` interface. Requests check the fixture store first;
on a miss, they block until a response is composed through the UI. The
web UI communicates over WebSocket for real-time updates.

---

## Roadmap

**In progress:**

- `streamGenerateContent` support with SSE chunking for streaming clients
- Fixture fuzzy matching — ignore volatile fields (timestamps, UUIDs,
  request IDs) when computing fixture hashes

**Planned:**

- Proxy recording mode — forward requests to the real Gemini API while
  saving request/response pairs as fixtures, enabling migration from
  live development to fully mocked development
- Response schema validation — validate composed responses against the
  Gemini `GenerateContentResponse` schema before sending, catching
  structural errors before they reach your client code
- Latency simulation — configurable response delays per fixture to test
  timeout handling and loading states in client code

---

## Contributing

Contributions are welcome. Bug fixes, fixture templates, documentation
improvements, and feature implementations are all appreciated.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-change`)
3. Commit your changes (`git commit -m 'Add my change'`)
4. Push to the branch (`git push origin feature/my-change`)
5. Open a Pull Request

Open an issue first for large changes so the approach can be discussed.

---

## License

MIT. See [LICENSE.md](LICENSE.md).