package appcore

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/HexmosTech/git-lrc/internal/decisionflow"
)

type decisionPrompt struct {
	Title       string
	Description string
	AllowCommit bool
	AllowAbort  bool
	AllowSkip   bool
	AllowVouch  bool
}

type decisionTUIModel struct {
	prompt  decisionPrompt
	output  chan<- int
	status  string
	decided bool
}

func (m decisionTUIModel) Init() tea.Cmd {
	return nil
}

func (m decisionTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.decided {
		return m, tea.Quit
	}

	switch v := msg.(type) {
	case tea.KeyPressMsg:
		key := strings.ToLower(v.String())
		switch key {
		case "enter":
			if m.prompt.AllowCommit {
				m.submit(decisionflow.DecisionCommit)
				return m, tea.Quit
			}
		case "ctrl+c", "q", "esc":
			if m.prompt.AllowAbort {
				m.submit(decisionflow.DecisionAbort)
				return m, tea.Quit
			}
		case "ctrl+s", "s":
			if m.prompt.AllowSkip {
				m.submit(decisionflow.DecisionSkip)
				return m, tea.Quit
			}
		case "ctrl+v", "ctrl+y", "v", "y":
			if m.prompt.AllowVouch {
				m.submit(decisionflow.DecisionVouch)
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m decisionTUIModel) View() tea.View {
	var lines []string
	lines = append(lines, "")
	lines = append(lines, "LiveReview Decision")
	lines = append(lines, "-------------------")
	if strings.TrimSpace(m.prompt.Title) != "" {
		lines = append(lines, m.prompt.Title)
	}
	if strings.TrimSpace(m.prompt.Description) != "" {
		lines = append(lines, m.prompt.Description)
	}
	lines = append(lines, "")
	lines = append(lines, "Keys:")
	if m.prompt.AllowCommit {
		lines = append(lines, "  Enter      Continue with commit")
	}
	if m.prompt.AllowSkip {
		lines = append(lines, "  Ctrl-S/S   Skip review and continue")
	}
	if m.prompt.AllowVouch {
		lines = append(lines, "  Ctrl-V/V   Vouch and continue")
	}
	if m.prompt.AllowAbort {
		lines = append(lines, "  Ctrl-C/Q   Abort")
	}
	if m.status != "" {
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("Status: %s", m.status))
	}
	lines = append(lines, "")
	return tea.NewView(strings.Join(lines, "\n"))
}

func (m *decisionTUIModel) submit(code int) {
	if m.decided {
		return
	}
	m.decided = true
	switch code {
	case decisionflow.DecisionCommit:
		m.status = "commit selected"
	case decisionflow.DecisionSkip:
		m.status = "skip selected"
	case decisionflow.DecisionVouch:
		m.status = "vouch selected"
	case decisionflow.DecisionAbort:
		m.status = "abort selected"
	}
	select {
	case m.output <- code:
	default:
	}
}

func startTerminalDecisionBubbleTea(prompt decisionPrompt) (<-chan int, func(), <-chan struct{}) {
	decisionCh := make(chan int, 1)
	doneCh := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer close(doneCh)
		defer close(decisionCh)

		model := decisionTUIModel{prompt: prompt, output: decisionCh}
		program := tea.NewProgram(model, tea.WithContext(ctx))
		_, _ = program.Run()
	}()

	stop := func() {
		cancel()
	}

	return decisionCh, stop, doneCh
}
