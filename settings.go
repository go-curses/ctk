package ctk

import (
	"regexp"
	"strings"
	"time"

	"github.com/go-curses/cdk"
	"github.com/go-curses/cdk/lib/enums"
)

const TypeSettings cdk.CTypeTag = "ctk-settings"

func init() {
	_ = cdk.TypesManager.AddType(TypeSettings, func() interface{} { return nil })
}

var ctkDefaultSettings *CSettings

// Settings Hierarchy:
//	Object
//	  +- Settings
type Settings interface {
	Object

	LoadFromString(rc string) (err error)
	GetAlternativeButtonOrder() (value bool)
	GetAlternativeSortArrows() (value bool)
	GetColorPalette() (value string)
	GetColorScheme() (value string)
	GetCursorBlink() (value bool)
	GetCursorBlinkTime() (value time.Duration)
	GetCursorBlinkTimeout() (value time.Duration)
	GetCursorThemeName() (value string)
	GetDndDragThreshold() (value time.Duration)
	GetDoubleClickDistance() (value int)
	GetDoubleClickTime() (value time.Duration)
	GetEnableAccels() (value bool)
	GetEnableMnemonics() (value bool)
	GetEnableTooltips() (value bool)
	GetEntryPasswordHintTimeout() (value time.Duration)
	GetEntrySelectOnFocus() (value bool)
	GetErrorBell() (value bool)
	GetFallbackIconTheme() (value string)
	GetFileChooserBackend() (value string)
	GetIconThemeName() (value string)
	GetImModule() (value string)
	GetImPreeditStyle() (value interface{})
	GetImStatusStyle() (value interface{})
	GetKeyThemeName() (value string)
	GetKeynavCursorOnly() (value bool)
	GetKeynavWrapAround() (value bool)
	GetLabelSelectOnFocus() (value bool)
	GetMenuBarAccel() (value string)
	GetMenuBarPopupDelay() (value time.Duration)
	GetMenuImages() (value bool)
	GetMenuPopdownDelay() (value time.Duration)
	GetMenuPopupDelay() (value time.Duration)
	GetModules() (value string)
	GetPrimaryButtonWarpsSlider() (value bool)
	GetScrolledWindowPlacement() (value interface{})
	GetShowInputMethodMenu() (value bool)
	GetShowUnicodeMenu() (value bool)
	GetThemeName() (value string)
	GetTimeoutExpand() (value time.Duration)
	GetTimeoutInitial() (value time.Duration)
	GetTimeoutRepeat() (value time.Duration)
	GetToolbarStyle() (value interface{})
	GetTooltipBrowseModeTimeout() (value time.Duration)
	GetTooltipBrowseTimeout() (value time.Duration)
	GetTooltipTimeout() (value time.Duration)
	GetTouchscreenMode() (value bool)
	SetCtkAlternativeButtonOrder(value bool)
	SetCtkAlternativeSortArrows(value bool)
	SetCtkColorPalette(value string)
	SetCtkColorScheme(value string)
	SetCtkCursorBlink(value bool)
	SetCtkCursorBlinkTime(value time.Duration)
	SetCtkCursorBlinkTimeout(value time.Duration)
	SetCtkCursorThemeName(value string)
	SetCtkDndDragThreshold(value time.Duration)
	SetCtkDoubleClickDistance(value int)
	SetCtkDoubleClickTime(value time.Duration)
	SetCtkEnableAccels(value bool)
	SetCtkEnableMnemonics(value bool)
	SetCtkEnableTooltips(value bool)
	SetCtkEntryPasswordHintTimeout(value time.Duration)
	SetCtkEntrySelectOnFocus(value bool)
	SetCtkErrorBell(value bool)
	SetCtkFallbackIconTheme(value string)
	SetCtkFileChooserBackend(value string)
	SetCtkIconThemeName(value string)
	SetCtkImModule(value string)
	SetCtkImPreeditStyle(value interface{})
	SetCtkImStatusStyle(value interface{})
	SetCtkKeyThemeName(value string)
	SetCtkKeynavCursorOnly(value bool)
	SetCtkKeynavWrapAround(value bool)
	SetCtkLabelSelectOnFocus(value bool)
	SetCtkMenuBarAccel(value string)
	SetCtkMenuBarPopupDelay(value time.Duration)
	SetCtkMenuImages(value bool)
	SetCtkMenuPopdownDelay(value time.Duration)
	SetCtkMenuPopupDelay(value time.Duration)
	SetCtkModules(value string)
	SetCtkPrimaryButtonWarpsSlider(value bool)
	SetCtkScrolledWindowPlacement(value interface{})
	SetCtkShowInputMethodMenu(value bool)
	SetCtkShowUnicodeMenu(value bool)
	SetCtkThemeName(value string)
	SetCtkTimeoutExpand(value time.Duration)
	SetCtkTimeoutInitial(value time.Duration)
	SetCtkTimeoutRepeat(value time.Duration)
	SetCtkToolbarStyle(value interface{})
	SetCtkTooltipBrowseModeTimeout(value time.Duration)
	SetCtkTooltipBrowseTimeout(value time.Duration)
	SetCtkTooltipTimeout(value time.Duration)
	SetCtkTouchscreenMode(value bool)
}

