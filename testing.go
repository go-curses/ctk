package ctk

import (
	"context"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/sync"
	"github.com/go-curses/cdk/log"
)

type AppFn func(app Application)

func WithApp(initFn cdk.SignalListenerFn, action AppFn) func() {
	return func() {
		app := NewApplication(
			"AppName", "AppUsage",
			"AppDesc", "v0.0.0",
			"app-tag", "AppTitle",
			cdk.OffscreenTtyPath,
		)
		app.Connect(cdk.SignalStartup, "testing-withapp-init-fn-handler", initFn)
		defer func() {
			if app != nil {
				app.Destroy()
			}
			app = nil
		}()
		app.SetupDisplay()
		ctx, cancel := context.WithCancel(context.Background())
		if f := app.Emit(cdk.SignalStartup, app, app.Display(), ctx, cancel, &sync.WaitGroup{}); f == enums.EVENT_STOP {
			log.ErrorF("WithApp startup listeners requested EVENT_STOP")
		} else {
			action(app)
		}
	}
}

func TestingWithCtkWindow(_ []interface{}, argv ...interface{}) enums.EventFlag {
	if _, d, _, _, _, ok := ArgvApplicationSignalStartup(argv...); ok {
		w := NewWindowWithTitle(d.GetTitle())
		d.SetActiveWindow(w)
	}
	return enums.EVENT_PASS
}
