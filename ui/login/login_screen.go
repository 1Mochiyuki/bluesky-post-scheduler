package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/1Mochiyuki/gosky/client"
	"github.com/1Mochiyuki/gosky/errs"
	"github.com/1Mochiyuki/gosky/ui/send"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type LoginScreenModel struct {
	error           error
	LoginComponents []textinput.Model
	focusState      int
	rememberState   bool
}

func InitLoginScreenModel() LoginScreenModel {
	m := LoginScreenModel{
		LoginComponents: make([]textinput.Model, 2), focusState: 0,
	}
	var t textinput.Model
	for i := range m.LoginComponents {
		t = textinput.New()
		t.Cursor.Style = focusedStyle
		t.Cursor.SetMode(cursor.CursorHide)

		switch i {
		case 0:
			t.Placeholder = "BlueSky Handle"
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
					return errors.New("App Pass Not Long Enough")
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
	handle = iota
	pass
	rememberMe
	submit
	filledCheckBox = "☑"
	emptyCheckbox  = "☐"
)

func (l LoginScreenModel) saveLogin() {
	if readErr := viper.ReadInConfig(); readErr != nil {
		l.error = errors.New("Issue reading config")
		return
	}
	viper.Set("gsky_app_pass", l.LoginComponents[pass].Value())
	if writeErr := viper.WriteConfig(); writeErr != nil {
		l.error = errors.New("Issue writing credentials")
	}
}

func (l LoginScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return l, tea.Quit
		case tea.KeyUp.String(), tea.KeyDown.String(), tea.KeyEnter.String(), tea.KeyCtrlJ.String(), tea.KeyCtrlK.String():
			s := msg.String()
			if s == tea.KeyEnter.String() {
				if l.focusState == submit {
					if l.rememberState {
						l.saveLogin()
					}
					agent := client.NewAgent(context.Background(), "", l.LoginComponents[handle].Value(), l.LoginComponents[pass].Value())
					if err := agent.Connect(); err != nil {
						if errors.Is(err, errs.IncorrectCredentials{}) {
							l.error = errs.NewIncorrectCredentialsError()
							return l, cmd
						}
						l.error = err
					}
					return send.Model(agent), cmd

				}
				if l.focusState == rememberMe {
					l.rememberState = !l.rememberState
				}
			}

			if s == tea.KeyUp.String() || s == tea.KeyCtrlK.String() {
				l.focusState--
			} else if s == tea.KeyDown.String() || s == tea.KeyCtrlJ.String() {
				l.focusState++
			}

			if l.focusState > len(l.LoginComponents)+1 {
				l.focusState = handle
			} else if l.focusState == submit || l.focusState < handle {
				l.focusState = submit
			} else if l.focusState == rememberMe {
				l.focusState = rememberMe
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
	if l.focusState == submit {
		button = &focusedButton
	}
	btn := fmt.Sprintf("%v", *button)
	rememberMeStr := blurredRememberMe
	if l.focusState == rememberMe {
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
		noStyle.Width(50).Border(lipgloss.RoundedBorder(), true, true, true, true).Align(lipgloss.Center).Render(
			"  [ GoSky Scheduler ]",
			fmt.Sprintf("\n%s\n", fmtNilErrAsEmpty(l.error)),

			noStyle.Margin(1, 0).Width(35).Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true, true, true, true).Align(lipgloss.Center).Render(

				l.LoginComponents[handle].View(),
			),
			noStyle.Margin(1, 0).Width(35).Align(lipgloss.Center).Border(lipgloss.RoundedBorder(), true, true, true, true).Align(lipgloss.Center).Render(
				l.LoginComponents[pass].View(),
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
