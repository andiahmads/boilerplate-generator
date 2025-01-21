package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/ease"
	"os/exec"
	"strconv"
	"time"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

func initModel() model {
	ti := textinput.New()
	ti.Placeholder = "project name:"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	m := model{Choice: 0,
		Chosen: false,
		Ticks:  60, Frames: 0, Progress: 0, Loaded: false, Quitting: false,
		textInput: ti, Submitted: false, Spinner: sp}

	return m
}

type (
	tickMsg         struct{}
	frameMsg        struct{}
	CreateFolderMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return frameMsg{}
	})

}

type model struct {
	Choice     int
	Chosen     bool
	Ticks      int
	Frames     int
	Progress   float64
	Loaded     bool
	Quitting   bool
	textInput  textinput.Model
	Submitted  bool
	Spinner    spinner.Model
	Done       bool
	FolderName string
}

func (m model) Init() tea.Cmd {
	return tick()
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			switch m.Choice {
			case 0:
				m.FolderName = m.textInput.Value()
				if m.FolderName != "" {
					m.Loaded = true
					m.Done = true
					return m, tea.Batch(m.Spinner.Tick, createFolderCmd(m.FolderName))
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case CreateFolderMsg:
		m.Loaded = false
		m.Done = true
		m.Quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	if !m.Chosen {
		return updateChoices(msg, m)
	}

	var cmd1, cmd2 tea.Cmd
	m.textInput, cmd1 = m.textInput.Update(msg)

	updatedModel, cmd2 := updateChosen(msg, m)

	// Gabungkan semua command agar dijalankan oleh Bubble Tea
	return updatedModel, tea.Batch(cmd1, cmd2)
}

func (m model) View() string {
	var s string
	if m.Quitting {
		data := fmt.Sprintf("\n🎉 Folder '%s' berhasil dibuat!\n", m.FolderName)
		return data + "\n See you letter!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
		if m.Done {
			s = loadAnimation(m)
		}
	}
	return mainStyle.Render("\n" + s + "\n\n")
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			return m, frame()
		}
	case tickMsg:
		if m.Ticks == 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}

	return m, nil
}

func choicesView(m model) string {
	c := m.Choice
	tpl := "Whate do you plan?\n\n"
	tpl += "%s\n\n"
	tpl += "Program quites in %s seconde\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle + subtleStyle.Render("enter: chose") + dotStyle + subtleStyle.Render("q, esc: quite")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		checkbox("1.Create BoilerPlate", c == 0),
		checkbox("2.add package modular", c == 1),
		checkbox("3.Read something", c == 2),
		checkbox("4.See friends", c == 3),
	)
	return fmt.Sprintf(tpl, choices, ticksStyle.Render(strconv.Itoa(m.Ticks)))
}

func chosenView(m model) string {
	var msg string
	switch m.Choice {
	case 0:
		// msg = fmt.Sprintf("create pkg?\n\nCool, we'll need %s and %s...", keywordStyle.Render("libgarden"), keywordStyle.Render("vegeutils"))
		msg = m.createBoilerPlate()

	case 1:
		msg = fmt.Sprintf("create pkg?\n\nCool, we'll need %s and %s...", keywordStyle.Render("libgarden"), keywordStyle.Render("vegeutils"))
	case 2:
		msg = fmt.Sprintf("Reading time?\n\nOkay, cool, then we’ll need a library. Yes, an %s.", keywordStyle.Render("actual library"))
	default:
		msg = fmt.Sprintf("It’s always good to see friends.\n\nFetching %s and %s...", keywordStyle.Render("social-skills"), keywordStyle.Render("conversationutils"))
	}
	// label := "Downloading..."
	// if m.Loaded {
	// 	label = fmt.Sprintf("Downloaded. Exiting in %s seconds...", ticksStyle.Render(strconv.Itoa(m.Ticks)))
	// }
	return msg

	// return msg + "\n\n" + label + "\n" + progressbar(m.Progress) + "%"
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case frameMsg:
		if !m.Loaded {
			m.Frames++
			m.Progress = ease.OutBounce(float64(m.Frames) / float64(100))
			if m.Progress >= 1 {
				m.Progress = 1
				m.Loaded = true
				m.Done = false
				m.Ticks = 5
				return m, tick()
			}
			return m, frame()
		}
	case tickMsg:
		if m.Loaded {
			if m.Ticks == 0 {
				m.Quitting = true
				return m, tea.Quit
			}
			m.Ticks--
			return m, tick()
		}
	}
	return m, nil
}

func (m model) createBoilerPlate() string {
	return fmt.Sprintf(
		"project name, please use `- or _` for separator?\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"

}

func createFolderCmd(folderName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("mkdir", folderName)
		err := cmd.Run()
		if err != nil {
			fmt.Println("Gagal membuat folder:", err)
		}
		return CreateFolderMsg{}
	}
}

func loadAnimation(m model) string {
	var msg string
	label := "generate package..."

	if m.Loaded {
		label = fmt.Sprintf("Downloaded. Exiting in %s seconds...", ticksStyle.Render(strconv.Itoa(m.Ticks)))
	}
	return msg + "\n\n" + label + "\n" + progressbar(m.Progress) + "%"
}