var _ Settings = (*CSettings)(nil)

// The CSettings structure implements the Settings interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Settings objects
type CSettings struct {
	CObject
}

func GetDefaultSettings() (settings Settings) {
	if ctkDefaultSettings == nil {
		ctkDefaultSettings = &CSettings{}
		ctkDefaultSettings.Init()
	}
	return ctkDefaultSettings
}

// Settings object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Settings instance
func (s *CSettings) Init() (already bool) {
	if s.InitTypeItem(TypeSettings, s) {
		return true
	}
	s.CObject.Init()
	_ = s.InstallProperty(PropertyCtkAlternativeButtonOrder, cdk.BoolProperty, true, false)
	_ = s.InstallProperty(PropertyCtkAlternativeSortArrows, cdk.BoolProperty, true, false)
	_ = s.InstallProperty(PropertyCtkColorPalette, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkColorScheme, cdk.StringProperty, true, "")
	_ = s.InstallProperty(PropertyCtkCursorBlink, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkCursorBlinkTime, cdk.TimeProperty, true, 1200*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkCursorBlinkTimeout, cdk.TimeProperty, true, 2147483647*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkCursorThemeName, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkDndDragThreshold, cdk.IntProperty, true, 2)
	_ = s.InstallProperty(PropertyCtkDoubleClickDistance, cdk.IntProperty, true, 1)
	_ = s.InstallProperty(PropertyCtkDoubleClickTime, cdk.TimeProperty, true, 250*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkEnableAccels, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkEnableMnemonics, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkEnableTooltips, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkEntryPasswordHintTimeout, cdk.TimeProperty, true, 0*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkEntrySelectOnFocus, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkErrorBell, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkFallbackIconTheme, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkFileChooserBackend, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkIconThemeName, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkImModule, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkImPreeditStyle, cdk.StructProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkImStatusStyle, cdk.StructProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkKeyThemeName, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkKeynavCursorOnly, cdk.BoolProperty, true, false)
	_ = s.InstallProperty(PropertyCtkKeynavWrapAround, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkLabelSelectOnFocus, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkMenuBarAccel, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkMenuBarPopupDelay, cdk.TimeProperty, true, 0*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkMenuImages, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkMenuPopdownDelay, cdk.TimeProperty, true, 1000*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkMenuPopupDelay, cdk.TimeProperty, true, 225*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkModules, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkPrimaryButtonWarpsSlider, cdk.BoolProperty, true, false)
	_ = s.InstallProperty(PropertyCtkScrolledWindowPlacement, cdk.StructProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkShowInputMethodMenu, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkShowUnicodeMenu, cdk.BoolProperty, true, true)
	_ = s.InstallProperty(PropertyCtkThemeName, cdk.StringProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkTimeoutExpand, cdk.TimeProperty, true, 500*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkTimeoutInitial, cdk.TimeProperty, true, 200*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkTimeoutRepeat, cdk.TimeProperty, true, 20*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkToolbarStyle, cdk.StructProperty, true, nil)
	_ = s.InstallProperty(PropertyCtkTooltipBrowseModeTimeout, cdk.TimeProperty, true, 500*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkTooltipBrowseTimeout, cdk.TimeProperty, true, 60*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkTooltipTimeout, cdk.TimeProperty, true, 500*time.Millisecond)
	_ = s.InstallProperty(PropertyCtkTouchscreenMode, cdk.BoolProperty, true, false)
	return false
}

var rxCtkSettingsParseLine = regexp.MustCompile(`^\s*([-a-z]+?)\s*=\s*(.+?)\s*$`)

// LoadFromString parses the given string for key=value pairs, matching the
// CTK settings property names.
func (s *CSettings) LoadFromString(rc string) (err error) {
	keys := ctkSettingsPropertyKeys()
	lines := strings.Split(rc, "\n")
	for _, line := range lines {
		if rxCtkSettingsParseLine.MatchString(line) {
			m := rxCtkSettingsParseLine.FindAllString(line, -1)
			if len(m) != 3 {
				s.LogError("error parsing rc line: %v", line)
				continue
			}
			mk := cdk.Property(string(m[1]))
			found := false
			for _, key := range keys {
				if cdk.Property(key) == mk {
					found = true
					break
				}
			}
			if found {
				if prop := s.GetProperty(mk); prop != nil {
					if err = prop.SetFromString(m[2]); err != nil {
						s.LogErr(err)
					}
				} else {
					s.LogError("actual property not found: %v", mk)
				}
			} else {
				s.LogError("encountered unknown key: %v", mk)
			}
		}
	}
	return
}

