package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"cheaptrick/internal/fixture"
	"cheaptrick/internal/store"
)

const (
	stateList = iota
	stateDetail
	stateComposer
	stateConfirmQuit
)

var (
	docStyle       = lipgloss.NewStyle().Margin(1, 2)
	activeBorder   = lipgloss.Color("62")
	inactiveBorder = lipgloss.Color("240")
	detailStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	composerStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
)

type tuiRequest struct {
	*store.Request
	Answered bool
}

func (t tuiRequest) Title() string {
	status := "[PENDING]"
	if t.Answered {
		status = "[ANSWERED]"
	}
	return fmt.Sprintf("%s: %s", status, t.ID)
}

func (t tuiRequest) Description() string {
	desc := ""
	if contents, ok := t.ParsedBody["contents"].([]interface{}); ok && len(contents) > 0 {
		if contentMap, ok := contents[0].(map[string]interface{}); ok {
			if parts, ok := contentMap["parts"].([]interface{}); ok && len(parts) > 0 {
				if partMap, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := partMap["text"].(string); ok {
						desc = strings.ReplaceAll(text, "\n", " ")
						if len(desc) > 80 {
							desc = desc[:77] + "..."
						}
					}
				}
			}
		}
	}
	return t.Timestamp.Format("15:04:05") + " | " + desc
}

func (t tuiRequest) FilterValue() string { return t.ID + " " + t.Model }

type model struct {
	reqCh       <-chan *store.Request
	eventCh     <-chan string
	respondedCh <-chan string
	reqStore    *store.Store
	fixturesDir string

	requests     []tuiRequest
	list         list.Model
	viewport     viewport.Model
	textarea     textarea.Model
	state        int
	notification string
}

type eventMsg string

func InitialModel(reqStore *store.Store, reqCh <-chan *store.Request, eventCh <-chan string, respondedCh <-chan string, fixturesDir string) model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Gemini Mock Requests"
	l.SetShowStatusBar(false)

	vp := viewport.New(0, 0)
	vp.SetContent("Select a request to view details.")

	ta := textarea.New()
	ta.Placeholder = "Write JSON response here..."
	ta.ShowLineNumbers = true

	return model{
		reqStore:    reqStore,
		reqCh:       reqCh,
		eventCh:     eventCh,
		respondedCh: respondedCh,
		fixturesDir: fixturesDir,
		list:        l,
		viewport:    vp,
		textarea:    ta,
		state:       stateList,
	}
}

func waitForRequest(sub <-chan *store.Request) tea.Cmd {
	return func() tea.Msg { return <-sub }
}

func waitForEvent(sub <-chan string) tea.Cmd {
	return func() tea.Msg { return eventMsg(<-sub) }
}

type respondedMsg string

