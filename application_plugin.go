package ctk

import (
	"context"
	"fmt"
	"plugin"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
	"github.com/go-curses/cdk/lib/sync"
)

// NewApplicationFromPlugin constructs an Application from a Go plugin shared
// object file. This file must export a number of variables and functions, as
// follows:
//
//  Variables
//  ---------
//   CtkName         the name of the program binary
//   CtkUsage        brief summary of what the program is used for
//   CtkDescription  long description of the program
//   CtkVersion      the version number to report
//   CtkTag          a machine tag name, ie: "tld.domain.name" ([-.a-zA-Z])
//   CtkTitle        the user-visible name of the program
//   CtkTtyPath      the unix tty device path for display capture
//
//  Functions
//  ---------
//   // early initialization stage
//   Init(Application)
//   // UI initialization stage
//   Startup(Application, cdk.Display, context.Context, context.CancelFunc, *sync.WaitGroup) cenums.EventFlag
//   // Application shutdown stage (display destroyed, logging still works, stdin/stdout restored)
//   Shutdown() cenums.EventFlag
//
// Two things must happen in the Startup function for CTK to correctly start
// rendering the user-interface. The first is to ensure that the function
// returns EVENT_PASS for the SignalStartup to complete. Returning EVENT_STOP
// will abort the entire startup process and immediately result in system
// shutdown. The second thing that must happen for startup to complete is that
// Application.NotifyStartupComplete must be called before returning EVENT_PASS.
// This will notify the Application and enable the Display to begin rending to
// the Screen instance.
func NewApplicationFromPlugin(path string) (app Application, err error) {
	var plug *plugin.Plugin
	if plug, err = plugin.Open(path); err != nil {
		return nil, err
	}
	var initFn ApplicationInitFn
	if initFn, err = lookupPluginInitFn(plug); err != nil {
		return nil, err
	}
	var startupFn ApplicationStartupFn
	if startupFn, err = lookupPluginStartupFn(plug); err != nil {
		return nil, err
	}
	var shutdownFn ApplicationShutdownFn
	if shutdownFn, err = lookupPluginShutdownFn(plug); err != nil {
		return nil, err
	}
	var name, usage, description, version, tag, title, ttyPath string
	if name, err = lookupPluginStringValue("CtkName", plug); err != nil {
		return nil, err
	}
	if usage, err = lookupPluginStringValue("CtkUsage", plug); err != nil {
		return nil, err
	}
	if description, err = lookupPluginStringValue("CtkDescription", plug); err != nil {
		return nil, err
	}
	if version, err = lookupPluginStringValue("CtkVersion", plug); err != nil {
		return nil, err
	}
	if tag, err = lookupPluginStringValue("CtkTag", plug); err != nil {
		return nil, err
	}
	if title, err = lookupPluginStringValue("CtkTitle", plug); err != nil {
		return nil, err
	}
	if ttyPath, err = lookupPluginStringValue("CtkTtyPath", plug); err != nil {
		return nil, err
	}
	app = NewApplication(name, usage, description, version, tag, title, ttyPath)
	if startupFn != nil {
		app.Connect(cdk.SignalStartup, ApplicationPluginStartupHandle, func(_ []interface{}, argv ...interface{}) cenums.EventFlag {
			_ = app.Disconnect(cdk.SignalStartup, ApplicationPluginStartupHandle)
			if app, display, ctx, cancel, wg, ok := ArgvApplicationSignalStartup(argv...); ok {
				return startupFn(app, display, ctx, cancel, wg)
			}
			return cenums.EVENT_STOP
		})
	}
	if shutdownFn != nil {
		app.Connect(cdk.SignalShutdown, ApplicationPluginShutdownHandle, func(_ []interface{}, argv ...interface{}) cenums.EventFlag {
			_ = app.Disconnect(cdk.SignalShutdown, ApplicationPluginShutdownHandle)
			return shutdownFn()
		})
	}
	if initFn != nil {
		initFn(app)
	}
	return
}

