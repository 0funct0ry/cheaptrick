<div align="center">

# 🎭 Cheaptrick

**A human-in-the-loop mock server for the Google Gemini API**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](CONTRIBUTING.md)
[![Built with Bubble Tea](https://img.shields.io/badge/Built_with-Bubble_Tea-ff69b4?style=flat-square)](https://github.com/charmbracelet/bubbletea)

Stop burning Gemini tokens while you're still writing `println!("here")`.

Cheaptrick sits between your app and the Gemini API, letting you intercept every request and hand-craft every response — all from a slick terminal UI.

[Getting Started](#-getting-started) · [Usage](#-usage) · [SDK Examples](#-connecting-your-app) · [Keybindings](#%EF%B8%8F-keybindings) · [Contributing](#-contributing)

</div>

---

## Why Cheaptrick?

Building an LLM-powered agent means hundreds of API round-trips during development. Most of them are throwaway calls while you're debugging parsing logic, testing tool-call flows, or iterating on prompts. That's money and rate-limit headroom you'll never get back.

Cheaptrick gives you a local Gemini-compatible endpoint where **you** decide what comes back. Pending requests appear in a terminal dashboard. You pick a template, tweak the JSON, and hit send. Your client code doesn't know the difference.

Once you're happy with a response, save it as a **fixture**. The next time the same request arrives, Cheaptrick replies instantly — no human in the loop needed. Over time, you build a deterministic replay layer that doubles as your test suite.

**In short:** develop fast, spend nothing, and ship with confidence.

---

## ✨ Highlights

| | Feature | Description |
|---|---|---|
| 🖥️ | **Terminal UI & Web UI** | Full Bubble Tea TUI *and* a modern React Web UI — request list, detail viewer, and response composer built right in |
| 📝 | **Response Templates** | One-key skeletons for text replies, function calls, rate-limit errors, and server errors |
| 💾 | **Fixture Replay** | Save a response once, auto-reply forever. Build a fixture library as you develop |
| 🐚 | **Interactive Shell** | Built-in REPL that talks to your mock server using the official `google.golang.org/genai` client |
| 🔒 | **TLS Support** | Spin up HTTPS when your client requires it |
| 📊 | **Request Logging** | Every request/response pair logged to JSONL for post-hoc debugging |
| 🔌 | **Drop-in Compatible** | Works with every official Gemini SDK — Go, Python, TypeScript, Rust — just swap the base URL |

---

## 🚀 Getting Started

### Prerequisites

- **Go 1.22** or later

### Install

```bash
# Clone and build
git clone https://github.com/yourusername/cheaptrick.git
cd cheaptrick
make build

# Binary lands in ./bin/cheaptrick
./bin/cheaptrick --help
```

Or install directly into your `$GOPATH/bin`:

```bash
make install
cheaptrick --help
```

---

## 📖 Usage

Cheaptrick has four main subcommands: **`start`**, **`web`**, **`shell`**, and **`fixtures`**.

### `start` — Launch the mock server + TUI

```bash
# Defaults to localhost:8080
cheaptrick start

# Custom port and fixture directory
cheaptrick start --port 9090 --fixtures ./my_fixtures

# With TLS
cheaptrick start --tls-cert cert.pem --tls-key key.pem
```

**The workflow in 60 seconds:**

Open a second terminal and fire a request at the mock server:

```bash
curl -s -X POST \
  http://localhost:8080/v1beta/models/gemini-2.0-flash:generateContent \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"Hello Gemini!"}]}]}'
```

The curl hangs. In the TUI, a **\[PENDING\]** request appears. Select it, press **Enter** to open the response composer, hit **F1** for a text-response template, edit the JSON to your liking, and press **Ctrl+S**. The curl returns your crafted response immediately.

Want to replay that same response automatically next time? Press **Ctrl+F** to save it as a fixture. Future identical requests get an instant **\[AUTO\]** reply.

### `web` — Launch the mock server + Web UI

Prefer the browser over the terminal? The `web` command starts the mock API and serves a beautiful React frontend embedded right in the binary. No Node.js required!

```bash
# Starts API on :8080 and Web UI on :3000
cheaptrick web

# Custom ports and auto-open browser
cheaptrick web --port 9090 --web-port 4000 --open
```

The Web UI comes with everything from the TUI: a request list, JSON syntax highlighting, response composing, templates (F1-F4), and fixture saving (Ctrl+S / Ctrl+F).

### `shell` — Interactive REPL for quick testing

Talk to your mock server directly without constructing curl payloads. The shell uses the official `google.golang.org/genai` client under the hood.

```bash
cheaptrick shell
cheaptrick shell --host 127.0.0.1 --port 9090 --model gemini-1.5-pro
```

| Flag | Default | Description |
|------|---------|-------------|
| `-H, --host` | `localhost` | Mock server host |
| `-p, --port` | `8080` | Mock server port |
| `-m, --model` | `gemini-2.0-flash` | Model name in requests |
| `--history-file` | OS temp dir | Readline history path |

Set `GEMINI_API_KEY` if your client requires one (Cheaptrick ignores it).

### `fixtures` — Generate starter fixture files

Bootstrap 30 predefined fixtures for common text and tool-call prompts, plus a `MANIFEST.md` index.

```bash
cheaptrick fixtures
cheaptrick fixtures --output-dir ./test_assets
```

---

## 🔌 Connecting Your App

Point any official Gemini SDK at `http://localhost:8080` and use any string as the API key. Cheaptrick doesn't validate credentials.

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

**Tip:** Use an environment variable to toggle between mock and production:

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
# Development
GEMINI_MOCK_URL=http://localhost:8080 cargo run

# Production
GEMINI_API_KEY=sk-real-key cargo run
```

---

## ⌨️ Keybindings

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
| `Ctrl+F` | Save response as auto-replay fixture |
| `Esc` | Exit composer / cancel |
| `q` | Quit (confirms if requests are pending) |

---

## 🏗️ How It Works

```
┌─────────────┐         ┌──────────────────────────────────────────┐
│  Your App   │  HTTP   │              Cheaptrick                  │
│  (any SDK)  │────────▶│                                          │
│             │         │  ┌────────┐  event hook   ┌───────────┐  │
│             │◀────────│  │ Server │──────────────▶│ TUI / Web │  │
│             │         │  │  :8080 │◀──────────────│   UI      │  │
└─────────────┘         │  └────────┘  response     └───────────┘  │
                        │       │                                  │
                        │       ▼                                  │
                        │  ┌────────┐                              │
                        │  │Fixtures│  (auto-reply if match found) │
                        │  └────────┘                              │
                        └──────────────────────────────────────────┘
```

The HTTP server runs in a goroutine. When a request arrives, it checks the fixture store first. On a cache miss, it notifies the UI observers (Bubble Tea TUI or React Web UI via WebSocket) and blocks until you compose a response. Everything is recorded via an internal shared Request Store.

---

## 🗺️ Roadmap

- [ ] Streaming response support (`streamGenerateContent` with SSE chunking)
- [ ] Fixture fuzzy matching (ignore timestamps, IDs)
- [ ] Request diffing (highlight what changed between similar requests)
- [ ] Multi-turn conversation threading in the TUI
- [ ] Export fixtures as Go/Python/Rust test helpers
- [ ] Plugin system for custom response logic

---

## 🤝 Contributing

Contributions are welcome and appreciated. Whether it's a bug fix, new feature, documentation improvement, or fixture template — all PRs are valued.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please open an issue first for large changes so we can discuss the approach.

---

## 📄 License

Distributed under the MIT License. See [LICENSE.md](LICENSE.md) for details.

