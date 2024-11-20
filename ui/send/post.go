package send

import (
	"fmt"

	"github.com/1Mochiyuki/gosky/api/posts"
	"github.com/1Mochiyuki/gosky/client"
	"github.com/1Mochiyuki/gosky/ui/picker"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PostModel struct {
	PostArea  textarea.Model
	err       error
	Agent     client.BskyAgent
	ResultMsg string
}

func Model(connectedAgent client.BskyAgent) PostModel {
	ta := textarea.New()
	ta.Placeholder = "Send A Post..."
	ta.CharLimit = 300

	ta.SetWidth(300)
	ta.SetHeight(6)
	ta.Focus()

	return PostModel{
		Agent:     connectedAgent,
		PostArea:  ta,
		err:       nil,
		ResultMsg: "",
	}
}

func (m PostModel) Init() tea.Cmd {
	return textinput.Blink
}

type (
	errMsg error
)

func (m PostModel) charCount() string {
	remainingChars := m.PostArea.CharLimit - m.PostArea.Length()
	return fmt.Sprintf("Chars Left: %d", remainingChars)
}

func (m PostModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String():
			return m, tea.Quit
		case tea.KeyCtrlP.String():

			post, err := posts.NewPostBuilder(m.PostArea.Value()).CreatePost()
			if err != nil {
				m.err = err
				return m, nil
			}
			_, _, err = m.Agent.CreatePost(post)
			if err != nil {
				fmt.Println(err)
				m.err = err
				m.ResultMsg = err.Error()
				return m, nil
			}
			m.ResultMsg = "Success :D"
			m.PostArea.SetValue("")
		case tea.KeyCtrlI.String():
			model := picker.NewPickerModel()
			return model, model.Init()
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	m.PostArea, cmd = m.PostArea.Update(msg)
	return m, cmd
}

func (m PostModel) View() string {
	return fmt.Sprintf(
		"What would you like to post to BlueSky?\n\n%s\n%s\n\n\n%s\n\n%s",
		m.PostArea.View(),
		m.charCount(),
		"(esc to quit)",
		m.ResultMsg,
	) + "\n"
}
