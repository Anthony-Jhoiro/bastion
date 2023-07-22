package discovery

import (
	"fmt"
)

type Host struct {
	Name string
	Ip   string
	Up   bool
}

func (h Host) Title() string {
	if h.Up {
		return fmt.Sprintf("%v ✅", h.Name)
	}
	return h.Name
}
func (h Host) Description() string { return h.Ip }
func (h Host) FilterValue() string { return h.Name }
