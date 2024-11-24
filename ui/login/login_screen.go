package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/1Mochiyuki/gosky/client"
	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/1Mochiyuki/gosky/db"
	"github.com/1Mochiyuki/gosky/errs"
	"github.com/1Mochiyuki/gosky/ui/send"
	"github.com/1Mochiyuki/gosky/ui/state"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginScreenModel struct {
	error           error
	state           *int
	LoginComponents []textinput.Model
	focusState      int
	rememberState   bool
}

func InitLoginScreenModel(handle string, state *int) LoginScreenModel {
	m := LoginScreenModel{
		LoginComponents: make([]textinput.Model, 2), focusState: 0, state: state,
	}
	var t textinput.Model
	for i := range m.LoginComponents {
		t = textinput.New()
		t.Cursor.Style = focusedStyle
		t.Cursor.SetMode(cursor.CursorHide)

		switch i {
		case 0:
			t.Placeholder = "BlueSky Handle"
			t.SetValue(handle)
			t.Focus()
			t.PlaceholderStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "BlueSky App Pass"
			t.EchoMode = textinput.EchoPassword
			t.CharLimit = 19
			t.EchoCharacter = '•'
			t.Prompt = ""
			t.Validate = func(s string) error {
				if len(s) < 19 {
					return errors.New("app pass not long enough")
				}

				return nil
			}
		}

		m.LoginComponents[i] = t
	}

	return m
}

func (l *LoginScreenModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(l.LoginComponents))

	for i := range l.LoginComponents {
		l.LoginComponents[i], cmds[i] = l.LoginComponents[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

var (
	noStyle       = lipgloss.NewStyle()
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#3c78fc"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#2953af"))
	blurredButton = noStyle.Render("Submit")
	focusedButton = focusedStyle.Render("Submit")

	blurredRememberMe = noStyle.Render(fmt.Sprintf("%s Remember Me", emptyCheckbox))
	focusedRememberMe = focusedStyle.Render(fmt.Sprintf("%s Remember Me", emptyCheckbox))
)

const (
	handlePos = iota
	passPos
	rememberMePos
	submitPos
	filledCheckBox = "☑"
	emptyCheckbox  = "☐"
)

var (
	log               = logger.Get()
	checkedCreds bool = false
)

func (l LoginScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case AnyUserExistsMsg:
		cmd = l.updateInputs(msg)
		if len(msg.Results) <= 0 {
			log.Debug().Msg("no credentials here")
			return l, cmd
		}
		if len(msg.Results) == 1 {
			creds := msg.Results[0]

			agent := client.NewAgent(context.Background(), "", creds.Handle, creds.AppPass)
			if err := agent.ConnectSave(); err != nil {
				var credErr *errs.CredentialsError
				if errors.As(err, &credErr) {
					l.error = err
					log.Error().Err(err).Msg("Incorrect Credentials")
					return l, cmd
				}
				l.error = errors.New("unhandled error, check logs file")
				log.Error().Err(err).Msg("Unknown Error")
				return l, cmd
			}
			return send.Model(agent), cmd
		}
		return NewMultiAccountLogin(msg.Results, l.state), cmd
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlLeft.String():
			*l.state = state.ACCOUNT_PICKER
			log.Debug().Int("state", *l.state).Msg("should go to account picker")
			return NewMultiAccountLogin(nil, l.state), AnyCredentialsExist
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			db.DB.Close()
			return l, tea.Quit
		case tea.KeyUp.String(), tea.KeyDown.String(), tea.KeyEnter.String(), tea.KeyCtrlJ.String(), tea.KeyCtrlK.String():
			s := msg.String()
			if s == tea.KeyEnter.String() {
				handle := l.LoginComponents[handlePos].Value()
				appPass := l.LoginComponents[passPos].Value()
				if l.focusState == rememberMePos {
					l.rememberState = !l.rememberState
				}
				if l.focusState == submitPos {

					agent := client.NewAgent(context.Background(), "", handle, appPass)
					var connErr error
					if l.rememberState {
						connErr = agent.ConnectSave()
					} else {
						connErr = agent.ConnectNoSave()
					}
					if connErr != nil {

						var credErr *errs.CredentialsError
						if errors.As(connErr, &credErr) {
							l.error = connErr
							log.Error().Err(connErr).Msg("Incorrect Credentials")
							return l, cmd
						}
						l.error = errors.New("unhandled error, check logs file")
						log.Error().Err(connErr).Msg("Unknown Error")
						return l, cmd
					}
					return send.Model(agent), cmd

				}
			}

			if s == tea.KeyUp.String() || s == tea.KeyCtrlK.String() {
				l.focusState--
			} else if s == tea.KeyDown.String() || s == tea.KeyCtrlJ.String() {
				l.focusState++
			}

			if l.focusState > len(l.LoginComponents)+1 {
				l.focusState = handlePos
			} else if l.focusState == submitPos || l.focusState < handlePos {
				l.focusState = submitPos
			} else if l.focusState == rememberMePos {
				l.focusState = rememberMePos
			}
			cmds := make([]tea.Cmd, len(l.LoginComponents))
			for i := 0; i <= len(l.LoginComponents)-1; i++ {
				if i == l.focusState {
					cmds[i] = l.LoginComponents[i].Focus()
					l.LoginComponents[i].PlaceholderStyle = focusedStyle
					l.LoginComponents[i].TextStyle = focusedStyle
					l.LoginComponents[i].Prompt = "> "
					continue
				}

				l.LoginComponents[i].Blur()
				l.LoginComponents[i].PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
				l.LoginComponents[i].Prompt = ""
				l.LoginComponents[i].TextStyle = noStyle
			}
			return l, tea.Batch(cmds...)
		}
	}

	cmd = l.updateInputs(msg)
	return l, cmd
}

func fmtNilErrAsEmpty(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func (l LoginScreenModel) View() string {
	button := &blurredButton
	if l.focusState == submitPos {
		button = &focusedButton
	}
	btn := fmt.Sprintf("%v", *button)
	rememberMeStr := blurredRememberMe
	if l.focusState == rememberMePos {
		switch l.rememberState {
		case true:
			rememberMeStr = focusedStyle.Render(fmt.Sprintf("%s Remember Me", filledCheckBox))
		case false:
			rememberMeStr = focusedRememberMe
		}
	} else {
		switch l.rememberState {
		case true:
			rememberMeStr = noStyle.Render(fmt.Sprintf("%s Remember Me", filledCheckBox))
		case false:
			rememberMeStr = blurredRememberMe
		}
	}
	rememberMeBox := fmt.Sprintf("\n%v\n", rememberMeStr)
	view := noStyle.PaddingLeft(50).Render(
		noStyle.Width(55).Border(lipgloss.RoundedBorder(), true, true, true, true).Align(lipgloss.Center).Render(
			"  [ GoSky Scheduler ]",
			fmt.Sprintf("\n%s\n", fmtNilErrAsEmpty(l.error)),

			noStyle.Margin(1, 0).Width(35).Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true, true, true, true).Align(lipgloss.Center).Render(

				l.LoginComponents[handlePos].View(),
			),
			noStyle.Margin(1, 0).Width(35).Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true, true, true, true).Align(lipgloss.Center).Render(
				l.LoginComponents[passPos].View(),
			),
			rememberMeBox,

			noStyle.Margin(1, 0).Width(13).Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true, true, true, true).Align(lipgloss.Center).Render(
				btn,
			),
		),
	)
	return view
}

func (l LoginScreenModel) Init() tea.Cmd {
	return textinput.Blink
}
