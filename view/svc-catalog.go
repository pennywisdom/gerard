package view

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewModel() (*model, error) {
	return &model{}, nil
}

type model struct {
	Event string
	F     Fn
}

type Fn func(cmd *cobra.Command, args []string) error

var _ tea.Model = (*model)(nil)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	//m.F()
	if m.Event != "" {
		return fmt.Sprintf("You've selected: %s", m.Event)
	}
	return "TODO" // We'll do this soon :)
}

// Update is called with a tea.Msg, representing something that happened within
// our application.
//
// This can be things like terminal resizing, keypresses, or custom IO.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Let's figure out what is in tea.Msg, and what we need to do.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// The terminal was resized.  We can access the new size with:
		_, _ = msg.Width, msg.Height
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