func (s *CSettings) GetAlternativeButtonOrder() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkAlternativeButtonOrder); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetAlternativeSortArrows() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkAlternativeSortArrows); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetColorPalette() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkColorPalette); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetColorScheme() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkColorScheme); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetCursorBlink() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkCursorBlink); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetCursorBlinkTime() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkCursorBlinkTime); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetCursorBlinkTimeout() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkCursorBlinkTimeout); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetCursorThemeName() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkCursorThemeName); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetDndDragThreshold() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkDndDragThreshold); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetDoubleClickDistance() (value int) {
	var err error
	if value, err = s.GetIntProperty(PropertyCtkDoubleClickDistance); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetDoubleClickTime() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkDoubleClickTime); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetEnableAccels() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkEnableAccels); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetEnableMnemonics() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkEnableMnemonics); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetEnableTooltips() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkEnableTooltips); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetEntryPasswordHintTimeout() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkEntryPasswordHintTimeout); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetEntrySelectOnFocus() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkEntrySelectOnFocus); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetErrorBell() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkErrorBell); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetFallbackIconTheme() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkFallbackIconTheme); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetFileChooserBackend() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkFileChooserBackend); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetIconThemeName() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkIconThemeName); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetImModule() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkImModule); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetImPreeditStyle() (value interface{}) {
	var err error
	if value, err = s.GetStructProperty(PropertyCtkImPreeditStyle); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetImStatusStyle() (value interface{}) {
	var err error
	if value, err = s.GetStructProperty(PropertyCtkImStatusStyle); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetKeyThemeName() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkKeyThemeName); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetKeynavCursorOnly() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkKeynavCursorOnly); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetKeynavWrapAround() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkKeynavWrapAround); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetLabelSelectOnFocus() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkLabelSelectOnFocus); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetMenuBarAccel() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkMenuBarAccel); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetMenuBarPopupDelay() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkMenuBarPopupDelay); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetMenuImages() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkMenuImages); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetMenuPopdownDelay() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkMenuPopdownDelay); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetMenuPopupDelay() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkMenuPopupDelay); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetModules() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkModules); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetPrimaryButtonWarpsSlider() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkPrimaryButtonWarpsSlider); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetScrolledWindowPlacement() (value interface{}) {
	var err error
	if value, err = s.GetStructProperty(PropertyCtkScrolledWindowPlacement); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetShowInputMethodMenu() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkShowInputMethodMenu); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetShowUnicodeMenu() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkShowUnicodeMenu); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetThemeName() (value string) {
	var err error
	if value, err = s.GetStringProperty(PropertyCtkThemeName); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetTimeoutExpand() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkTimeoutExpand); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetTimeoutInitial() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkTimeoutInitial); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetTimeoutRepeat() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkTimeoutRepeat); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetToolbarStyle() (value interface{}) {
	var err error
	if value, err = s.GetStructProperty(PropertyCtkToolbarStyle); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetTooltipBrowseModeTimeout() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkTooltipBrowseModeTimeout); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetTooltipBrowseTimeout() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkTooltipBrowseTimeout); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetTooltipTimeout() (value time.Duration) {
	var err error
	if value, err = s.GetTimeProperty(PropertyCtkTooltipTimeout); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) GetTouchscreenMode() (value bool) {
	var err error
	if value, err = s.GetBoolProperty(PropertyCtkTouchscreenMode); err != nil {
		s.LogErr(err)
	}
	return
}

