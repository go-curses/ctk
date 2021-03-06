// Code generated by "bitmasker -output enums_bitmask.go -kebab -type AccelFlags,CalendarDisplayOptions,CellRendererState,ButtonAction,DebugFlag,DialogFlags,AttachOptions,StateType,FileFilterFlags,PrivateFlags,RBNodeColor,RcFlags,RecentFilterFlags,TextSearchFlags,TreeModelFlags,TreeViewFlags,UIManagerItemType,WidgetFlags,ParamFlags"; DO NOT EDIT.

package enums

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[ACCEL_VISIBLE-1]
	_ = x[ACCEL_LOCKED-2]
	_ = x[ACCEL_MASK-0]
}

const _AccelFlags_name = "maskvisiblelocked"

var _AccelFlags_index = [...]uint8{0, 4, 11, 17}

func (i AccelFlags) String() string {
	if i >= AccelFlags(len(_AccelFlags_index)-1) {
		return "AccelFlags(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AccelFlags_name[_AccelFlags_index[i]:_AccelFlags_index[i+1]]
}

// Has returns TRUE if the given flag is present in the bitmask
func (i AccelFlags) Has(m AccelFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i AccelFlags) Set(m AccelFlags) AccelFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i AccelFlags) Clear(m AccelFlags) AccelFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i AccelFlags) Toggle(m AccelFlags) AccelFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[CALENDAR_SHOW_HEADING-1]
	_ = x[CALENDAR_SHOW_DAY_NAMES-2]
	_ = x[CALENDAR_NO_MONTH_CHANGE-4]
	_ = x[CALENDAR_SHOW_WEEK_NUMBERS-8]
	_ = x[CALENDAR_WEEK_START_MONDAY-16]
	_ = x[CALENDAR_SHOW_DETAILS-32]
}

const (
	_CalendarDisplayOptions_name_0 = "show-headingshow-day-names"
	_CalendarDisplayOptions_name_1 = "no-month-change"
	_CalendarDisplayOptions_name_2 = "show-week-numbers"
	_CalendarDisplayOptions_name_3 = "week-start-monday"
	_CalendarDisplayOptions_name_4 = "show-details"
)

var (
	_CalendarDisplayOptions_index_0 = [...]uint8{0, 12, 26}
)

