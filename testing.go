package ctk

import (
	"github.com/go-curses/cdk"
)

func TestingWithCtkWindow(d cdk.Display) error {
	w := NewWindowWithTitle(d.GetTitle())
	d.SetActiveWindow(w)
	return nil
}
