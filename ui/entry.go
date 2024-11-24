package ui

import (
	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/1Mochiyuki/gosky/ui/login"
	tea "github.com/charmbracelet/bubbletea"
)

type AppEntry struct {
	state *int
}

// TODO:
// need to do db call here for account picker and remove the logic from login_screen.go

var l = logger.Get()

func (a AppEntry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l.Debug().Int("state", *a.state).Msg("app entry")
	switch msg := msg.(type) {
	case login.AnyUserExistsMsg:
		credsLen := len(msg.Results)
		if credsLen == 1 {
			return login.InitLoginScreenModel(msg.Results[0].Handle, a.state), nil
		} else if credsLen > 1 {
			return login.NewMultiAccountLogin(msg.Results, a.state), nil
		} else {
			return login.InitLoginScreenModel("", a.state), nil
		}
	}

	return a, login.AnyCredentialsExist
}

func (a AppEntry) View() string {
	return " "
}

func (a AppEntry) Init() tea.Cmd {
	return nil
}

func NewAppEntry() AppEntry {
	init := 0
	return AppEntry{
		state: &init,
	}
}
