package ui

import "github.com/charmbracelet/lipgloss"

var CenterStyle = lipgloss.NewStyle().Width(50).Align(lipgloss.Left)
var TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(1, 0, 1, 0).Margin(2, 1, 1, 1).Align(lipgloss.Center)
