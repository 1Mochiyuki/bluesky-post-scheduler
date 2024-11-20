package picker

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ImagePickerModel struct {
	fp            filepicker.Model
	err           error
	SelectedFiles []string
	quitting      bool
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m ImagePickerModel) Init() tea.Cmd {
	return m.fp.Init()
}

var l = logger.Get()

func (m ImagePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.fp, cmd = m.fp.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.fp.DidSelectFile(msg); didSelect {
		l.Debug().Int("length", len(m.SelectedFiles)).Msg("selected files length")
		// Get the path of the selected file.
		if len(m.SelectedFiles) > 0 {
			for idx, s := range m.SelectedFiles {
				l.Debug().Str("selected", path).Str("file", s).Msg("in for loop")
				if path == s {
					m.SelectedFiles = append(m.SelectedFiles[:idx], m.SelectedFiles[idx+1:]...)
					return m, cmd
				}
			}
		}
		if len(m.SelectedFiles) < 4 {
			l.Debug().Msg("length less than 4 ")
			m.SelectedFiles = append(m.SelectedFiles, path)
		} else {
			l.Debug().Msg("max files chosen")
			m.err = errors.New("maximum files chosen")
			return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
		}

		return m, cmd

	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.fp.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m ImagePickerModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\nFiles Selected: ")
	selectedFiles := m.SelectedFiles
	m.fp.Styles.DisabledFile = lipgloss.NewStyle().Foreground(lipgloss.Color(""))

	if len(selectedFiles) > 0 {
		for _, file := range m.SelectedFiles {

			name := fileName(file)
			s.WriteString(name + ", ")
		}
	}
	if m.err != nil {
		s.WriteString(m.fp.Styles.DisabledFile.Render(m.err.Error()))
	}
	s.WriteString("\n\n" + m.fp.View() + "\n")
	return s.String()
}

func fileName(absFilePath string) string {
	parts := strings.Split(absFilePath, string(os.PathSeparator))
	return parts[len(parts)-1]
}

func NewPickerModel() ImagePickerModel {
	fp := filepicker.New()
	fp.ShowHidden = false

	fp.DirAllowed = true
	fp.Height = 10
	fp.AllowedTypes = []string{".png", ".jpg", ".jpeg", ".mp4", ".mov"}

	home, _ := os.UserHomeDir()
	fp.CurrentDirectory = filepath.Join(home, "Pictures", "backgrounds")

	m := ImagePickerModel{
		fp:            fp,
		SelectedFiles: make([]string, 0, 4),
	}
	return m
}
