package banner

import (
	"bastion/colors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

var hydraBanner = colors.Heading("  ______           _   _             \n  | ___ \\         | | (_)            \n  | |_/ / __ _ ___| |_ _  ___  _ __  \n  | ___ \\/ _` / __| __| |/ _ \\| '_ \\ \n  | |_/ / (_| \\__ \\ |_| | (_) | | | |\n  \\____/ \\__,_|___/\\__|_|\\___/|_| |_|\n                                     ")

type model struct {
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("%v\n\n", hydraBanner)
}

func NewBannerModel() tea.Model {
	return model{}
}