func (s *CSettings) SetCtkAlternativeButtonOrder(value bool) {
	if f := s.Emit(SignalSetCtkAlternativeButtonOrder, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkAlternativeButtonOrder, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkAlternativeSortArrows(value bool) {
	if f := s.Emit(SignalSetCtkAlternativeSortArrows, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkAlternativeSortArrows, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkColorPalette(value string) {
	if f := s.Emit(SignalSetCtkColorPalette, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkColorPalette, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkColorScheme(value string) {
	if f := s.Emit(SignalSetCtkColorScheme, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkColorScheme, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkCursorBlink(value bool) {
	if f := s.Emit(SignalSetCtkCursorBlink, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkCursorBlink, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkCursorBlinkTime(value time.Duration) {
	if f := s.Emit(SignalSetCtkCursorBlinkTime, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkCursorBlinkTime, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkCursorBlinkTimeout(value time.Duration) {
	if f := s.Emit(SignalSetCtkCursorBlinkTimeout, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkCursorBlinkTimeout, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkCursorThemeName(value string) {
	if f := s.Emit(SignalSetCtkCursorThemeName, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkCursorThemeName, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkDndDragThreshold(value time.Duration) {
	if f := s.Emit(SignalSetCtkDndDragThreshold, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkDndDragThreshold, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkDoubleClickDistance(value int) {
	if f := s.Emit(SignalSetCtkDoubleClickDistance, value); f == enums.EVENT_PASS {
		if err := s.SetIntProperty(PropertyCtkDoubleClickDistance, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkDoubleClickTime(value time.Duration) {
	if f := s.Emit(SignalSetCtkDoubleClickTime, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkDoubleClickTime, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkEnableAccels(value bool) {
	if f := s.Emit(SignalSetCtkEnableAccels, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkEnableAccels, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkEnableMnemonics(value bool) {
	if f := s.Emit(SignalSetCtkEnableMnemonics, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkEnableMnemonics, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkEnableTooltips(value bool) {
	if f := s.Emit(SignalSetCtkEnableTooltips, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkEnableTooltips, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkEntryPasswordHintTimeout(value time.Duration) {
	if f := s.Emit(SignalSetCtkEntryPasswordHintTimeout, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkEntryPasswordHintTimeout, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkEntrySelectOnFocus(value bool) {
	if f := s.Emit(SignalSetCtkEntrySelectOnFocus, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkEntrySelectOnFocus, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkErrorBell(value bool) {
	if f := s.Emit(SignalSetCtkErrorBell, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkErrorBell, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkFallbackIconTheme(value string) {
	if f := s.Emit(SignalSetCtkFallbackIconTheme, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkFallbackIconTheme, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkFileChooserBackend(value string) {
	if f := s.Emit(SignalSetCtkFileChooserBackend, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkFileChooserBackend, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkIconThemeName(value string) {
	if f := s.Emit(SignalSetCtkIconThemeName, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkIconThemeName, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkImModule(value string) {
	if f := s.Emit(SignalSetCtkImModule, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkImModule, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkImPreeditStyle(value interface{}) {
	if f := s.Emit(SignalSetCtkImPreeditStyle, value); f == enums.EVENT_PASS {
		if err := s.SetStructProperty(PropertyCtkImPreeditStyle, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkImStatusStyle(value interface{}) {
	if f := s.Emit(SignalSetCtkImStatusStyle, value); f == enums.EVENT_PASS {
		if err := s.SetStructProperty(PropertyCtkImStatusStyle, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkKeyThemeName(value string) {
	if f := s.Emit(SignalSetCtkKeyThemeName, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkKeyThemeName, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkKeynavCursorOnly(value bool) {
	if f := s.Emit(SignalSetCtkKeynavCursorOnly, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkKeynavCursorOnly, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkKeynavWrapAround(value bool) {
	if f := s.Emit(SignalSetCtkKeynavWrapAround, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkKeynavWrapAround, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkLabelSelectOnFocus(value bool) {
	if f := s.Emit(SignalSetCtkLabelSelectOnFocus, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkLabelSelectOnFocus, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkMenuBarAccel(value string) {
	if f := s.Emit(SignalSetCtkMenuBarAccel, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkMenuBarAccel, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkMenuBarPopupDelay(value time.Duration) {
	if f := s.Emit(SignalSetCtkMenuBarPopupDelay, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkMenuBarPopupDelay, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkMenuImages(value bool) {
	if f := s.Emit(SignalSetCtkMenuImages, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkMenuImages, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkMenuPopdownDelay(value time.Duration) {
	if f := s.Emit(SignalSetCtkMenuPopdownDelay, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkMenuPopdownDelay, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkMenuPopupDelay(value time.Duration) {
	if f := s.Emit(SignalSetCtkMenuPopupDelay, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkMenuPopupDelay, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkModules(value string) {
	if f := s.Emit(SignalSetCtkModules, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkModules, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkPrimaryButtonWarpsSlider(value bool) {
	if f := s.Emit(SignalSetCtkPrimaryButtonWarpsSlider, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkPrimaryButtonWarpsSlider, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkScrolledWindowPlacement(value interface{}) {
	if f := s.Emit(SignalSetCtkScrolledWindowPlacement, value); f == enums.EVENT_PASS {
		if err := s.SetStructProperty(PropertyCtkScrolledWindowPlacement, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkShowInputMethodMenu(value bool) {
	if f := s.Emit(SignalSetCtkShowInputMethodMenu, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkShowInputMethodMenu, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkShowUnicodeMenu(value bool) {
	if f := s.Emit(SignalSetCtkShowUnicodeMenu, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkShowUnicodeMenu, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkThemeName(value string) {
	if f := s.Emit(SignalSetCtkThemeName, value); f == enums.EVENT_PASS {
		if err := s.SetStringProperty(PropertyCtkThemeName, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkTimeoutExpand(value time.Duration) {
	if f := s.Emit(SignalSetCtkTimeoutExpand, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkTimeoutExpand, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkTimeoutInitial(value time.Duration) {
	if f := s.Emit(SignalSetCtkTimeoutInitial, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkTimeoutInitial, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkTimeoutRepeat(value time.Duration) {
	if f := s.Emit(SignalSetCtkTimeoutRepeat, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkTimeoutRepeat, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkToolbarStyle(value interface{}) {
	if f := s.Emit(SignalSetCtkToolbarStyle, value); f == enums.EVENT_PASS {
		if err := s.SetStructProperty(PropertyCtkToolbarStyle, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkTooltipBrowseModeTimeout(value time.Duration) {
	if f := s.Emit(SignalSetCtkTooltipBrowseModeTimeout, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkTooltipBrowseModeTimeout, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkTooltipBrowseTimeout(value time.Duration) {
	if f := s.Emit(SignalSetCtkTooltipBrowseTimeout, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkTooltipBrowseTimeout, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkTooltipTimeout(value time.Duration) {
	if f := s.Emit(SignalSetCtkTooltipTimeout, value); f == enums.EVENT_PASS {
		if err := s.SetTimeProperty(PropertyCtkTooltipTimeout, value); err != nil {
			s.LogErr(err)
		}
	}
}

func (s *CSettings) SetCtkTouchscreenMode(value bool) {
	if f := s.Emit(SignalSetCtkTouchscreenMode, value); f == enums.EVENT_PASS {
		if err := s.SetBoolProperty(PropertyCtkTouchscreenMode, value); err != nil {
			s.LogErr(err)
		}
	}
}

func ctkSettingsPropertyKeys() []cdk.Property {
	return []cdk.Property{
		PropertyCtkAlternativeButtonOrder,
		PropertyCtkAlternativeSortArrows,
		PropertyCtkColorPalette,
		PropertyCtkColorScheme,
		PropertyCtkCursorBlink,
		PropertyCtkCursorBlinkTime,
		PropertyCtkCursorBlinkTimeout,
		PropertyCtkCursorThemeName,
		PropertyCtkDndDragThreshold,
		PropertyCtkDoubleClickDistance,
		PropertyCtkDoubleClickTime,
		PropertyCtkEnableAccels,
		PropertyCtkEnableMnemonics,
		PropertyCtkEnableTooltips,
		PropertyCtkEntryPasswordHintTimeout,
		PropertyCtkEntrySelectOnFocus,
		PropertyCtkErrorBell,
		PropertyCtkFallbackIconTheme,
		PropertyCtkFileChooserBackend,
		PropertyCtkIconThemeName,
		PropertyCtkImModule,
		PropertyCtkImPreeditStyle,
		PropertyCtkImStatusStyle,
		PropertyCtkKeyThemeName,
		PropertyCtkKeynavCursorOnly,
		PropertyCtkKeynavWrapAround,
		PropertyCtkLabelSelectOnFocus,
		PropertyCtkMenuBarAccel,
		PropertyCtkMenuBarPopupDelay,
		PropertyCtkMenuImages,
		PropertyCtkMenuPopdownDelay,
		PropertyCtkMenuPopupDelay,
		PropertyCtkModules,
		PropertyCtkPrimaryButtonWarpsSlider,
		PropertyCtkScrolledWindowPlacement,
		PropertyCtkShowInputMethodMenu,
		PropertyCtkShowUnicodeMenu,
		PropertyCtkThemeName,
		PropertyCtkTimeoutExpand,
		PropertyCtkTimeoutInitial,
		PropertyCtkTimeoutRepeat,
		PropertyCtkToolbarStyle,
		PropertyCtkTooltipBrowseModeTimeout,
		PropertyCtkTooltipBrowseTimeout,
		PropertyCtkTooltipTimeout,
		PropertyCtkTouchscreenMode,
	}
}

// Whether buttons in dialogs should use the alternative button order.
// Flags: Read / Write
// Default value: FALSE
const PropertyCtkAlternativeButtonOrder cdk.Property = "ctk-alternative-button-order"

// Controls the direction of the sort indicators in sorted list and tree
// views. By default an arrow pointing down means the column is sorted in
// ascending order. When set to TRUE, this order will be inverted.
// Flags: Read / Write
// Default value: FALSE
const PropertyCtkAlternativeSortArrows cdk.Property = "ctk-alternative-sort-arrows"

// Palette to use in the color selector.
// Flags: Read / Write
// Default value: "black:white:gray50:red:purple:blue:light blue:green:yellow:orange:lavender:brown:goldenrod4:dodger blue:pink:light green:gray10:gray30:gray75:gray90"
const PropertyCtkColorPalette cdk.Property = "ctk-color-palette"

// A palette of named colors for use in themes. The format of the string is
// Color names must be acceptable as identifiers in the color specifications
// must be in the format accepted by ColorParse. Note that due to the
// way the color tables from different sources are merged, color
// specifications will be converted to hexadecimal form when getting this
// property. Starting with CTK 2.12, the entries can alternatively be
// separated by ';' instead of newlines:
// Flags: Read / Write
// Default value: ""
const PropertyCtkColorScheme cdk.Property = "ctk-color-scheme"

// Whether the cursor should blink. Also see the
// “ctk-cursor-blink-timeout” setting, which allows more flexible control
// over cursor blinking.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkCursorBlink cdk.Property = "ctk-cursor-blink"

// Length of the cursor blink cycle, in milliseconds.
// Flags: Read / Write
// Allowed values: >= 100
// Default value: 1200
const PropertyCtkCursorBlinkTime cdk.Property = "ctk-cursor-blink-time"

// Time after which the cursor stops blinking, in seconds. The timer is reset
// after each user interaction. Setting this to zero has the same effect as
// setting “ctk-cursor-blink” to FALSE.
// Flags: Read / Write
// Allowed values: >= 1
// Default value: 2147483647
const PropertyCtkCursorBlinkTimeout cdk.Property = "ctk-cursor-blink-timeout"

// Name of the cursor theme to use, or NULL to use the default theme.
// Flags: Read / Write
// Default value: NULL
const PropertyCtkCursorThemeName cdk.Property = "ctk-cursor-theme-name"

// Number of pixels the cursor can move before dragging.
// Flags: Read / Write
// Allowed values: >= 1
// Default value: 8
const PropertyCtkDndDragThreshold cdk.Property = "ctk-dnd-drag-threshold"

// Maximum distance allowed between two clicks for them to be considered a
// double click (in pixels).
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 5
const PropertyCtkDoubleClickDistance cdk.Property = "ctk-double-click-distance"

// Maximum time allowed between two clicks for them to be considered a double
// click (in milliseconds).
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 250
const PropertyCtkDoubleClickTime cdk.Property = "ctk-double-click-time"

// Whether menu items should have visible accelerators which can be
// activated.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkEnableAccels cdk.Property = "ctk-enable-accels"

// Whether labels and menu items should have visible mnemonics which can be
// activated.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkEnableMnemonics cdk.Property = "ctk-enable-mnemonics"

// Whether tooltips should be shown on widgets.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkEnableTooltips cdk.Property = "ctk-enable-tooltips"

// How long to show the last input character in hidden entries. This value is
// in milliseconds. 0 disables showing the last char. 600 is a good value for
// enabling it.
// Flags: Read / Write
// Default value: 0
const PropertyCtkEntryPasswordHintTimeout cdk.Property = "ctk-entry-password-hint-timeout"

// Whether to select the contents of an entry when it is focused.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkEntrySelectOnFocus cdk.Property = "ctk-entry-select-on-focus"

// When TRUE, keyboard navigation and other input-related errors will cause a
// beep. Since the error bell is implemented using WindowBeep, the
// windowing system may offer ways to configure the error bell in many ways,
// such as flashing the window or similar visual effects.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkErrorBell cdk.Property = "ctk-error-bell"

// Name of a icon theme to fall back to.
// Flags: Read / Write
// Default value: NULL
const PropertyCtkFallbackIconTheme cdk.Property = "ctk-fallback-icon-theme"

// Name of the FileChooser backend to use by default.
// Flags: Read / Write
// Default value: NULL
const PropertyCtkFileChooserBackend cdk.Property = "ctk-file-chooser-backend"

// Name of icon theme to use.
// Flags: Read / Write
// Default value: "hicolor"
const PropertyCtkIconThemeName cdk.Property = "ctk-icon-theme-name"

// Which IM (input method) module should be used by default. This is the
// input method that will be used if the user has not explicitly chosen
// another input method from the IM context menu. This also can be a
// colon-separated list of input methods, which CTK will try in turn until
// it finds one available on the system. See IMContext and see the
// “ctk-show-input-method-menu” property.
// Flags: Read / Write
// Default value: NULL
const PropertyCtkImModule cdk.Property = "ctk-im-module"

// How to draw the input method preedit string.
// Flags: Read / Write
// Default value: ctk_IM_PREEDIT_CALLBACK
const PropertyCtkImPreeditStyle cdk.Property = "ctk-im-preedit-style"

// How to draw the input method statusbar.
// Flags: Read / Write
// Default value: ctk_IM_STATUS_CALLBACK
const PropertyCtkImStatusStyle cdk.Property = "ctk-im-status-style"

// Name of key theme RC file to load.
// Flags: Read / Write
// Default value: NULL
const PropertyCtkKeyThemeName cdk.Property = "ctk-key-theme-name"

// When TRUE, keyboard navigation should be able to reach all widgets by
// using the cursor keys only. Tab, Shift etc. keys can't be expected to be
// present on the used input device.
// Flags: Read / Write
// Default value: FALSE
const PropertyCtkKeynavCursorOnly cdk.Property = "ctk-keynav-cursor-only"

// When TRUE, some widgets will wrap around when doing keyboard navigation,
// such as menus, menubars and notebooks.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkKeynavWrapAround cdk.Property = "ctk-keynav-wrap-around"

// Whether to select the contents of a selectable label when it is focused.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkLabelSelectOnFocus cdk.Property = "ctk-label-select-on-focus"

// Keybinding to activate the menu bar.
// Flags: Read / Write
// Default value: "F10"
const PropertyCtkMenuBarAccel cdk.Property = "ctk-menu-bar-accel"

// Delay before the submenus of a menu bar appear.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 0
const PropertyCtkMenuBarPopupDelay cdk.Property = "ctk-menu-bar-popup-delay"

// Whether images should be shown in menus.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkMenuImages cdk.Property = "ctk-menu-images"

// The time before hiding a submenu when the pointer is moving towards the
// submenu.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 1000
const PropertyCtkMenuPopdownDelay cdk.Property = "ctk-menu-popdown-delay"

// Minimum time the pointer must stay over a menu item before the submenu
// appear.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 225
const PropertyCtkMenuPopupDelay cdk.Property = "ctk-menu-popup-delay"

// List of currently active ctk modules.
// Flags: Read / Write
// Default value: NULL
const PropertyCtkModules cdk.Property = "ctk-modules"

// Whether a click in a Range trough should scroll to the click position
// or scroll by a single page in the respective direction.
// Flags: Read / Write
// Default value: FALSE
const PropertyCtkPrimaryButtonWarpsSlider cdk.Property = "ctk-primary-button-warps-slider"

// Where the contents of scrolled windows are located with respect to the
// scrollbars, if not overridden by the scrolled window's own placement.
// Flags: Read / Write
// Default value: ctk_CORNER_TOP_LEFT
const PropertyCtkScrolledWindowPlacement cdk.Property = "ctk-scrolled-window-placement"

// Whether the context menus of entries and text views should offer to change
// the input method.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkShowInputMethodMenu cdk.Property = "ctk-show-input-method-menu"

// Whether the context menus of entries and text views should offer to insert
// control characters.
// Flags: Read / Write
// Default value: TRUE
const PropertyCtkShowUnicodeMenu cdk.Property = "ctk-show-unicode-menu"

// Name of theme RC file to load.
// Flags: Read / Write
// Default value: "Raleigh"
const PropertyCtkThemeName cdk.Property = "ctk-theme-name"

// Expand value for timeouts, when a widget is expanding a new region.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 500
const PropertyCtkTimeoutExpand cdk.Property = "ctk-timeout-expand"

// Starting value for timeouts, when button is pressed.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 200
const PropertyCtkTimeoutInitial cdk.Property = "ctk-timeout-initial"

// Repeat value for timeouts, when button is pressed.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 20
const PropertyCtkTimeoutRepeat cdk.Property = "ctk-timeout-repeat"

// Whether default toolbars have text only, text and icons, icons only, etc.
// Flags: Read / Write
// Default value: ctk_TOOLBAR_BOTH
const PropertyCtkToolbarStyle cdk.Property = "ctk-toolbar-style"

// Amount of time, in milliseconds, after which the browse mode will be
// disabled. See “ctk-tooltip-browse-timeout” for more information about
// browse mode.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 500
const PropertyCtkTooltipBrowseModeTimeout cdk.Property = "ctk-tooltip-browse-mode-timeout"

// Controls the time after which tooltips will appear when browse mode is
// enabled, in milliseconds. Browse mode is enabled when the mouse pointer
// moves off an object where a tooltip was currently being displayed. If the
// mouse pointer hits another object before the browse mode timeout expires
// (see “ctk-tooltip-browse-mode-timeout”), it will take the amount of
// milliseconds specified by this setting to popup the tooltip for the new
// object.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 60
const PropertyCtkTooltipBrowseTimeout cdk.Property = "ctk-tooltip-browse-timeout"

// Time, in milliseconds, after which a tooltip could appear if the cursor is
// hovering on top of a widget.
// Flags: Read / Write
// Allowed values: >= 0
// Default value: 500
const PropertyCtkTooltipTimeout cdk.Property = "ctk-tooltip-timeout"

// When TRUE, there are no motion notify events delivered on this screen, and
// widgets can't use the pointer hovering them for any essential
// functionality.
// Flags: Read / Write
// Default value: FALSE
const PropertyCtkTouchscreenMode cdk.Property = "ctk-touchscreen-mode"

const SignalSetCtkAlternativeButtonOrder cdk.Signal = "ctk-alternative-button-order"
const SignalSetCtkAlternativeSortArrows cdk.Signal = "ctk-alternative-sort-arrows"
const SignalSetCtkColorPalette cdk.Signal = "ctk-color-palette"
const SignalSetCtkColorScheme cdk.Signal = "ctk-color-scheme"
const SignalSetCtkCursorBlink cdk.Signal = "ctk-cursor-blink"
const SignalSetCtkCursorBlinkTime cdk.Signal = "ctk-cursor-blink-time"
const SignalSetCtkCursorBlinkTimeout cdk.Signal = "ctk-cursor-blink-timeout"
const SignalSetCtkCursorThemeName cdk.Signal = "ctk-cursor-theme-name"
const SignalSetCtkDndDragThreshold cdk.Signal = "ctk-dnd-drag-threshold"
const SignalSetCtkDoubleClickDistance cdk.Signal = "ctk-double-click-distance"
const SignalSetCtkDoubleClickTime cdk.Signal = "ctk-double-click-time"
const SignalSetCtkEnableAccels cdk.Signal = "ctk-enable-accels"
const SignalSetCtkEnableMnemonics cdk.Signal = "ctk-enable-mnemonics"
const SignalSetCtkEnableTooltips cdk.Signal = "ctk-enable-tooltips"
const SignalSetCtkEntryPasswordHintTimeout cdk.Signal = "ctk-entry-password-hint-timeout"
const SignalSetCtkEntrySelectOnFocus cdk.Signal = "ctk-entry-select-on-focus"
const SignalSetCtkErrorBell cdk.Signal = "ctk-error-bell"
const SignalSetCtkFallbackIconTheme cdk.Signal = "ctk-fallback-icon-theme"
const SignalSetCtkFileChooserBackend cdk.Signal = "ctk-file-chooser-backend"
const SignalSetCtkIconThemeName cdk.Signal = "ctk-icon-theme-name"
const SignalSetCtkImModule cdk.Signal = "ctk-im-module"
const SignalSetCtkImPreeditStyle cdk.Signal = "ctk-im-preedit-style"
const SignalSetCtkImStatusStyle cdk.Signal = "ctk-im-status-style"
const SignalSetCtkKeyThemeName cdk.Signal = "ctk-key-theme-name"
const SignalSetCtkKeynavCursorOnly cdk.Signal = "ctk-keynav-cursor-only"
const SignalSetCtkKeynavWrapAround cdk.Signal = "ctk-keynav-wrap-around"
const SignalSetCtkLabelSelectOnFocus cdk.Signal = "ctk-label-select-on-focus"
const SignalSetCtkMenuBarAccel cdk.Signal = "ctk-menu-bar-accel"
const SignalSetCtkMenuBarPopupDelay cdk.Signal = "ctk-menu-bar-popup-delay"
const SignalSetCtkMenuImages cdk.Signal = "ctk-menu-images"
const SignalSetCtkMenuPopdownDelay cdk.Signal = "ctk-menu-popdown-delay"
const SignalSetCtkMenuPopupDelay cdk.Signal = "ctk-menu-popup-delay"
const SignalSetCtkModules cdk.Signal = "ctk-modules"
const SignalSetCtkPrimaryButtonWarpsSlider cdk.Signal = "ctk-primary-button-warps-slider"
const SignalSetCtkScrolledWindowPlacement cdk.Signal = "ctk-scrolled-window-placement"
const SignalSetCtkShowInputMethodMenu cdk.Signal = "ctk-show-input-method-menu"
const SignalSetCtkShowUnicodeMenu cdk.Signal = "ctk-show-unicode-menu"
const SignalSetCtkThemeName cdk.Signal = "ctk-theme-name"
const SignalSetCtkTimeoutExpand cdk.Signal = "ctk-timeout-expand"
const SignalSetCtkTimeoutInitial cdk.Signal = "ctk-timeout-initial"
const SignalSetCtkTimeoutRepeat cdk.Signal = "ctk-timeout-repeat"
const SignalSetCtkToolbarStyle cdk.Signal = "ctk-toolbar-style"
const SignalSetCtkTooltipBrowseModeTimeout cdk.Signal = "ctk-tooltip-browse-mode-timeout"
const SignalSetCtkTooltipBrowseTimeout cdk.Signal = "ctk-tooltip-browse-timeout"
const SignalSetCtkTooltipTimeout cdk.Signal = "ctk-tooltip-timeout"
const SignalSetCtkTouchscreenMode cdk.Signal = "ctk-touchscreen-mode"