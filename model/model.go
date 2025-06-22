package model

import (
	"fmt"
	"github.com/Daty26/pomodoro/data"
	"github.com/Daty26/pomodoro/storage"
	"github.com/Daty26/pomodoro/ui"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os/exec"
	"time"
)

type Model struct {
	Remaining     int
	Running       bool
	Working       bool
	Sessions      int
	Pause         bool
	Progress      progress.Model
	WorkDuration  int
	ShortDuration int
	LongDuration  int
	Logs          []data.Log
}

type tickMsg struct{}

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
func showDial() {
	err := exec.Command("osascript", "-e", `display dialog "Session finished!" with title "Pomodoro" buttons {"OK"}`).Run()
	if err != nil {
		fmt.Println("Warning: Failed to show dialog:", err)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.Running && !m.Pause && m.Remaining > 0 {
			m.Remaining--
			return m, Tick()
		}
		if m.Running && m.Remaining == 0 {
			compDuration := 0
			finishedWorking := m.Working
			if m.Working {
				compDuration = m.WorkDuration
				m.Sessions++
				if m.Sessions%4 == 0 {
					m.Remaining = m.LongDuration
				} else {
					m.Remaining = m.ShortDuration
				}
				m.Working = false
			} else {
				if m.Sessions%4 == 0 {
					compDuration = m.LongDuration
				} else {
					compDuration = m.ShortDuration
				}
				m.Working = true
				m.Remaining = m.WorkDuration
			}
			logging := data.Log{
				Timestamp: time.Now(),
				Working:   finishedWorking,
				Duration:  compDuration,
			}
			//fmt.Println(logging)
			showDial()
			m.Logs = append(m.Logs, logging)
			err := storage.FileSave(m.Logs)
			if err != nil {
				log.Fatal(err)
				return nil, nil
			}
			return m, Tick()
		}
		return m, Tick()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if !m.Running {
				m.Running = true
				m.Working = true
				m.Pause = false
				m.Remaining = m.WorkDuration
				return m, Tick()
			}
		case " ":
			if m.Running {
				m.Pause = !m.Pause
			}
		}
	}
	return m, nil
}
func (m Model) View() string {
	minutes := m.Remaining / 60
	seconds := m.Remaining % 60
	mode := "Waiting"
	if m.Running {
		mode = "Resting"
		if m.Working {
			mode = "Working"
		}
	}
	timerText := fmt.Sprintf("%s: %02d:%02d", mode, minutes, seconds)

	var total int
	if m.Working {
		total = m.WorkDuration
	} else {
		if m.Sessions%4 == 0 && !m.Working && m.Sessions > 0 {
			total = m.LongDuration
		} else {
			total = m.ShortDuration
		}
	}
	percent := 1.0
	if m.Running {
		percent = 1.0 - float64(m.Remaining)/float64(total)
	}
	bar := m.Progress.ViewAs(percent)

	view := ui.TitleStyle.Render("Pomodoro timer") + "\n\n" + timerText + "\n\n" + bar + "\n\n" + "Press 'Enter' to start, 'Space' to pause, q to exit"
	return ui.CenterStyle.Render(view)
}