func ArgvApplicationSignalStartup(argv ...interface{}) (app Application, display cdk.Display, ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ok bool) {
	if len(argv) == 5 {
		if app, ok = argv[0].(Application); ok {
			if display, ok = argv[1].(cdk.Display); ok {
				if ctx, ok = argv[2].(context.Context); ok {
					if cancel, ok = argv[3].(context.CancelFunc); ok {
						if wg, ok = argv[4].(*sync.WaitGroup); ok {
							return
						}
						cancel = nil
					}
					ctx = nil
				}
				display = nil
			}
			app = nil
		}
	}
	return
}

func WithArgvApplicationSignalStartup(startupFn ApplicationStartupFn) cdk.SignalListenerFn {
	return func(_ []interface{}, argv ...interface{}) cenums.EventFlag {
		if app, display, ctx, cancel, wg, ok := ArgvApplicationSignalStartup(argv...); ok {
			return startupFn(app, display, ctx, cancel, wg)
		}
		return cenums.EVENT_STOP
	}
}

func WithArgvNoneWithFlagsSignal(fn func() cenums.EventFlag) cdk.SignalListenerFn {
	return func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
		return fn()
	}
}

func WithArgvNoneSignal(fn func(), eventFlag cenums.EventFlag) cdk.SignalListenerFn {
	return func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
		fn()
		return eventFlag
	}
}

func lookupPluginStringValue(key string, plug *plugin.Plugin) (value string, err error) {
	var symbol plugin.Symbol
	if symbol, err = plug.Lookup(key); err != nil {
		return "", err
	}
	var ok bool
	var ptr *string
	if ptr, ok = symbol.(*string); !ok {
		return "", fmt.Errorf("%v value stored in plugin is not of *string type: %v (%T)", key, symbol, symbol)
	}
	value = *ptr
	return
}

func lookupPluginInitFn(plug *plugin.Plugin) (initFn ApplicationInitFn, err error) {
	var symbol plugin.Symbol
	if symbol, err = plug.Lookup(ApplicationPluginInitFnName); err != nil {
		return nil, err
	}
	if symPtr, ok := symbol.(ApplicationInitFn); !ok {
		return nil, fmt.Errorf(
			"%v func stored in plugin is not of ApplicationInitFn type: %v (%T)\n",
			ApplicationPluginInitFnName,
			symbol,
			symbol,
		)
	} else {
		initFn = symPtr
	}
	return
}

func lookupPluginStartupFn(plug *plugin.Plugin) (startupFn ApplicationStartupFn, err error) {
	var symbol plugin.Symbol
	if symbol, err = plug.Lookup(ApplicationPluginStartupFnName); err != nil {
		return nil, err
	}
	if symPtr, ok := symbol.(ApplicationStartupFn); !ok {
		return nil, fmt.Errorf(
			"%v func stored in plugin is not of ApplicationStartupFn type: %v (%T)\n",
			ApplicationPluginStartupFnName,
			symbol,
			symbol,
		)
	} else {
		startupFn = symPtr
	}
	return
}

func lookupPluginShutdownFn(plug *plugin.Plugin) (shutdownFn ApplicationShutdownFn, err error) {
	var symbol plugin.Symbol
	if symbol, err = plug.Lookup(ApplicationPluginShutdownFnName); err != nil {
		return nil, err
	}
	if symPtr, ok := symbol.(ApplicationShutdownFn); !ok {
		return nil, fmt.Errorf(
			"%v func stored in plugin is not of ApplicationShutdownFn type: %v (%T)\n",
			ApplicationPluginShutdownFnName,
			symbol,
			symbol,
		)
	} else {
		shutdownFn = symPtr
	}
	return
}

var ApplicationPluginInitFnName = "CtkInit"
var ApplicationPluginStartupFnName = "CtkStartup"
var ApplicationPluginShutdownFnName = "CtkShutdown"

const ApplicationPluginStartupHandle = "application-plugin-startup-handler"
const ApplicationPluginShutdownHandle = "application-plugin-shutdown-handler"

type ApplicationInitFn = func(app Application)

type ApplicationStartupFn = func(
	app Application,
	display cdk.Display,
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
) cenums.EventFlag

type ApplicationShutdownFn = func() cenums.EventFlag
