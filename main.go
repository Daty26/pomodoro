package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"time"
)

type model struct {
	remaining     int
	running       bool
	working       bool
	sessions      int
	pause         bool
	progress      progress.Model
	workDuration  int
	shortDuration int
	longDuration  int
}
type tickMsg struct{}

var centerStyle = lipgloss.NewStyle().Width(50).Align(lipgloss.Left)
var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).
	//Background(lipgloss.Color("#7D56F4")).
	Padding(2, 0, 1, 0).Margin(0, 0).Align(lipgloss.Center)

func playSound() {
	exec.Command("say", "Time is up").Run()
}

func sendNot(title, message string) {
	exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "%s"`, title, message)).Run()
}

func (m model) Init() tea.Cmd {
	return nil
}
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.running && !m.pause && m.remaining > 0 {
			m.remaining--
			return m, tick()
		}
		if m.running && m.remaining == 0 {
			if m.working {
				m.sessions++
				if m.sessions%4 == 0 {
					m.remaining = m.longDuration
				} else {
					m.remaining = m.shortDuration
				}
				m.working = false
			} else {
				m.working = true
				m.remaining = m.workDuration
			}
			playSound()
			sendNot("Pomodoro", "Work session finished")
			return m, tick()
		}
		return m, tick()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if !m.running {
				m.running = true
				m.working = true
				m.pause = false
				m.remaining = m.workDuration
				return m, tick()
			}
		case "space":
			if m.running {
				m.pause = !m.pause
			}
		}
	}
	return m, nil
}
func (m model) View() string {
	minutes := m.remaining / 60
	seconds := m.remaining % 60
	mode := "Waiting"
	if m.running {
		mode = "Resting"
		if m.working {
			mode = "Working"
		}
	}
	timerText := fmt.Sprintf("%s: %02d:%02d", mode, minutes, seconds)

	var total int
	if m.working {
		total = m.workDuration
	} else {
		if m.sessions%4 == 0 && !m.working && m.sessions > 0 {
			total = m.longDuration
		} else {
			total = m.shortDuration
		}
	}
	percent := 1.0
	if m.running {
		percent = 1.0 - float64(m.remaining)/float64(total)
	}
	bar := m.progress.ViewAs(percent)

	view := titleStyle.Render("Pomodoro timer") + "\n\n" + timerText + "\n\n" + bar + "\n\n" + "Press s to start, p to pause, q to exit"
	return centerStyle.Render(view)
}
func main() {
	work := flag.Int("work", 25, "Work session minutes")
	short := flag.Int("short", 5, "sshort session minutes")
	long := flag.Int("long", 15, "long session minutes")

	flag.Parse()
	p := tea.NewProgram(model{
		workDuration:  *work * 60,
		shortDuration: *short * 60,
		longDuration:  *long * 60,
		progress:      progress.New(),
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("There has been an error %v", err)
		os.Exit(1)
	}
}
