<div align="center">

# 🎭 Cheaptrick

**A human-in-the-loop mock server for the Google Gemini API**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](CONTRIBUTING.md)

Cheaptrick intercepts Google Gemini API requests, lets you craft responses
by hand, and replays them from fixtures — so you can develop and debug
LLM-powered agents locally without spending tokens.

[Getting Started](#getting-started) · [Usage](#usage) · [Fixtures](#fixtures) · [Tool-Call Debugging](#tool-call-debugging-with-the-shell) · [SDK Examples](#connecting-your-app) · [Keybindings](#keybindings) · [Contributing](#contributing)

</div>

---

## The Problem

Developing an LLM agent means hundreds of API round-trips: testing parsing
logic, iterating on tool-call schemas, handling edge cases in multi-turn
flows. Each call costs tokens and hits rate limits, but the responses you
need during development are often predictable.

Cheaptrick replaces the Gemini API with a local endpoint where you control
every response. Requests appear in a web dashboard. You compose a JSON
response (or pick a template), send it, and your client receives it as if
it came from Gemini. Save a response as a fixture, and identical future
requests are replayed automatically.

The result: deterministic, reproducible agent development at zero API cost.

---

## Features

- **Web UI** — A React-based web dashboard embedded in the binary (no
  Node.js required). Request inspection, response composing, fixture
  management, and real-time WebSocket updates — all in the browser.
- **Response templates** — Pre-built skeletons for text responses, function
  calls, 429 rate-limit errors, and 500 server errors.
- **Fixture replay** — Save any response as a fixture keyed by request
  hash. Matching requests are auto-replied without human intervention.
- **Interactive shell** — A REPL backed by the official `google.golang.org/genai`
  client that sends prompts to your mock server and handles tool-call
  loops with canned response files.
- **Sample tool responses** — Generate 20 pre-built canned tool response
  sets covering weather, search, email, SQL, file I/O, translation, and
  more with a single command.
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

### Quick start

```bash
# Generate sample tool responses
cheaptrick tools -o ./mock_tools

# Start the mock server and web UI
cheaptrick web

# In another terminal, start the interactive shell
cheaptrick shell --tool-responses ./mock_tools
```

---

## Usage

Cheaptrick provides four subcommands: `web`, `shell`, `tools`, and
`fixtures`.

### `web` — Mock server with browser UI

```bash
cheaptrick web
cheaptrick web --port 9090 --web-port 4000 --open
cheaptrick web --tls-cert cert.pem --tls-key key.pem
```

Starts the mock Gemini API on `:8080` and serves the React frontend on
`:3000`. The frontend is embedded in the binary — no separate build step
or Node.js installation is needed.

Send a request from another terminal:

```bash
curl -s -X POST \
  http://localhost:8080/v1beta/models/gemini-2.0-flash:generateContent \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"Hello Gemini!"}]}]}'
```

The request blocks. In the web UI, a pending request appears. Select it,
compose a response (or click a template button), and send. The curl
command returns with your response.

Click **Save Fixture** to save the response for auto-replay. Subsequent
identical requests are answered instantly without human input.

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | Gemini mock server port |
| `--web-port` | `3000` | Web UI server port |
| `--fixtures` | `./fixtures` | Fixture directory path |
| `--log` | `mock_log.jsonl` | JSONL log file path |
| `--tls-cert` | | TLS certificate file |
| `--tls-key` | | TLS key file |
| `--open` | `true` | Auto-open browser on startup |

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

### `tools` — Generate sample canned tool responses

```bash
cheaptrick tools
cheaptrick tools -o ./my_tools
```

Generates 20 pre-built canned tool response sets with 62 total files,
covering a range of common tool-calling patterns. The generated files
are ready to use with `cheaptrick shell --tool-responses`. See
[Sample Tool Responses](#sample-tool-responses) for the full list.

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output-dir` | `./mock_tools` | Output directory for generated files |

### `fixtures` — Generate starter fixtures

```bash
cheaptrick fixtures
cheaptrick fixtures --output-dir ./test_assets
```

Generates 30 predefined fixture files for common text and tool-call
prompts, plus a `MANIFEST.md` index. These are Gemini API response
fixtures used for auto-replay by the mock server — distinct from the
tool response files generated by `cheaptrick tools`.

---

## Fixtures

Cheaptrick's core feature is the ability to save and replay Gemini API responses as fixtures. This enables deterministic, reproducible agent development at zero API cost.

### How Fixtures Work

A fixture is a JSON file containing a `GenerateContentResponse` payload. When the mock server receives an incoming request, it calculates a **request hash** and checks the fixtures directory for a matching file.

Matching follows these steps:
1. **Hashing**: The server identifies the "canonical" request content:
   - If the request has a `contents` array with at least one part containing `text`, that text is used.
   - Otherwise, the entire raw JSON request body is used.
   - The SHA256 hash of this content becomes the fixture key.
2. **Lookup**: The server looks for `<hash>.json` in the directory specified by `--fixtures` (defaulting to `./fixtures`).
3. **Auto-Replay**: If found, the server immediately returns the fixture's content as the response, without human-in-the-loop intervention.

### Creating and Managing Fixtures

#### Saving from the Web UI
The easiest way to create fixtures is while you are debugging your application:
- Point your SDK at the mock server.
- When a request appears in the **Web UI**, compose a response (or pick a template).
- Press **Ctrl+F** (or click "Save Fixture") to save the response and Press **Ctrl+S** to send it.
- Future identical requests will now be auto-replied.

#### Generating Starter Fixtures
Run the `fixtures` command to generate a set of common text and tool-call prompts:
```bash
cheaptrick fixtures --output-dir ./fixtures
```
This generates 30 predefined fixture files and a `MANIFEST.md` that acts as a reference for which prompts map to which hashes.

#### Manual Management
Fixtures are just plain JSON files. You can:
- **Edit** them manually to tweak the model's tone, schema, or error messages.
- **Copy** them between projects or share them with team members.
- **Version control** them alongside your code to ensure your agent tests always have the necessary mocks.

### Fixture Example

If you send a request for "What is the capital of France?", the server computes the hash `115049a298532be2f181edb03f766770c0db84c22aff39003fec340deaec7545`. If `./fixtures/115049a298532be2f181edb03f766770c0db84c22aff39003fec340deaec7545.json` exists with the following content:

```json
{
  "candidates": [
    {
      "content": {
        "role": "model",
        "parts": [{"text": "The capital of France is Paris."}]
      },
      "finishReason": "STOP"
    }
  ]
}
```

The mock server will reply with "The capital of France is Paris." instantly every time that prompt is received.

---

## Sample Tool Responses

The `cheaptrick tools` command generates canned tool response files that
the shell uses to automatically reply to `FunctionCall` parts. Each tool
demonstrates a different resolution strategy (static, argument-matched,
or sequenced), giving developers a working reference for every pattern.

```bash
cheaptrick tools -o ./mock_tools
cheaptrick shell --tool-responses ./mock_tools
```

### Generated tools

| # | Function | Type | Description |
|---|----------|------|-------------|
| 1 | `get_weather` | Argument-matched (city) | Weather data for Paris, Tokyo, New York, London, São Paulo with forecasts |
| 2 | `web_search` | Sequenced (3 pages) | Paginated search results; third page returns empty |
| 3 | `send_email` | Static | Confirms email sent with message ID and timestamp |
| 4 | `execute_sql` | Argument-matched (query_type) | SELECT returns sample rows; INSERT/UPDATE/DELETE return affected counts |
| 5 | `read_file` | Static | Returns file content, metadata, and MIME type |
| 6 | `write_file` | Static | Confirms bytes written with path and timestamp |
| 7 | `create_calendar_event` | Static | Confirms event creation with calendar link |
| 8 | `get_stock_price` | Argument-matched (symbol) | Quotes for AAPL, GOOGL, MSFT with market data |
| 9 | `translate_text` | Argument-matched (target_language) | Translations to Spanish, French, German, Japanese |
| 10 | `get_directions` | Static | Turn-by-turn driving directions with distance and duration |
| 11 | `create_ticket` | Static | Jira-style ticket creation with ID and URL |
| 12 | `send_slack_message` | Static | Confirms message posted with permalink |
| 13 | `get_exchange_rate` | Argument-matched (to) | USD to EUR, GBP, JPY exchange rates |
| 14 | `dns_lookup` | Static | A, AAAA, MX, NS, TXT records |
| 15 | `scrape_url` | Static | Page title, text content, links, and response headers |
| 16 | `get_github_repo` | Argument-matched (repo) | Repository metadata for bubbletea and rig |
| 17 | `get_user_profile` | Argument-matched (user_id) | User profiles with preferences and login history |
| 18 | `http_request` | Argument-matched (method) | GET/POST/PUT/DELETE with appropriate status codes |
| 19 | `get_news` | Sequenced (3 pages) | Paginated news articles; third page returns empty |
| 20 | `run_shell_command` | Argument-matched (command) | Output for ls, pwd, whoami |

### Resolution strategies

The three types demonstrated by the generated tools:

**Static** — A single `<function>.json` file. Every call to that function
returns the same response, regardless of arguments. Files can include
`{{args.fieldname}}` placeholders for dynamic substitution. Good for
tools with predictable output like `send_email` or `write_file`.

**Argument-matched** — A `<function>/` directory containing a `_match.json`
routing file that selects a response based on a top-level argument value.
Includes a `_default.json` fallback for unrecognized values. Good for
tools where the response depends on a key parameter like `get_weather`
(city) or `execute_sql` (query_type).

**Sequenced** — Numbered files `<function>.N.json` (1-indexed) that return
different responses on successive calls to the same function. When the
call count exceeds the highest numbered file, the response clamps to the
last file. Good for pagination, polling, and retry patterns. Call
counters reset on `/clear`.

### Generated directory structure

```
mock_tools/
  MANIFEST.md                    # Documents every tool and its type
  send_email.json                # Static tools
  read_file.json
  write_file.json
  create_calendar_event.json
  get_directions.json
  create_ticket.json
  send_slack_message.json
  dns_lookup.json
  scrape_url.json
  web_search.1.json              # Sequenced tools
  web_search.2.json
  web_search.3.json
  get_news.1.json
  get_news.2.json
  get_news.3.json
  get_weather/                   # Argument-matched tools
    _match.json
    paris.json
    tokyo.json
    new_york.json
    london.json
    sao_paulo.json
    _default.json
  execute_sql/
    _match.json
    select.json
    insert.json
    update.json
    delete.json
    _default.json
  get_stock_price/
    _match.json
    aapl.json
    googl.json
    msft.json
    _default.json
  translate_text/
    _match.json
    spanish.json
    french.json
    german.json
    japanese.json
    _default.json
  get_exchange_rate/
    _match.json
    eur.json
    gbp.json
    jpy.json
    _default.json
  get_github_repo/
    _match.json
    bubbletea.json
    rig.json
    _default.json
  get_user_profile/
    _match.json
    alice.json
    bob.json
    _default.json
  http_request/
    _match.json
    get.json
    post.json
    put.json
    delete.json
  run_shell_command/
    _match.json
    ls.json
    pwd.json
    whoami.json
    _default.json
```

### Adding your own tools

Create a new JSON file in the tool responses directory following the
naming conventions above. The shell picks up files at startup. Use
`/clear` between conversations to reset sequenced call counters.

---

## Tool-Call Debugging with the Shell

The shell is designed for debugging multi-turn tool-calling flows. In this
setup, the developer operates both sides of the conversation:

- **The shell** acts as the tool executor — it resolves canned responses
  from disk and sends `FunctionResponse` parts automatically.
- **The developer** (via the web UI) acts as the model — crafting
  `FunctionCall` or text responses for each turn.

This two-window workflow lets you simulate complete agentic loops, control
every decision boundary, and build fixture libraries that capture
multi-step flows.

### How it works

```
 Terminal (shell)                      Browser (web UI)
 ────────────────                      ────────────────
 Prompt> What's the weather?
   → GenerateContent sent to :8080
   → blocks, waiting...
                                       [PENDING] req-01 appears
                                       Click Function Call template,
                                       compose:
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
                                       Click Text Response template,
                                       type:
                                         "It's 18°C and cloudy."
                                       Ctrl+S to send

 Receives text response:
   "It's 18°C and cloudy."
 Prints to terminal. Chain complete.

 ─── Chain Complete (3 turns, 12.4s) ──────────
  1  user  │ What's the weather?
  2  mock  │ ƒ get_weather(city="Paris")     [via web]
  3  tool  │ → {"temp":18,"condition":"cloudy"} [canned: paris.json]
  4  mock  │ It's 18°C and cloudy.           [via web]
 ─────────────────────────────────────────────

 Prompt> _
```

### Setting up canned tool responses

Create a directory for your tool response files (or generate one with
`cheaptrick tools`) and pass it to the shell with `--tool-responses`:

```bash
cheaptrick tools -o ./mock_tools
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
the web UI, the shell matches on the `city` field and loads `paris.json`.

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

These exported sequences provide a complete record of a multi-turn conversation that can be used for manual inspection or as a base for creating individual fixtures.

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

These apply to the Web UI when the browser window is focused.

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

---

## Architecture

```
┌─────────────┐         ┌──────────────────────────────────────────┐
│  Your App   │  HTTP   │              Cheaptrick                  │
│  (any SDK)  │────────▶│                                          │
│             │         │  ┌────────┐  observer     ┌───────────┐  │
│             │◀────────│  │ Mock   │──────────────▶│  Web UI   │  │
│             │         │  │ Server │◀──────────────│  (:3000)  │  │
└─────────────┘         │  │ (:8080)│  response ch  └───────────┘  │
                        │  └────────┘       ▲ WebSocket             │
                        │       │           │                      │
                        │       ▼           │                      │
                        │  ┌────────┐  ┌────┴────┐                 │
                        │  │Fixtures│  │ Browser │                 │
                        │  └────────┘  └─────────┘                 │
                        └──────────────────────────────────────────┘
```

The `cheaptrick web` command starts two HTTP servers in goroutines: the
mock Gemini API (`:8080`) and the Gin-based web server (`:3000`) which
serves the embedded React SPA and exposes a REST API under `/api/`.
Both servers share a `RequestStore` through a `RequestObserver`
interface. Requests check the fixture store first; on a miss, they block
until a response is composed through the web UI. The browser receives
real-time updates over a WebSocket connection at `/ws`.

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