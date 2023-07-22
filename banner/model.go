package banner

import (
	"fmt"
	"github.com/Anthony-Jhoiro/bastion/colors"
	tea "github.com/charmbracelet/bubbletea"
)

var hydraBanner = colors.Heading("  ______           _   _             \n  | ___ \\         | | (_)            \n  | |_/ / __ _ ___| |_ _  ___  _ __  \n  | ___ \\/ _` / __| __| |/ _ \\| '_ \\ \n  | |_/ / (_| \\__ \\ |_| | (_) | | | |\n  \\____/ \\__,_|___/\\__|_|\\___/|_| |_|\n                                     ")

type model struct {
}

// NewBannerModel is a bubbletea model that displays the "Bastion" banner. It does not handle any events nor requires
// init
func NewBannerModel() tea.Model {
	return model{}
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
