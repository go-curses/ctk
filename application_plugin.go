package ctk

import (
	"context"
	"fmt"
	"plugin"
	"sync"

	"github.com/go-curses/cdk"
	cenums "github.com/go-curses/cdk/lib/enums"
)

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

func WrapArgvApplicationSignalStartup(startupFn ApplicationPluginStartupFn) cdk.SignalListenerFn {
	return func(_ []interface{}, argv ...interface{}) cenums.EventFlag {
		if app, display, ctx, cancel, wg, ok := ArgvApplicationSignalStartup(argv...); ok {
			return startupFn(app, display, ctx, cancel, wg)
		}
		return cenums.EVENT_STOP
	}
}

func WrapArgvNoneWithFlagsSignal(fn func() cenums.EventFlag) cdk.SignalListenerFn {
	return func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
		return fn()
	}
}

func WrapArgvNoneSignal(fn func(), eventFlag cenums.EventFlag) cdk.SignalListenerFn {
	return func(_ []interface{}, _ ...interface{}) cenums.EventFlag {
		fn()
		return eventFlag
	}
}

func NewApplicationFromPlugin(path string) (app Application, err error) {
	var plug *plugin.Plugin
	if plug, err = plugin.Open(path); err != nil {
		return nil, err
	}
	var initFn ApplicationPluginInitFn
	if initFn, err = lookupPluginInitFn(plug); err != nil {
		return nil, err
	}
	var startupFn ApplicationPluginStartupFn
	if startupFn, err = lookupPluginStartupFn(plug); err != nil {
		return nil, err
	}
	var shutdownFn ApplicationPluginShutdownFn
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

func lookupPluginInitFn(plug *plugin.Plugin) (initFn ApplicationPluginInitFn, err error) {
	var symbol plugin.Symbol
	if symbol, err = plug.Lookup(ApplicationPluginInitFnName); err != nil {
		return nil, err
	}
	if symPtr, ok := symbol.(ApplicationPluginInitFn); !ok {
		return nil, fmt.Errorf(
			"%v func stored in plugin is not of ApplicationPluginInitFn type: %v (%T)\n",
			ApplicationPluginInitFnName,
			symbol,
			symbol,
		)
	} else {
		initFn = symPtr
	}
	return
}

func lookupPluginStartupFn(plug *plugin.Plugin) (startupFn ApplicationPluginStartupFn, err error) {
	var symbol plugin.Symbol
	if symbol, err = plug.Lookup(ApplicationPluginStartupFnName); err != nil {
		return nil, err
	}
	if symPtr, ok := symbol.(ApplicationPluginStartupFn); !ok {
		return nil, fmt.Errorf(
			"%v func stored in plugin is not of ApplicationPluginStartupFn type: %v (%T)\n",
			ApplicationPluginStartupFnName,
			symbol,
			symbol,
		)
	} else {
		startupFn = symPtr
	}
	return
}

func lookupPluginShutdownFn(plug *plugin.Plugin) (shutdownFn ApplicationPluginShutdownFn, err error) {
	var symbol plugin.Symbol
	if symbol, err = plug.Lookup(ApplicationPluginShutdownFnName); err != nil {
		return nil, err
	}
	if symPtr, ok := symbol.(ApplicationPluginShutdownFn); !ok {
		return nil, fmt.Errorf(
			"%v func stored in plugin is not of ApplicationPluginShutdownFn type: %v (%T)\n",
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

type ApplicationPluginInitFn = func(app Application)

var ApplicationPluginStartupFnName = "CtkStartup"

type ApplicationPluginStartupFn = func(
	app Application,
	display cdk.Display,
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
) cenums.EventFlag

var ApplicationPluginShutdownFnName = "CtkShutdown"

type ApplicationPluginShutdownFn = func() cenums.EventFlag

const ApplicationPluginStartupHandle = "application-plugin-startup-handler"

const ApplicationPluginShutdownHandle = "application-plugin-shutdown-handler"