func waitForResponded(sub <-chan string) tea.Cmd {
	return func() tea.Msg { return respondedMsg(<-sub) }
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForRequest(m.reqCh),
		waitForEvent(m.eventCh),
		waitForResponded(m.respondedCh),
		textarea.Blink,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width/3-h, msg.Height-v-3) // 1/3 width

		vpWidth := msg.Width*2/3 - h - 4
		vpHeight := msg.Height - v - 3

		m.viewport.Width = vpWidth
		m.viewport.Height = vpHeight

		m.textarea.SetWidth(vpWidth)
		m.textarea.SetHeight(vpHeight)
		return m, nil

	case eventMsg:
		m.notification = string(msg)
		cmds = append(cmds, waitForEvent(m.eventCh))

	case respondedMsg:
		id := string(msg)
		for i, r := range m.requests {
			if r.ID == id {
				m.requests[i].Answered = true
			}
		}
		var items []list.Item
		for _, r := range m.requests {
			items = append(items, r)
		}
		m.list.SetItems(items)
		cmds = append(cmds, waitForResponded(m.respondedCh))

	case *store.Request:
		req := tuiRequest{Request: msg, Answered: false}
		m.requests = append(m.requests, req)

		var items []list.Item
		for _, r := range m.requests {
			items = append(items, r)
		}
		m.list.SetItems(items)
		cmds = append(cmds, waitForRequest(m.reqCh))
	}

	if m.state == stateConfirmQuit {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "y" || msg.String() == "Y" {
				return m, tea.Quit
			}
			m.state = stateList
		}
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" && m.state == stateList {
			for _, r := range m.requests {
				if !r.Answered {
					m.state = stateConfirmQuit
					return m, tea.Batch(cmds...)
				}
			}
			return m, tea.Quit
		}

		if msg.String() == "tab" {
			if m.state == stateList {
				m.state = stateDetail
			} else if m.state == stateDetail {
				m.state = stateList
			}
			return m, tea.Batch(cmds...)
		}

		switch m.state {
		case stateList:
			if msg.Type == tea.KeyEnter && len(m.requests) > 0 {
				m.state = stateComposer
				m.textarea.Focus()
				if len(strings.TrimSpace(m.textarea.Value())) == 0 {
					m.textarea.SetValue(fixture.TemplateText())
				}
			}
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)

			if len(m.requests) > 0 {
				idx := m.list.Index()
				if idx >= 0 && idx < len(m.requests) {
					req := m.requests[idx]
					b, _ := json.MarshalIndent(req.ParsedBody, "", "  ")
					m.viewport.SetContent(string(b))
				}
			}

		case stateDetail:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)

		case stateComposer:
			switch msg.String() {
			case "esc":
				m.state = stateList
				m.textarea.Blur()
			case "f1":
				m.textarea.SetValue(fixture.TemplateText())
			case "f2":
				m.textarea.SetValue(fixture.TemplateFunctionCall(m.requests[m.list.Index()].ParsedBody))
			case "f3":
				m.textarea.SetValue(fixture.Template429())
			case "f4":
				m.textarea.SetValue(fixture.Template500())
			case "ctrl+s":
				idx := m.list.Index()
				req := m.requests[idx]
				if !req.Answered {
					m.reqStore.MarkResponded(req.ID, "manual")
					req.ResponseCh <- m.textarea.Value()
					m.requests[idx].Answered = true

					var items []list.Item
					for _, r := range m.requests {
						items = append(items, r)
					}
					m.list.SetItems(items)
					m.textarea.SetValue("")
				}
				m.state = stateList
				m.textarea.Blur()
			case "ctrl+f":
				idx := m.list.Index()
				req := m.requests[idx]
				if m.fixturesDir != "" {
					err := fixture.SaveFixture(m.fixturesDir, req.Hash, m.textarea.Value())
					if err == nil {
						m.notification = "Saved fixture " + req.Hash[:8]
					} else {
						m.notification = "Failed to save: " + err.Error()
					}
				} else {
					m.notification = "--fixtures dir not set"
				}
			default:
				m.textarea, cmd = m.textarea.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.state == stateConfirmQuit {
		return "\n\n  Pending requests exist. Are you sure you want to quit? (y/N)\n\n"
	}

	leftFrame := docStyle.Render(m.list.View())

	var rightFrame string
	if m.state == stateComposer {
		composerStyle = composerStyle.BorderForeground(activeBorder)
		rightFrame = composerStyle.Width(m.textarea.Width() + 2).Height(m.textarea.Height() + 2).Render(m.textarea.View())
	} else {
		if m.state == stateDetail {
			detailStyle = detailStyle.BorderForeground(activeBorder)
		} else {
			detailStyle = detailStyle.BorderForeground(inactiveBorder)
		}
		rightFrame = detailStyle.Width(m.viewport.Width + 2).Height(m.viewport.Height + 2).Render(m.viewport.View())
	}

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftFrame, rightFrame)

	notiBar := lipgloss.NewStyle().
		Background(lipgloss.Color("235")).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1).
		Width(m.list.Width() + m.viewport.Width + 5 + 4).
		Render(m.notification)

	return lipgloss.JoinVertical(lipgloss.Left, mainView, notiBar)
}
