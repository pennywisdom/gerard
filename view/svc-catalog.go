package view

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func NewModel() (*model, error) {
	// create a new spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	return &model{
		spinner: s,
	}, nil
}

type model struct {
	quitting bool
	spinner  spinner.Model
	result   string
	err      error
	F        Fn
	width    int
	height   int
}

type Fn func() error

var _ tea.Model = (*model)(nil)

func (m model) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		err := m.F()
		if err != nil {
			return err
		}
		return "done"
	}, m.spinner.Tick)
}

func (m model) View() string {
	if m.err != nil {
		return wordwrap.String(m.err.Error(), m.width) + "\n Press Ctrl+C to quit"
	}
	if m.quitting {
		return fmt.Sprintf("%s\nQuitting...", m.result)
	} else {
		return m.spinner.View()
	}
}

// Update is called with a tea.Msg, representing something that happened within
// our application.
//
// This can be things like terminal resizing, keypresses, or custom IO.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Let's figure out what is in tea.Msg, and what we need to do.
	// log.Printf("msg: %v", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// The terminal was resized.  We can access the new size with:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		// msg is a keypress.  We can handle each key combo uniquely, and update
		// our state:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlBackslash:
			// In this case, ctrl+c or ctrl+backslash quits the app by sending a
			// tea.Quit cmd.  This is a Bubbletea builtin which terminates the
			// overall framework which renders our model.
			//
			// Unfortunately, if you don't include this quitting can be, uh,
			// frustrating, as bubbletea catches every key combo by default.
			return m, tea.Quit

		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case error:
		m.err = msg
		return m, nil
	case string:
		var cmd tea.Cmd
		m.quitting = true
		m.result = msg
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	// We return an updated model to Bubbletea for rendering here.  This allows
	// us to mutate state so that Bubbletea can render an updated view.
	//
	// We also return "commands".  A command is something that you need to do
	// after rendering.  Each command produces a tea.Msg which is its *result*.
	// Bubbletea calls this Update function again with the tea.Msg - this is our
	// render loop.
	//
	// For now, we have no commands to run given the message is not a keyboard
	// quit combo.
	return m, nil
}
