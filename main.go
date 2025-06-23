package main

import (
	"flag"
	"fmt"
	"github.com/Daty26/pomodoro/model"
	"github.com/Daty26/pomodoro/stats"
	"github.com/Daty26/pomodoro/storage"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	work := flag.Int("work", 25, "Work session minutes")
	short := flag.Int("short", 5, "short session minutes")
	long := flag.Int("long", 15, "long session minutes")
	reset := flag.Bool("reset", false, "Reset all logs")

	flag.Parse()
	if *reset {
		err := os.Remove("pomodoro_logs.json")
		if err != nil {
			fmt.Println("Error occured: " + err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}
	logs := storage.LoadLogs()
	stats.ShowStats(logs)
	p := tea.NewProgram(model.Model{
		WorkDuration:  *work * 60,
		ShortDuration: *short * 60,
		LongDuration:  *long * 60,
		Progress:      progress.New(),
		Logs:          logs,
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("There has been an error %v", err)
		os.Exit(1)
	}
}
