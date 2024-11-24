package login

import (
	"github.com/1Mochiyuki/gosky/db"
	"github.com/1Mochiyuki/gosky/ui/state"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Account struct {
	title string
}

func (a Account) Title() string {
	return a.title
}

func (a Account) Description() string {
	return ""
}

func (a Account) FilterValue() string {
	return a.title
}

type MultiAccountLogin struct {
	AccountList list.Model
	state       *int
}

func (m MultiAccountLogin) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():

			db.DB.Close()
			return m, tea.Quit
		case tea.KeyEnter.String():
			handle := m.AccountList.SelectedItem().(Account).title
			*m.state = state.MAIN_ACCOUNT_LOGIN
			l.Debug().Int("state", *m.state).Msg("should go to login screen")

			return InitLoginScreenModel(handle, m.state), nil
		}
	}
	var cmd tea.Cmd
	m.AccountList, cmd = m.AccountList.Update(msg)
	return m, cmd
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m MultiAccountLogin) View() string {
	return docStyle.Render(m.AccountList.View())
}

func (m MultiAccountLogin) Init() tea.Cmd {
	return AnyCredentialsExist
}

func NewMultiAccountLogin(logins []Credentials, state *int) MultiAccountLogin {
	creds := []list.Item{}
	for handle := range LOGIN_CACHCE {
		creds = append(creds, Account{handle})
	}
	credList := list.New(creds, list.NewDefaultDelegate(), 25, 15)
	credList.Title = "Stored Accounts"
	return MultiAccountLogin{
		AccountList: credList,
		state:       state,
	}
}