func (i CalendarDisplayOptions) String() (value string) {
	update := func(t CalendarDisplayOptions, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(CalendarDisplayOptions(1), _CalendarDisplayOptions_name_0[_CalendarDisplayOptions_index_0[0]:_CalendarDisplayOptions_index_0[0+1]])
	update(CalendarDisplayOptions(2), _CalendarDisplayOptions_name_0[_CalendarDisplayOptions_index_0[1]:_CalendarDisplayOptions_index_0[1+1]])
	update(CalendarDisplayOptions(4), _CalendarDisplayOptions_name_1)
	update(CalendarDisplayOptions(8), _CalendarDisplayOptions_name_2)
	update(CalendarDisplayOptions(16), _CalendarDisplayOptions_name_3)
	update(CalendarDisplayOptions(32), _CalendarDisplayOptions_name_4)
	if value == "" {
		return "CalendarDisplayOptions(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i CalendarDisplayOptions) Has(m CalendarDisplayOptions) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i CalendarDisplayOptions) Set(m CalendarDisplayOptions) CalendarDisplayOptions {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i CalendarDisplayOptions) Clear(m CalendarDisplayOptions) CalendarDisplayOptions {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i CalendarDisplayOptions) Toggle(m CalendarDisplayOptions) CalendarDisplayOptions {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[CELL_RENDERER_SELECTED-1]
	_ = x[CELL_RENDERER_PRELIT-2]
	_ = x[CELL_RENDERER_INSENSITIVE-4]
	_ = x[CELL_RENDERER_SORTED-8]
	_ = x[CELL_RENDERER_FOCUSED-16]
}

const (
	_CellRendererState_name_0 = "renderer-selectedrenderer-prelit"
	_CellRendererState_name_1 = "renderer-insensitive"
	_CellRendererState_name_2 = "renderer-sorted"
	_CellRendererState_name_3 = "renderer-focused"
)

var (
	_CellRendererState_index_0 = [...]uint8{0, 17, 32}
)

func (i CellRendererState) String() (value string) {
	update := func(t CellRendererState, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(CellRendererState(1), _CellRendererState_name_0[_CellRendererState_index_0[0]:_CellRendererState_index_0[0+1]])
	update(CellRendererState(2), _CellRendererState_name_0[_CellRendererState_index_0[1]:_CellRendererState_index_0[1+1]])
	update(CellRendererState(4), _CellRendererState_name_1)
	update(CellRendererState(8), _CellRendererState_name_2)
	update(CellRendererState(16), _CellRendererState_name_3)
	if value == "" {
		return "CellRendererState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i CellRendererState) Has(m CellRendererState) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i CellRendererState) Set(m CellRendererState) CellRendererState {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i CellRendererState) Clear(m CellRendererState) CellRendererState {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i CellRendererState) Toggle(m CellRendererState) CellRendererState {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[BUTTON_IGNORED-1]
	_ = x[BUTTON_SELECTS-2]
	_ = x[BUTTON_DRAGS-4]
	_ = x[BUTTON_EXPANDS-8]
}

const (
	_ButtonAction_name_0 = "ignoredselects"
	_ButtonAction_name_1 = "drags"
	_ButtonAction_name_2 = "expands"
)

var (
	_ButtonAction_index_0 = [...]uint8{0, 7, 14}
)

func (i ButtonAction) String() (value string) {
	update := func(t ButtonAction, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(ButtonAction(1), _ButtonAction_name_0[_ButtonAction_index_0[0]:_ButtonAction_index_0[0+1]])
	update(ButtonAction(2), _ButtonAction_name_0[_ButtonAction_index_0[1]:_ButtonAction_index_0[1+1]])
	update(ButtonAction(4), _ButtonAction_name_1)
	update(ButtonAction(8), _ButtonAction_name_2)
	if value == "" {
		return "ButtonAction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i ButtonAction) Has(m ButtonAction) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i ButtonAction) Set(m ButtonAction) ButtonAction {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i ButtonAction) Clear(m ButtonAction) ButtonAction {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i ButtonAction) Toggle(m ButtonAction) ButtonAction {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[DEBUG_MISC-1]
	_ = x[DEBUG_PLUGSOCKET-2]
	_ = x[DEBUG_TEXT-4]
	_ = x[DEBUG_TREE-8]
	_ = x[DEBUG_UPDATES-16]
	_ = x[DEBUG_KEYBINDINGS-32]
	_ = x[DEBUG_MULTIHEAD-64]
	_ = x[DEBUG_MODULES-128]
	_ = x[DEBUG_GEOMETRY-256]
	_ = x[DEBUG_ICONTHEME-512]
	_ = x[DEBUG_PRINTING-1024]
	_ = x[DEBUG_BUILDER-2048]
}

const _DebugFlag_name = "DEBUG_MISCDEBUG_PLUGSOCKETDEBUG_TEXTDEBUG_TREEDEBUG_UPDATESDEBUG_KEYBINDINGSDEBUG_MULTIHEADDEBUG_MODULESDEBUG_GEOMETRYDEBUG_ICONTHEMEDEBUG_PRINTINGDEBUG_BUILDER"

var _DebugFlag_map = map[DebugFlag]string{
	1:    _DebugFlag_name[0:10],
	2:    _DebugFlag_name[10:26],
	4:    _DebugFlag_name[26:36],
	8:    _DebugFlag_name[36:46],
	16:   _DebugFlag_name[46:59],
	32:   _DebugFlag_name[59:76],
	64:   _DebugFlag_name[76:91],
	128:  _DebugFlag_name[91:104],
	256:  _DebugFlag_name[104:118],
	512:  _DebugFlag_name[118:133],
	1024: _DebugFlag_name[133:147],
	2048: _DebugFlag_name[147:160],
}

func (i DebugFlag) String() string {
	if str, ok := _DebugFlag_map[i]; ok {
		return str
	}
	return "DebugFlag(" + strconv.FormatInt(int64(i), 10) + ")"
}

// Has returns TRUE if the given flag is present in the bitmask
func (i DebugFlag) Has(m DebugFlag) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i DebugFlag) Set(m DebugFlag) DebugFlag {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i DebugFlag) Clear(m DebugFlag) DebugFlag {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i DebugFlag) Toggle(m DebugFlag) DebugFlag {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[DialogModal-1]
	_ = x[DialogDestroyWithParent-2]
	_ = x[DialogNoSeparator-4]
}

const (
	_DialogFlags_name_0 = "modaldestroy-with-parent"
	_DialogFlags_name_1 = "no-separator"
)

var (
	_DialogFlags_index_0 = [...]uint8{0, 5, 24}
)

func (i DialogFlags) String() (value string) {
	update := func(t DialogFlags, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(DialogFlags(1), _DialogFlags_name_0[_DialogFlags_index_0[0]:_DialogFlags_index_0[0+1]])
	update(DialogFlags(2), _DialogFlags_name_0[_DialogFlags_index_0[1]:_DialogFlags_index_0[1+1]])
	update(DialogFlags(4), _DialogFlags_name_1)
	if value == "" {
		return "DialogFlags(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i DialogFlags) Has(m DialogFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i DialogFlags) Set(m DialogFlags) DialogFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i DialogFlags) Clear(m DialogFlags) DialogFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i DialogFlags) Toggle(m DialogFlags) DialogFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[EXPAND-1]
	_ = x[SHRINK-2]
	_ = x[FILL-4]
}

const (
	_AttachOptions_name_0 = "expandshrink"
	_AttachOptions_name_1 = "fill"
)

var (
	_AttachOptions_index_0 = [...]uint8{0, 6, 12}
)

func (i AttachOptions) String() (value string) {
	update := func(t AttachOptions, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(AttachOptions(1), _AttachOptions_name_0[_AttachOptions_index_0[0]:_AttachOptions_index_0[0+1]])
	update(AttachOptions(2), _AttachOptions_name_0[_AttachOptions_index_0[1]:_AttachOptions_index_0[1+1]])
	update(AttachOptions(4), _AttachOptions_name_1)
	if value == "" {
		return "AttachOptions(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i AttachOptions) Has(m AttachOptions) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i AttachOptions) Set(m AttachOptions) AttachOptions {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i AttachOptions) Clear(m AttachOptions) AttachOptions {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i AttachOptions) Toggle(m AttachOptions) AttachOptions {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[StateNone-0]
	_ = x[StateNormal-2]
	_ = x[StateActive-4]
	_ = x[StatePrelight-8]
	_ = x[StateSelected-16]
	_ = x[StateInsensitive-32]
}

const (
	_StateType_name_0 = "none"
	_StateType_name_1 = "normal"
	_StateType_name_2 = "active"
	_StateType_name_3 = "prelight"
	_StateType_name_4 = "selected"
	_StateType_name_5 = "insensitive"
)

func (i StateType) String() (value string) {
	update := func(t StateType, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(StateType(0), _StateType_name_0)
	update(StateType(2), _StateType_name_1)
	update(StateType(4), _StateType_name_2)
	update(StateType(8), _StateType_name_3)
	update(StateType(16), _StateType_name_4)
	update(StateType(32), _StateType_name_5)
	if value == "" {
		return "StateType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i StateType) Has(m StateType) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i StateType) Set(m StateType) StateType {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i StateType) Clear(m StateType) StateType {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i StateType) Toggle(m StateType) StateType {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[FILE_FILTER_FILENAME-1]
	_ = x[FILE_FILTER_URI-2]
	_ = x[FILE_FILTER_DISPLAY_NAME-4]
	_ = x[FILE_FILTER_MIME_TYPE-8]
}

const (
	_FileFilterFlags_name_0 = "filter-filenamefilter-uri"
	_FileFilterFlags_name_1 = "filter-display-name"
	_FileFilterFlags_name_2 = "filter-mime-type"
)

var (
	_FileFilterFlags_index_0 = [...]uint8{0, 15, 25}
)

func (i FileFilterFlags) String() (value string) {
	update := func(t FileFilterFlags, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(FileFilterFlags(1), _FileFilterFlags_name_0[_FileFilterFlags_index_0[0]:_FileFilterFlags_index_0[0+1]])
	update(FileFilterFlags(2), _FileFilterFlags_name_0[_FileFilterFlags_index_0[1]:_FileFilterFlags_index_0[1+1]])
	update(FileFilterFlags(4), _FileFilterFlags_name_1)
	update(FileFilterFlags(8), _FileFilterFlags_name_2)
	if value == "" {
		return "FileFilterFlags(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i FileFilterFlags) Has(m FileFilterFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i FileFilterFlags) Set(m FileFilterFlags) FileFilterFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i FileFilterFlags) Clear(m FileFilterFlags) FileFilterFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i FileFilterFlags) Toggle(m FileFilterFlags) FileFilterFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[PRIVATE_USER_STYLE-1]
	_ = x[PRIVATE_RESIZE_PENDING-4]
	_ = x[PRIVATE_HAS_POINTER-8]
	_ = x[PRIVATE_SHADOWED-16]
	_ = x[PRIVATE_HAS_SHAPE_MASK-32]
	_ = x[PRIVATE_IN_REPARENT-64]
	_ = x[PRIVATE_DIRECTION_SET-128]
	_ = x[PRIVATE_DIRECTION_LTR-256]
	_ = x[PRIVATE_ANCHORED-512]
	_ = x[PRIVATE_CHILD_VISIBLE-1024]
	_ = x[PRIVATE_REDRAW_ON_ALLOC-2048]
	_ = x[PRIVATE_ALLOC_NEEDED-4096]
	_ = x[PRIVATE_REQUEST_NEEDED-8192]
}

const _PrivateFlags_name = "PRIVATE_USER_STYLEPRIVATE_RESIZE_PENDINGPRIVATE_HAS_POINTERPRIVATE_SHADOWEDPRIVATE_HAS_SHAPE_MASKPRIVATE_IN_REPARENTPRIVATE_DIRECTION_SETPRIVATE_DIRECTION_LTRPRIVATE_ANCHOREDPRIVATE_CHILD_VISIBLEPRIVATE_REDRAW_ON_ALLOCPRIVATE_ALLOC_NEEDEDPRIVATE_REQUEST_NEEDED"

var _PrivateFlags_map = map[PrivateFlags]string{
	1:    _PrivateFlags_name[0:18],
	4:    _PrivateFlags_name[18:40],
	8:    _PrivateFlags_name[40:59],
	16:   _PrivateFlags_name[59:75],
	32:   _PrivateFlags_name[75:97],
	64:   _PrivateFlags_name[97:116],
	128:  _PrivateFlags_name[116:137],
	256:  _PrivateFlags_name[137:158],
	512:  _PrivateFlags_name[158:174],
	1024: _PrivateFlags_name[174:195],
	2048: _PrivateFlags_name[195:218],
	4096: _PrivateFlags_name[218:238],
	8192: _PrivateFlags_name[238:260],
}

func (i PrivateFlags) String() string {
	if str, ok := _PrivateFlags_map[i]; ok {
		return str
	}
	return "PrivateFlags(" + strconv.FormatInt(int64(i), 10) + ")"
}

// Has returns TRUE if the given flag is present in the bitmask
func (i PrivateFlags) Has(m PrivateFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i PrivateFlags) Set(m PrivateFlags) PrivateFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i PrivateFlags) Clear(m PrivateFlags) PrivateFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i PrivateFlags) Toggle(m PrivateFlags) PrivateFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[RBNODE_BLACK-1]
	_ = x[RBNODE_RED-2]
	_ = x[RBNODE_IS_PARENT-4]
	_ = x[RBNODE_IS_SELECTED-8]
	_ = x[RBNODE_IS_PRELIT-16]
	_ = x[RBNODE_IS_SEMI_COLLAPSED-32]
	_ = x[RBNODE_IS_SEMI_EXPANDED-64]
	_ = x[RBNODE_INVALID-128]
	_ = x[RBNODE_COLUMN_INVALID-256]
	_ = x[RBNODE_DESCENDANTS_INVALID-512]
	_ = x[RBNODE_NON_COLORS-4]
}

const (
	_RBNodeColor_name_0 = "node-blacknode-red"
	_RBNodeColor_name_1 = "node-is-parent"
	_RBNodeColor_name_2 = "node-is-selected"
	_RBNodeColor_name_3 = "node-is-prelit"
	_RBNodeColor_name_4 = "node-is-semi-collapsed"
	_RBNodeColor_name_5 = "node-is-semi-expanded"
	_RBNodeColor_name_6 = "node-invalid"
	_RBNodeColor_name_7 = "node-column-invalid"
	_RBNodeColor_name_8 = "node-descendants-invalid"
)

var (
	_RBNodeColor_index_0 = [...]uint8{0, 10, 18}
)

func (i RBNodeColor) String() (value string) {
	update := func(t RBNodeColor, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(RBNodeColor(1), _RBNodeColor_name_0[_RBNodeColor_index_0[0]:_RBNodeColor_index_0[0+1]])
	update(RBNodeColor(2), _RBNodeColor_name_0[_RBNodeColor_index_0[1]:_RBNodeColor_index_0[1+1]])
	update(RBNodeColor(4), _RBNodeColor_name_1)
	update(RBNodeColor(8), _RBNodeColor_name_2)
	update(RBNodeColor(16), _RBNodeColor_name_3)
	update(RBNodeColor(32), _RBNodeColor_name_4)
	update(RBNodeColor(64), _RBNodeColor_name_5)
	update(RBNodeColor(128), _RBNodeColor_name_6)
	update(RBNodeColor(256), _RBNodeColor_name_7)
	update(RBNodeColor(512), _RBNodeColor_name_8)
	if value == "" {
		return "RBNodeColor(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i RBNodeColor) Has(m RBNodeColor) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i RBNodeColor) Set(m RBNodeColor) RBNodeColor {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i RBNodeColor) Clear(m RBNodeColor) RBNodeColor {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i RBNodeColor) Toggle(m RBNodeColor) RBNodeColor {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[RC_FG-1]
	_ = x[RC_BG-2]
	_ = x[RC_TEXT-4]
	_ = x[RC_BASE-8]
}

const (
	_RcFlags_name_0 = "fgbg"
	_RcFlags_name_1 = "text"
	_RcFlags_name_2 = "base"
)

var (
	_RcFlags_index_0 = [...]uint8{0, 2, 4}
)

func (i RcFlags) String() (value string) {
	update := func(t RcFlags, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(RcFlags(1), _RcFlags_name_0[_RcFlags_index_0[0]:_RcFlags_index_0[0+1]])
	update(RcFlags(2), _RcFlags_name_0[_RcFlags_index_0[1]:_RcFlags_index_0[1+1]])
	update(RcFlags(4), _RcFlags_name_1)
	update(RcFlags(8), _RcFlags_name_2)
	if value == "" {
		return "RcFlags(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i RcFlags) Has(m RcFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i RcFlags) Set(m RcFlags) RcFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i RcFlags) Clear(m RcFlags) RcFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i RcFlags) Toggle(m RcFlags) RcFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[RECENT_FILTER_URI-1]
	_ = x[RECENT_FILTER_DISPLAY_NAME-2]
	_ = x[RECENT_FILTER_MIME_TYPE-4]
	_ = x[RECENT_FILTER_APPLICATION-8]
	_ = x[RECENT_FILTER_GROUP-16]
	_ = x[RECENT_FILTER_AGE-32]
}

const (
	_RecentFilterFlags_name_0 = "filter-urifilter-display-name"
	_RecentFilterFlags_name_1 = "filter-mime-type"
	_RecentFilterFlags_name_2 = "filter-application"
	_RecentFilterFlags_name_3 = "filter-group"
	_RecentFilterFlags_name_4 = "filter-age"
)

var (
	_RecentFilterFlags_index_0 = [...]uint8{0, 10, 29}
)

func (i RecentFilterFlags) String() (value string) {
	update := func(t RecentFilterFlags, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(RecentFilterFlags(1), _RecentFilterFlags_name_0[_RecentFilterFlags_index_0[0]:_RecentFilterFlags_index_0[0+1]])
	update(RecentFilterFlags(2), _RecentFilterFlags_name_0[_RecentFilterFlags_index_0[1]:_RecentFilterFlags_index_0[1+1]])
	update(RecentFilterFlags(4), _RecentFilterFlags_name_1)
	update(RecentFilterFlags(8), _RecentFilterFlags_name_2)
	update(RecentFilterFlags(16), _RecentFilterFlags_name_3)
	update(RecentFilterFlags(32), _RecentFilterFlags_name_4)
	if value == "" {
		return "RecentFilterFlags(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i RecentFilterFlags) Has(m RecentFilterFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i RecentFilterFlags) Set(m RecentFilterFlags) RecentFilterFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i RecentFilterFlags) Clear(m RecentFilterFlags) RecentFilterFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i RecentFilterFlags) Toggle(m RecentFilterFlags) RecentFilterFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[TEXT_SEARCH_VISIBLE_ONLY-1]
	_ = x[TEXT_SEARCH_TEXT_ONLY-2]
}

const _TextSearchFlags_name = "search-visible-onlysearch-text-only"

var _TextSearchFlags_index = [...]uint8{0, 19, 35}

func (i TextSearchFlags) String() string {
	i -= 1
	if i >= TextSearchFlags(len(_TextSearchFlags_index)-1) {
		return "TextSearchFlags(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TextSearchFlags_name[_TextSearchFlags_index[i]:_TextSearchFlags_index[i+1]]
}

// Has returns TRUE if the given flag is present in the bitmask
func (i TextSearchFlags) Has(m TextSearchFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i TextSearchFlags) Set(m TextSearchFlags) TextSearchFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i TextSearchFlags) Clear(m TextSearchFlags) TextSearchFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i TextSearchFlags) Toggle(m TextSearchFlags) TextSearchFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[TREE_MODEL_ITERS_PERSIST-1]
	_ = x[TREE_MODEL_LIST_ONLY-2]
}

const _TreeModelFlags_name = "model-iters-persistmodel-list-only"

var _TreeModelFlags_index = [...]uint8{0, 19, 34}

func (i TreeModelFlags) String() string {
	i -= 1
	if i >= TreeModelFlags(len(_TreeModelFlags_index)-1) {
		return "TreeModelFlags(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TreeModelFlags_name[_TreeModelFlags_index[i]:_TreeModelFlags_index[i+1]]
}

// Has returns TRUE if the given flag is present in the bitmask
func (i TreeModelFlags) Has(m TreeModelFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i TreeModelFlags) Set(m TreeModelFlags) TreeModelFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i TreeModelFlags) Clear(m TreeModelFlags) TreeModelFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i TreeModelFlags) Toggle(m TreeModelFlags) TreeModelFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[TREE_VIEW_IS_LIST-1]
	_ = x[TREE_VIEW_SHOW_EXPANDERS-2]
	_ = x[TREE_VIEW_IN_COLUMN_RESIZE-4]
	_ = x[TREE_VIEW_ARROW_PRELIT-8]
	_ = x[TREE_VIEW_HEADERS_VISIBLE-16]
	_ = x[TREE_VIEW_DRAW_KEYFOCUS-32]
	_ = x[TREE_VIEW_MODEL_SETUP-64]
	_ = x[TREE_VIEW_IN_COLUMN_DRAG-128]
}

const (
	_TreeViewFlags_name_0 = "view-is-listview-show-expanders"
	_TreeViewFlags_name_1 = "view-in-column-resize"
	_TreeViewFlags_name_2 = "view-arrow-prelit"
	_TreeViewFlags_name_3 = "view-headers-visible"
	_TreeViewFlags_name_4 = "view-draw-keyfocus"
	_TreeViewFlags_name_5 = "view-model-setup"
	_TreeViewFlags_name_6 = "view-in-column-drag"
)

var (
	_TreeViewFlags_index_0 = [...]uint8{0, 12, 31}
)

func (i TreeViewFlags) String() (value string) {
	update := func(t TreeViewFlags, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(TreeViewFlags(1), _TreeViewFlags_name_0[_TreeViewFlags_index_0[0]:_TreeViewFlags_index_0[0+1]])
	update(TreeViewFlags(2), _TreeViewFlags_name_0[_TreeViewFlags_index_0[1]:_TreeViewFlags_index_0[1+1]])
	update(TreeViewFlags(4), _TreeViewFlags_name_1)
	update(TreeViewFlags(8), _TreeViewFlags_name_2)
	update(TreeViewFlags(16), _TreeViewFlags_name_3)
	update(TreeViewFlags(32), _TreeViewFlags_name_4)
	update(TreeViewFlags(64), _TreeViewFlags_name_5)
	update(TreeViewFlags(128), _TreeViewFlags_name_6)
	if value == "" {
		return "TreeViewFlags(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i TreeViewFlags) Has(m TreeViewFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i TreeViewFlags) Set(m TreeViewFlags) TreeViewFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i TreeViewFlags) Clear(m TreeViewFlags) TreeViewFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i TreeViewFlags) Toggle(m TreeViewFlags) TreeViewFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[UI_MANAGER_AUTO-0]
	_ = x[UI_MANAGER_MENUBAR-1]
	_ = x[UI_MANAGER_MENU-4]
	_ = x[UI_MANAGER_TOOLBAR-8]
	_ = x[UI_MANAGER_PLACEHOLDER-16]
	_ = x[UI_MANAGER_POPUP-32]
	_ = x[UI_MANAGER_MENUITEM-64]
	_ = x[UI_MANAGER_TOOLITEM-128]
	_ = x[UI_MANAGER_SEPARATOR-256]
	_ = x[UI_MANAGER_ACCELERATOR-512]
	_ = x[UI_MANAGER_POPUP_WITH_ACCELS-1024]
}

const (
	_UIManagerItemType_name_0 = "manager-automanager-menubar"
	_UIManagerItemType_name_1 = "manager-menu"
	_UIManagerItemType_name_2 = "manager-toolbar"
	_UIManagerItemType_name_3 = "manager-placeholder"
	_UIManagerItemType_name_4 = "manager-popup"
	_UIManagerItemType_name_5 = "manager-menuitem"
	_UIManagerItemType_name_6 = "manager-toolitem"
	_UIManagerItemType_name_7 = "manager-separator"
	_UIManagerItemType_name_8 = "manager-accelerator"
	_UIManagerItemType_name_9 = "manager-popup-with-accels"
)

var (
	_UIManagerItemType_index_0 = [...]uint8{0, 12, 27}
)

func (i UIManagerItemType) String() (value string) {
	update := func(t UIManagerItemType, n string) {
		if i.Has(t) {
			if len(value) > 0 {
				value += " | "
			}
			value += n
		}
	}
	update(UIManagerItemType(0), _UIManagerItemType_name_0[_UIManagerItemType_index_0[0]:_UIManagerItemType_index_0[0+1]])
	update(UIManagerItemType(1), _UIManagerItemType_name_0[_UIManagerItemType_index_0[1]:_UIManagerItemType_index_0[1+1]])
	update(UIManagerItemType(4), _UIManagerItemType_name_1)
	update(UIManagerItemType(8), _UIManagerItemType_name_2)
	update(UIManagerItemType(16), _UIManagerItemType_name_3)
	update(UIManagerItemType(32), _UIManagerItemType_name_4)
	update(UIManagerItemType(64), _UIManagerItemType_name_5)
	update(UIManagerItemType(128), _UIManagerItemType_name_6)
	update(UIManagerItemType(256), _UIManagerItemType_name_7)
	update(UIManagerItemType(512), _UIManagerItemType_name_8)
	update(UIManagerItemType(1024), _UIManagerItemType_name_9)
	if value == "" {
		return "UIManagerItemType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return
}

// Has returns TRUE if the given flag is present in the bitmask
func (i UIManagerItemType) Has(m UIManagerItemType) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i UIManagerItemType) Set(m UIManagerItemType) UIManagerItemType {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i UIManagerItemType) Clear(m UIManagerItemType) UIManagerItemType {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i UIManagerItemType) Toggle(m UIManagerItemType) UIManagerItemType {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[NULL_WIDGET_FLAG-0]
	_ = x[TOPLEVEL-8]
	_ = x[NO_WINDOW-16]
	_ = x[REALIZED-32]
	_ = x[MAPPED-64]
	_ = x[VISIBLE-128]
	_ = x[SENSITIVE-256]
	_ = x[PARENT_SENSITIVE-512]
	_ = x[CAN_FOCUS-1024]
	_ = x[HAS_FOCUS-2048]
	_ = x[CAN_DEFAULT-4096]
	_ = x[HAS_DEFAULT-8192]
	_ = x[HAS_GRAB-16384]
	_ = x[RC_STYLE-32768]
	_ = x[COMPOSITE_CHILD-65536]
	_ = x[NO_REPARENT-131072]
	_ = x[APP_PAINTABLE-262144]
	_ = x[RECEIVES_DEFAULT-524288]
	_ = x[DOUBLE_BUFFERED-1048576]
	_ = x[NO_SHOW_ALL-2097152]
	_ = x[COMPOSITE_PARENT-4194304]
	_ = x[INVALID_WIDGET_FLAG-8388608]
}

const _WidgetFlags_name = "NULL_WIDGET_FLAGTOPLEVELNO_WINDOWREALIZEDMAPPEDVISIBLESENSITIVEPARENT_SENSITIVECAN_FOCUSHAS_FOCUSCAN_DEFAULTHAS_DEFAULTHAS_GRABRC_STYLECOMPOSITE_CHILDNO_REPARENTAPP_PAINTABLERECEIVES_DEFAULTDOUBLE_BUFFEREDNO_SHOW_ALLCOMPOSITE_PARENTINVALID_WIDGET_FLAG"

var _WidgetFlags_map = map[WidgetFlags]string{
	0:       _WidgetFlags_name[0:16],
	8:       _WidgetFlags_name[16:24],
	16:      _WidgetFlags_name[24:33],
	32:      _WidgetFlags_name[33:41],
	64:      _WidgetFlags_name[41:47],
	128:     _WidgetFlags_name[47:54],
	256:     _WidgetFlags_name[54:63],
	512:     _WidgetFlags_name[63:79],
	1024:    _WidgetFlags_name[79:88],
	2048:    _WidgetFlags_name[88:97],
	4096:    _WidgetFlags_name[97:108],
	8192:    _WidgetFlags_name[108:119],
	16384:   _WidgetFlags_name[119:127],
	32768:   _WidgetFlags_name[127:135],
	65536:   _WidgetFlags_name[135:150],
	131072:  _WidgetFlags_name[150:161],
	262144:  _WidgetFlags_name[161:174],
	524288:  _WidgetFlags_name[174:190],
	1048576: _WidgetFlags_name[190:205],
	2097152: _WidgetFlags_name[205:216],
	4194304: _WidgetFlags_name[216:232],
	8388608: _WidgetFlags_name[232:251],
}

func (i WidgetFlags) String() string {
	if str, ok := _WidgetFlags_map[i]; ok {
		return str
	}
	return "WidgetFlags(" + strconv.FormatInt(int64(i), 10) + ")"
}

// Has returns TRUE if the given flag is present in the bitmask
func (i WidgetFlags) Has(m WidgetFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i WidgetFlags) Set(m WidgetFlags) WidgetFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i WidgetFlags) Clear(m WidgetFlags) WidgetFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i WidgetFlags) Toggle(m WidgetFlags) WidgetFlags {
	return i ^ m
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the bitmasker command to generate them again.
	var x [1]struct{}
	_ = x[PARAM_READABLE-1]
	_ = x[PARAM_WRITABLE-2]
	_ = x[PARAM_READWRITE-3]
	_ = x[PARAM_CONSTRUCT-3]
	_ = x[PARAM_CONSTRUCT_ONLY-3]
	_ = x[PARAM_LAX_VALIDATION-3]
	_ = x[PARAM_STATIC_NAME-3]
	_ = x[PARAM_PRIVATE-3]
	_ = x[PARAM_STATIC_NICK-3]
	_ = x[PARAM_STATIC_BLURB-3]
	_ = x[PARAM_EXPLICIT_NOTIFY-3]
}

const _ParamFlags_name = "readablewritablereadwrite"

var _ParamFlags_index = [...]uint8{0, 8, 16, 25}

func (i ParamFlags) String() string {
	i -= 1
	if i >= ParamFlags(len(_ParamFlags_index)-1) {
		return "ParamFlags(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ParamFlags_name[_ParamFlags_index[i]:_ParamFlags_index[i+1]]
}

// Has returns TRUE if the given flag is present in the bitmask
func (i ParamFlags) Has(m ParamFlags) bool {
	if i == m {
		return true
	}
	return i&m != 0
}

// Set returns the bitmask with the given flag set
func (i ParamFlags) Set(m ParamFlags) ParamFlags {
	return i | m
}

// Clear returns the bitmask with the given flag removed
func (i ParamFlags) Clear(m ParamFlags) ParamFlags {
	return i &^ m
}

// Toggle returns the bitmask with the given flag toggled
func (i ParamFlags) Toggle(m ParamFlags) ParamFlags {
	return i ^ m
}
