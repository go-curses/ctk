package ctk

import (
	"fmt"
	"strconv"
	"strings"

	cbits "github.com/go-curses/cdk/lib/bits"
)

type EnumFromString interface {
	FromString(value string) (enum interface{}, err error)
}

/* Accel flags */
type AccelFlags uint64

const (
	ACCEL_VISIBLE AccelFlags = 1 << 0
	ACCEL_LOCKED  AccelFlags = 1 << iota
	ACCEL_MASK    AccelFlags = 0
)

/* Assistant page type */
type AssistantPageType uint64

const (
	ASSISTANT_PAGE_CONTENT AssistantPageType = iota
	ASSISTANT_PAGE_INTRO
	ASSISTANT_PAGE_CONFIRM
	ASSISTANT_PAGE_SUMMARY
	ASSISTANT_PAGE_PROGRESS
)

/* Builder error */
type BuilderError uint64

const (
	BUILDER_ERROR_INVALID_TYPE_FUNCTION BuilderError = iota
	BUILDER_ERROR_UNHANDLED_TAG
	BUILDER_ERROR_MISSINC_ATTRIBUTE
	BUILDER_ERROR_INVALID_ATTRIBUTE
	BUILDER_ERROR_INVALID_TAG
	BUILDER_ERROR_MISSINC_PROPERTY_VALUE
	BUILDER_ERROR_INVALID_VALUE
	BUILDER_ERROR_VERSION_MISMATCH
	BUILDER_ERROR_DUPLICATE_ID
)

/* Calendar display options */
type CalendarDisplayOptions uint64

const (
	CALENDAR_SHOW_HEADING   CalendarDisplayOptions = 1 << 0
	CALENDAR_SHOW_DAY_NAMES CalendarDisplayOptions = 1 << iota
	CALENDAR_NO_MONTH_CHANGE
	CALENDAR_SHOW_WEEK_NUMBERS
	CALENDAR_WEEK_START_MONDAY
	CALENDAR_SHOW_DETAILS
)

/* Cell renderer state */
type CellRendererState uint64

const (
	CELL_RENDERER_SELECTED CellRendererState = 1 << 0
	CELL_RENDERER_PRELIT   CellRendererState = 1 << iota
	CELL_RENDERER_INSENSITIVE
	CELL_RENDERER_SORTED
	CELL_RENDERER_FOCUSED
)

/* Cell renderer mode */
type CellRendererMode uint64

const (
	CELL_RENDERER_MODE_INERT CellRendererMode = iota
	CELL_RENDERER_MODE_ACTIVATABLE
	CELL_RENDERER_MODE_EDITABLE
)

/* Cell renderer accel mode */
type CellRendererAccelMode uint64

const (
	CELL_RENDERER_ACCEL_MODE_CTK CellRendererAccelMode = iota
	CELL_RENDERER_ACCEL_MODE_OTHER
)

/* Cell type */
type CellType uint64

const (
	CELL_EMPTY CellType = iota
	CELL_TEXT
	CELL_PIXMAP
	CELL_PIXTEXT
	CELL_WIDGET
)

/* List drag pos */
type CListDragPos uint64

const (
	CLIST_DRAC_NONE CListDragPos = iota
	CLIST_DRAC_BEFORE
	CLIST_DRAC_INTO
	CLIST_DRAC_AFTER
)

/* Button action */
type ButtonAction uint64

const (
	BUTTON_IGNORED ButtonAction = 0
	BUTTON_SELECTS ButtonAction = 1 << 0
	BUTTON_DRAGS   ButtonAction = 1 << iota
	BUTTON_EXPANDS
)

/* Tree pos */
type CTreePos uint64

const (
	CTREE_POS_BEFORE CTreePos = iota
	CTREE_POS_AS_CHILD
	CTREE_POS_AFTER
)

/* Tree line style */
type CTreeLineStyle uint64

const (
	CTREE_LINES_NONE CTreeLineStyle = iota
	CTREE_LINES_SOLID
	CTREE_LINES_DOTTED
	CTREE_LINES_TABBED
)

/* Tree expander style */
type CTreeExpanderStyle uint64

const (
	CTREE_EXPANDER_NONE CTreeExpanderStyle = iota
	CTREE_EXPANDER_SQUARE
	CTREE_EXPANDER_TRIANGLE
	CTREE_EXPANDER_CIRCULAR
)

/* Tree expansion type */
type CTreeExpansionType uint64

const (
	CTREE_EXPANSION_EXPAND CTreeExpansionType = iota
	CTREE_EXPANSION_EXPAND_RECURSIVE
	CTREE_EXPANSION_COLLAPSE
	CTREE_EXPANSION_COLLAPSE_RECURSIVE
	CTREE_EXPANSION_TOGGLE
	CTREE_EXPANSION_TOGGLE_RECURSIVE
)

/* Debug flag */
type DebugFlag uint64

const (
	DEBUG_MISC       DebugFlag = 1 << 0
	DEBUG_PLUGSOCKET DebugFlag = 1 << iota
	DEBUG_TEXT
	DEBUG_TREE
	DEBUG_UPDATES
	DEBUG_KEYBINDINGS
	DEBUG_MULTIHEAD
	DEBUG_MODULES
	DEBUG_GEOMETRY
	DEBUG_ICONTHEME
	DEBUG_PRINTING
	DEBUG_BUILDER
)

/* Dialog flags */
type DialogFlags uint64

const (
	DialogModal             DialogFlags = 1 << 0
	DialogDestroyWithParent DialogFlags = 1 << iota
	DialogNoSeparator
)

/* Response type */
type ResponseType int

const (
	ResponseNone        ResponseType = -1
	ResponseReject      ResponseType = -2
	ResponseAccept      ResponseType = -3
	ResponseDeleteEvent ResponseType = -4
	ResponseOk          ResponseType = -5
	ResponseCancel      ResponseType = -6
	ResponseClose       ResponseType = -7
	ResponseYes         ResponseType = -8
	ResponseNo          ResponseType = -9
	ResponseApply       ResponseType = -10
	ResponseHelp        ResponseType = -11
)

var (
	responseTypes = map[ResponseType]string{
		ResponseNone:        "none",
		ResponseReject:      "reject",
		ResponseAccept:      "accept",
		ResponseDeleteEvent: "delete-event",
		ResponseOk:          "ok",
		ResponseCancel:      "cancel",
		ResponseClose:       "close",
		ResponseYes:         "yes",
		ResponseNo:          "no",
		ResponseApply:       "apply",
		ResponseHelp:        "help",
	}
	responseNames = map[string]ResponseType{
		"none":         ResponseNone,
		"reject":       ResponseReject,
		"accept":       ResponseAccept,
		"delete-event": ResponseDeleteEvent,
		"ok":           ResponseOk,
		"cancel":       ResponseCancel,
		"close":        ResponseClose,
		"yes":          ResponseYes,
		"no":           ResponseNo,
		"apply":        ResponseApply,
		"help":         ResponseHelp,
	}
)

func (r ResponseType) String() string {
	if v, ok := responseTypes[r]; ok {
		return v
	}
	return strconv.Itoa(int(r))
}

func ResponseTypeFromName(name string) ResponseType {
	name = strings.ToLower(name)
	if v, ok := responseNames[name]; ok {
		return v
	}
	return ResponseNone
}

/* Entry icon position */
type EntryIconPosition uint64

const (
	ENTRY_ICON_PRIMARY EntryIconPosition = iota
	ENTRY_ICON_SECONDARY
)

/* Anchor type */
type AnchorType uint64

const (
	ANCHOR_CENTER AnchorType = iota
	ANCHOR_NORTH
	ANCHOR_NORTH_WEST
	ANCHOR_NORTH_EAST
	ANCHOR_SOUTH
	ANCHOR_SOUTH_WEST
	ANCHOR_SOUTH_EAST
	ANCHOR_WEST
	ANCHOR_EAST
	ANCHOR_N  = ANCHOR_NORTH
	ANCHOR_NW = ANCHOR_NORTH_WEST
	ANCHOR_NE = ANCHOR_NORTH_EAST
	ANCHOR_S  = ANCHOR_SOUTH
	ANCHOR_SW = ANCHOR_SOUTH_WEST
	ANCHOR_SE = ANCHOR_SOUTH_EAST
	ANCHOR_W  = ANCHOR_WEST
	ANCHOR_E  = ANCHOR_EAST
)

/* Arrow placement */
type ArrowPlacement uint64

const (
	ARROWS_BOTH ArrowPlacement = iota
	ARROWS_START
	ARROWS_END
)

/* Arrow type */
type ArrowType uint64

const (
	ArrowUp ArrowType = iota
	ArrowDown
	ArrowLeft
	ArrowRight
	ArrowNone
)

/* Attach options */
type AttachOptions uint64

const (
	EXPAND AttachOptions = 1 << 0
	SHRINK AttachOptions = 1 << iota
	FILL
)

/* Button box style */
type ButtonBoxStyle uint64

const (
	BUTTONBOX_DEFAULT_STYLE ButtonBoxStyle = iota
	// Spread the buttons evenly and centered away from the edges
	BUTTONBOX_SPREAD
	// Spread the buttons evenly and centered from edge to edge
	BUTTONBOX_EDGE
	// Group buttons at the start
	BUTTONBOX_START
	// Group buttons at the end
	BUTTONBOX_END
	// Group buttons at the center
	BUTTONBOX_CENTER
	// Buttons are expanded to evenly consume all space available
	BUTTONBOX_EXPAND
)

/* Curve type */
type CurveType uint64

const (
	CURVE_TYPE_LINEAR CurveType = iota
	CURVE_TYPE_SPLINE
	CURVE_TYPE_FREE
)

/* Delete type */
type DeleteType uint64

const (
	DELETE_CHARS DeleteType = iota
	DELETE_WORD_ENDS
	DELETE_WORDS
	DELETE_DISPLAY_LINES
	DELETE_DISPLAY_LINE_ENDS
	DELETE_PARAGRAPH_ENDS
	DELETE_PARAGRAPHS
	DELETE_WHITESPACE
)

/* Direction type */
type DirectionType uint64

const (
	DIR_TAB_FORWARD DirectionType = iota
	DIR_TAB_BACKWARD
	DIR_UP
	DIR_DOWN
	DIR_LEFT
	DIR_RIGHT
)

/* Expander style */
type ExpanderStyle uint64

const (
	EXPANDER_COLLAPSED ExpanderStyle = iota
	EXPANDER_SEMI_COLLAPSED
	EXPANDER_SEMI_EXPANDED
	EXPANDER_EXPANDED
)

/* Icon size */
type IconSize uint64

const (
	ICON_SIZE_INVALID IconSize = iota
	ICON_SIZE_MENU
	ICON_SIZE_SMALL_TOOLBAR
	ICON_SIZE_LARGE_TOOLBAR
	ICON_SIZE_BUTTON
	ICON_SIZE_DND
	ICON_SIZE_DIALOG
)

/* Sensitivity type */
type SensitivityType uint64

const (
	SensitivityAuto SensitivityType = iota
	SensitivityOn
	SensitivityOff
)

/* Side type */
type SideType uint64

const (
	SIDE_TOP SideType = iota
	SIDE_BOTTOM
	SIDE_LEFT
	SIDE_RIGHT
)

/* Text direction */
type TextDirection uint64

const (
	TextDirNone TextDirection = iota
	TextDirLtr
	TextDirRtl
)

/* Match type */
type MatchType uint64

const (
	MATCH_ALL MatchType = iota
	MATCH_ALL_TAIL
	MATCH_HEAD
	MATCH_TAIL
	MATCH_EXACT
	MATCH_LAST
)

/* Menu direction type */
type MenuDirectionType uint64

const (
	MENU_DIR_PARENT MenuDirectionType = iota
	MENU_DIR_CHILD
	MENU_DIR_NEXT
	MENU_DIR_PREV
)

/* Message type */
type MessageType uint64

const (
	MESSAGE_INFO MessageType = iota
	MESSAGE_WARNING
	MESSAGE_QUESTION
	MESSAGE_ERROR
	MESSAGE_OTHER
)

/* Metric type */
type MetricType uint64

const (
	PIXELS MetricType = iota
	INCHES
	CENTIMETERS
)

type ModifierType uint64

const (
	NullModMask ModifierType = 0
	ShiftMask   ModifierType = 1 << iota
	LockMask
	ControlMask
	Mod1Mask
	Mod2Mask
	Mod3Mask
	Mod4Mask
	Mod5Mask
	Button1Mask
	Button2Mask
	Button3Mask
	Button4Mask
	Button5Mask
	SuperMask
	HyperMask
	MetaMask
	ReleaseMask
	ModifierMask
)

func (m ModifierType) HasBit(b ModifierType) bool {
	return cbits.Has(uint64(m), uint64(b))
}

func (m ModifierType) String() string {
	v := ""
	if m.HasBit(SuperMask) || m.HasBit(MetaMask) {
		v += "<Super>"
	}
	if m.HasBit(ControlMask) {
		v += "<Control>"
	}
	if m.HasBit(Mod1Mask) {
		v += "<Mod1>"
	}
	if m.HasBit(Mod2Mask) {
		v += "<Mod2>"
	}
	if m.HasBit(Mod3Mask) {
		v += "<Mod3>"
	}
	if m.HasBit(Mod4Mask) {
		v += "<Mod4>"
	}
	if m.HasBit(Mod5Mask) {
		v += "<Mod5>"
	}
	if m.HasBit(ShiftMask) || m.HasBit(LockMask) {
		v += "<Shift>"
	}
	if m.HasBit(Button1Mask) {
		v += "button1"
	}
	if m.HasBit(Button2Mask) {
		v += "button2"
	}
	if m.HasBit(Button3Mask) {
		v += "button3"
	}
	if m.HasBit(Button4Mask) {
		v += "button4"
	}
	if m.HasBit(Button5Mask) {
		v += "button5"
	}
	return v
}

/* Movement step */
type MovementStep uint64

const (
	MOVEMENT_LOGICAL_POSITIONS MovementStep = iota
	MOVEMENT_VISUAL_POSITIONS
	MOVEMENT_WORDS
	MOVEMENT_DISPLAY_LINES
	MOVEMENT_DISPLAY_LINE_ENDS
	MOVEMENT_PARAGRAPHS
	MOVEMENT_PARAGRAPH_ENDS
	MOVEMENT_PAGES
	MOVEMENT_BUFFER_ENDS
	MOVEMENT_HORIZONTAL_PAGES
)

/* Scroll step */
type ScrollStep uint64

const (
	SCROLL_STEPS ScrollStep = iota
	SCROLL_PAGES
	SCROLL_ENDS
	SCROLL_HORIZONTAL_STEPS
	SCROLL_HORIZONTAL_PAGES
	SCROLL_HORIZONTAL_ENDS
)

/* Corner type */
type CornerType uint64

const (
	CornerTopLeft CornerType = iota
	CornerBottomLeft
	CornerTopRight
	CornerBottomRight
)

/* packing type */
type PackType uint64

const (
	PackStart PackType = iota
	PackEnd
)

/* layout style */
type LayoutStyle uint64

const (
	LayoutStart LayoutStyle = iota
	LayoutEnd
)

func (l LayoutStyle) FromString(value string) (enum interface{}, err error) {
	switch strings.ToLower(value) {
	case "start":
		enum = LayoutStart
	case "end":
		enum = LayoutEnd
	default:
		err = fmt.Errorf("unknown value for LayoutStyle.FromString(%v)", value)
	}
	return
}

/* Path priority type */
type PathPriorityType uint64

const (
	PATH_PRIO_LOWEST      PathPriorityType = 0
	PATH_PRIO_CTK         PathPriorityType = 4
	PATH_PRIO_APPLICATION PathPriorityType = 8
	PATH_PRIO_THEME       PathPriorityType = 10
	PATH_PRIO_RC          PathPriorityType = 12
	PATH_PRIO_HIGHEST     PathPriorityType = 15
)

/* Path type */
type PathType uint64

const (
	PATH_WIDGET PathType = iota
	PATH_WIDGET_CLASS
	PATH_CLASS
)

/* Policy type */
type PolicyType uint64

const (
	PolicyAlways PolicyType = iota
	PolicyAutomatic
	PolicyNever
)

func (p PolicyType) FromString(value string) (enum interface{}, err error) {
	switch strings.ToLower(value) {
	case "automatic":
		return PolicyAutomatic, nil
	case "always":
		return PolicyAlways, nil
	case "never":
		return PolicyNever, nil
	}
	return nil, fmt.Errorf("unknown value for WindowType.FromString(%v)", value)
}

/* Position type */
type PositionType uint64

const (
	POS_LEFT PositionType = iota
	POS_RIGHT
	POS_TOP
	POS_BOTTOM
)

/* Relief style */
type ReliefStyle uint64

const (
	RELIEF_NORMAL ReliefStyle = iota
	RELIEF_HALF
	RELIEF_NONE
)

/* Scroll type */
type ScrollType uint64

const (
	SCROLL_NONE ScrollType = iota
	SCROLL_JUMP
	SCROLL_STEP_BACKWARD
	SCROLL_STEP_FORWARD
	SCROLL_PAGE_BACKWARD
	SCROLL_PAGE_FORWARD
	SCROLL_STEP_UP
	SCROLL_STEP_DOWN
	SCROLL_PAGE_UP
	SCROLL_PAGE_DOWN
	SCROLL_STEP_LEFT
	SCROLL_STEP_RIGHT
	SCROLL_PAGE_LEFT
	SCROLL_PAGE_RIGHT
	SCROLL_START
	SCROLL_END
)

/* Selection mode */
type SelectionMode uint64

const (
	SELECTION_NONE SelectionMode = iota
	SELECTION_SINGLE
	SELECTION_BROWSE
	SELECTION_MULTIPLE
	SELECTION_EXTENDED SelectionMode = SelectionMode(SELECTION_MULTIPLE)
)

/* Shadow type */
type ShadowType uint64

const (
	SHADOW_NONE ShadowType = iota
	SHADOW_IN
	SHADOW_OUT
	SHADOW_ETCHED_IN
	SHADOW_ETCHED_OUT
)

type StateType uint64

const (
	StateNormal StateType = 1 << iota
	StateActive
	StatePrelight
	StateSelected
	StateInsensitive
)

// StateTypeFromString returns the StateType equivalent for the given named
// string. If the name given is not a valid name, returns StateNormal.
func StateTypeFromString(name string) (state StateType) {
	switch strings.ToLower(name) {
	case "active":
		return StateActive
	case "prelight":
		return StatePrelight
	case "selected":
		return StateSelected
	case "insensitive":
		return StateInsensitive
	case "normal":
		fallthrough
	default:
		return StateNormal
	}
}

func (s StateType) HasBit(state StateType) bool {
	return cbits.Has(uint64(s), uint64(state))
}

func (s StateType) String() (label string) {
	label = ""
	update := func(state StateType, name string) {
		if s.HasBit(state) {
			if len(label) > 0 {
				label += " | "
			}
			label += name
		}
	}
	update(StateNormal, "normal")
	update(StateActive, "active")
	update(StatePrelight, "prelight")
	update(StateInsensitive, "insensitive")
	update(StateSelected, "selected")
	if label == "" {
		label = "unknown"
	}
	return
}

/* Submenu direction */
type SubmenuDirection uint64

const (
	DIRECTION_LEFT SubmenuDirection = iota
	DIRECTION_RIGHT
)

/* Submenu placement */
type SubmenuPlacement uint64

const (
	TOP_BOTTOM SubmenuPlacement = iota
	LEFT_RIGHT
)

/* Toolbar style */
type ToolbarStyle uint64

const (
	TOOLBAR_ICONS ToolbarStyle = iota
	TOOLBAR_TEXT
	TOOLBAR_BOTH
	TOOLBAR_BOTH_HORIZ
)

/* Update type */
type UpdateType uint64

const (
	UpdateContinuous UpdateType = iota
	UpdateDiscontinuous
	UpdateDelayed
)

func (t UpdateType) FromString(value string) (enum interface{}, err error) {
	switch strings.ToLower(value) {
	case "continuous":
		return UpdateContinuous, nil
	case "discontinuous":
		return UpdateDiscontinuous, nil
	case "delayed":
		return UpdateDelayed, nil
	}
	return nil, fmt.Errorf("unknown value for WindowType.FromString(%v)", value)
}

/* Visibility */
type Visibility uint64

const (
	VISIBILITY_NONE Visibility = iota
	VISIBILITY_PARTIAL
	VISIBILITY_FULL
)

/* Window type */
type WindowType uint64

const (
	WindowTopLevel WindowType = iota
	WindowPopup
)

func (t WindowType) FromString(value string) (enum interface{}, err error) {
	switch strings.ToLower(value) {
	case "top-level", "toplevel":
		return WindowTopLevel, nil
	case "popup":
		return WindowPopup, nil
	}
	return nil, fmt.Errorf("unknown value for WindowType.FromString(%v)", value)
}

type WindowTypeHint uint64

const (
	// A normal toplevel window.
	WindowTypeHintNormal WindowTypeHint = iota
	// A dialog window.
	WindowTypeHintDialog
	// A window used to implement a menu.
	WindowTypeHintMenu
	// A window used to implement a toolbar.
	WindowTypeHintToolbar
	// A window used to implement a splash screen
	WindowTypeHintSplashscreen
	//
	WindowTypeHintUtility
	// A window used to implement a docking bar.
	WindowTypeHintDock
	// A window used to implement a desktop.
	WindowTypeHintDesktop
	// A menu that belongs to a menubar.
	WindowTypeHintDropdownMenu
	// A menu that does not belong to a menubar, e.g. a context menu.
	WindowTypeHintPopupMenu
	// A tooltip.
	WindowTypeHintTooltip
	// A notification - typically a "bubble" that belongs to a status icon.
	WindowTypeHintNotification
	// A popup from a combo box.
	WindowTypeHintCombo
	// A window that is used to implement a DND cursor.
	WindowTypeHintDND
)

func (t WindowTypeHint) FromString(value string) (enum interface{}, err error) {
	switch strings.ToLower(value) {
	case "normal":
		return WindowTypeHintNormal, nil
	case "dialog":
		return WindowTypeHintDialog, nil
	case "menu":
		return WindowTypeHintMenu, nil
	case "toolbar":
		return WindowTypeHintToolbar, nil
	case "splashscreen":
		return WindowTypeHintSplashscreen, nil
	case "utility":
		return WindowTypeHintUtility, nil
	case "dock":
		return WindowTypeHintDock, nil
	case "desktop":
		return WindowTypeHintDesktop, nil
	case "dropdown-menu":
		return WindowTypeHintDropdownMenu, nil
	case "popup-menu":
		return WindowTypeHintPopupMenu, nil
	case "tooltip":
		return WindowTypeHintTooltip, nil
	case "notification":
		return WindowTypeHintNotification, nil
	case "combo":
		return WindowTypeHintCombo, nil
	case "dnd":
		return WindowTypeHintDND, nil
	}
	return nil, fmt.Errorf("unknown value for WindowTypeHint.FromString(%v)", value)
}

type WindowEdge uint64

const (
	WindowEdgeNone WindowEdge = iota
	WindowEdgeNorthWest
	WindowEdgeNorth
	WindowEdgeNorthEast
	WindowEdgeWest
	WindowEdgeEast
	WindowEdgeSouthWest
	WindowEdgeSouth
	WindowEdgeSouthEast
)

func (e WindowEdge) FromString(value string) (enum interface{}, err error) {
	switch strings.ToLower(value) {
	case "north-west":
		return WindowEdgeNorthWest, nil
	case "north":
		return WindowEdgeNorth, nil
	case "north-east":
		return WindowEdgeNorthEast, nil
	case "west":
		return WindowEdgeWest, nil
	case "east":
		return WindowEdgeEast, nil
	case "south-west":
		return WindowEdgeSouthWest, nil
	case "south":
		return WindowEdgeSouth, nil
	case "south-east":
		return WindowEdgeSouthEast, nil
	}
	return nil, fmt.Errorf("unknown value for WindowEdge.FromString(%v)", value)
}

type Gravity uint64

const (
	// The reference point is at the top left corner.
	GravityNorthWest Gravity = iota
	// The reference point is in the middle of the top edge.
	GravityNorth
	// The reference point is at the top right corner.
	GravityNorthEast
	// The reference point is at the middle of the left edge.
	GravityWest
	// The reference point is at the center of the window.
	GravityCenter
	// The reference point is at the middle of the right edge.
	GravityEast
	// The reference point is at the lower left corner.
	GravitySouthWest
	// The reference point is at the middle of the lower edge.
	GravitySouth
	// The reference point is at the lower right corner.
	GravitySouthEast
	// The reference point is at the top left corner of the window itself, ignoring window manager decorations.
	GravityStatic
)

func (t Gravity) FromString(value string) (enum interface{}, err error) {
	switch strings.ToLower(value) {
	case "north-west", "northwest", "nw":
		return GravityNorthWest, nil
	case "north", "n":
		return GravityNorth, nil
	case "north-east", "northeast", "ne":
		return GravityNorthEast, nil
	case "west", "w":
		return GravityWest, nil
	case "center":
		return GravityCenter, nil
	case "east", "e":
		return GravityEast, nil
	case "south-west", "southwest", "sw":
		return GravitySouthWest, nil
	case "south", "s":
		return GravitySouth, nil
	case "south-east", "southeast", "se":
		return GravitySouthEast, nil
	case "static":
		return GravityStatic, nil
	}
	return nil, fmt.Errorf("unknown value for Gravity.FromString(%v)", value)
}

/* Window position */
type WindowPosition uint64

const (
	WinPosNone WindowPosition = iota
	WinPosCenter
	WinPosMouse
	WinPosCenterAlways
	WinPosCenterOnParent
)

/* Sort type */
type SortType uint64

const (
	SORT_ASCENDING SortType = iota
	SORT_DESCENDING
)

/* Preedit style */
type IMPreeditStyle uint64

const (
	IM_PREEDIT_NOTHING IMPreeditStyle = iota
	IM_PREEDIT_CALLBACK
	IM_PREEDIT_NONE
)

/* Status style */
type IMStatusStyle uint64

const (
	IM_STATUS_NOTHING IMStatusStyle = iota
	IM_STATUS_CALLBACK
	IM_STATUS_NONE
)

/* packing direction */
type PackDirection uint64

const (
	PACK_DIRECTION_LTR PackDirection = iota
	PACK_DIRECTION_RTL
	PACK_DIRECTION_TTB
	PACK_DIRECTION_BTT
)

/* Print pages */
type PrintPages uint64

const (
	PRINT_PAGES_ALL PrintPages = iota
	PRINT_PAGES_CURRENT
	PRINT_PAGES_RANGES
	PRINT_PAGES_SELECTION
)

/* Page set */
type PageSet uint64

const (
	PAGE_SET_ALL PageSet = iota
	PAGE_SET_EVEN
	PAGE_SET_ODD
)

/* Number up layout */
type NumberUpLayout uint64

const (
	NUMBER_UP_LAYOUT_LEFT_TO_RIGHT_TOP_TO_BOTTOM NumberUpLayout = iota
	NUMBER_UP_LAYOUT_LEFT_TO_RIGHT_BOTTOM_TO_TOP
	NUMBER_UP_LAYOUT_RIGHT_TO_LEFT_TOP_TO_BOTTOM
	NUMBER_UP_LAYOUT_RIGHT_TO_LEFT_BOTTOM_TO_TOP
	NUMBER_UP_LAYOUT_TOP_TO_BOTTOM_LEFT_TO_RIGHT
	NUMBER_UP_LAYOUT_TOP_TO_BOTTOM_RIGHT_TO_LEFT
	NUMBER_UP_LAYOUT_BOTTOM_TO_TOP_LEFT_TO_RIGHT
	NUMBER_UP_LAYOUT_BOTTOM_TO_TOP_RIGHT_TO_LEFT
)

/* Unit */
type Unit uint64

const (
	UNIT_PIXEL Unit = iota
	UNIT_POINTS
	UNIT_INCH
	UNIT_MM
)

/* Tree view grid lines */
type TreeViewGridLines uint64

const (
	TREE_VIEW_GRID_LINES_NONE TreeViewGridLines = iota
	TREE_VIEW_GRID_LINES_HORIZONTAL
	TREE_VIEW_GRID_LINES_VERTICAL
	TREE_VIEW_GRID_LINES_BOTH
)

/* File chooser action */
type FileChooserAction uint64

const (
	FILE_CHOOSER_ACTION_OPEN FileChooserAction = iota
	FILE_CHOOSER_ACTION_SAVE
	FILE_CHOOSER_ACTION_SELECT_FOLDER
	FILE_CHOOSER_ACTION_CREATE_FOLDER
)

/* File chooser confirmation */
type FileChooserConfirmation uint64

const (
	FILE_CHOOSER_CONFIRMATION_CONFIRM FileChooserConfirmation = iota
	FILE_CHOOSER_CONFIRMATION_ACCEPT_FILENAME
	FILE_CHOOSER_CONFIRMATION_SELECT_AGAIN
)

/* File chooser error */
type FileChooserError uint64

const (
	FILE_CHOOSER_ERROR_NONEXISTENT FileChooserError = iota
	FILE_CHOOSER_ERROR_BAD_FILENAME
	FILE_CHOOSER_ERROR_ALREADY_EXISTS
	FILE_CHOOSER_ERROR_INCOMPLETE_HOSTNAME
)

/* Load state */
type LoadState uint64

const (
	LOAD_EMPTY LoadState = iota
	LOAD_PRELOAD
	LOAD_LOADING
	LOAD_FINISHED
)

/* Reload state */
type ReloadState uint64

const (
	RELOAD_EMPTY ReloadState = iota
	RELOAD_HAS_FOLDER
)

/* Location mode */
type LocationMode uint64

const (
	LOCATION_MODE_PATH_BAR LocationMode = iota
	LOCATION_MODE_FILENAME_ENTRY
)

/* Operation mode */
type OperationMode uint64

const (
	OPERATION_MODE_BROWSE OperationMode = iota
	OPERATION_MODE_SEARCH
	OPERATION_MODE_RECENT
)

/* Startup mode */
type StartupMode uint64

const (
	STARTUP_MODE_RECENT StartupMode = iota
	STARTUP_MODE_CWD
)

/* File chooser prop */
type FileChooserProp uint64

const (
	FILE_CHOOSER_PROP_FIRST               FileChooserProp = 0
	FILE_CHOOSER_PROP_ACTION              FileChooserProp = FileChooserProp(FILE_CHOOSER_PROP_FIRST)
	FILE_CHOOSER_PROP_FILE_SYSTEM_BACKEND FileChooserProp = iota
	FILE_CHOOSER_PROP_FILTER
	FILE_CHOOSER_PROP_LOCAL_ONLY
	FILE_CHOOSER_PROP_PREVIEW_WIDGET
	FILE_CHOOSER_PROP_PREVIEW_WIDGET_ACTIVE
	FILE_CHOOSER_PROP_USE_PREVIEW_LABEL
	FILE_CHOOSER_PROP_EXTRA_WIDGET
	FILE_CHOOSER_PROP_SELECT_MULTIPLE
	FILE_CHOOSER_PROP_SHOW_HIDDEN
	FILE_CHOOSER_PROP_DO_OVERWRITE_CONFIRMATION
	FILE_CHOOSER_PROP_CREATE_FOLDERS
	FILE_CHOOSER_PROP_LAST FileChooserProp = FileChooserProp(FILE_CHOOSER_PROP_CREATE_FOLDERS)
)

/* File filter flags */
type FileFilterFlags uint64

const (
	FILE_FILTER_FILENAME FileFilterFlags = 1 << 0
	FILE_FILTER_URI      FileFilterFlags = 1 << iota
	FILE_FILTER_DISPLAY_NAME
	FILE_FILTER_MIME_TYPE
)

/* Icon lookup flags */
type IconLookupFlags uint64

const (
	ICON_LOOKUP_NO_SVG    IconLookupFlags = 1 << 0
	ICON_LOOKUP_FORCE_SVG IconLookupFlags = 1 << iota
	ICON_LOOKUP_USE_BUILTIN
	ICON_LOOKUP_GENERIC_FALLBACK
	ICON_LOOKUP_FORCE_SIZE
)

/* Icon style error */
type IconThemeError uint64

const (
	ICON_THEME_NOT_FOUND IconThemeError = iota
	ICON_THEME_FAILED
)

/* Icon view drop position */
type IconViewDropPosition uint64

const (
	ICON_VIEW_NO_DROP IconViewDropPosition = iota
	ICON_VIEW_DROP_INTO
	ICON_VIEW_DROP_LEFT
	ICON_VIEW_DROP_RIGHT
	ICON_VIEW_DROP_ABOVE
	ICON_VIEW_DROP_BELOW
)

/* Image type */
type ImageType uint64

const (
	IMAGE_EMPTY ImageType = iota
	IMAGE_PIXMAP
	IMAGE_IMAGE
	IMAGE_PIXBUF
	IMAGE_STOCK
	IMAGE_ICON_SET
	IMAGE_ANIMATION
	IMAGE_ICON_NAME
	IMAGE_GICON
)

/* Buttons type */
type ButtonsType uint64

const (
	BUTTONS_NONE ButtonsType = iota
	BUTTONS_OK
	BUTTONS_CLOSE
	BUTTONS_CANCEL
	BUTTONS_YES_NO
	BUTTONS_OK_CANCEL
)

/* Notebook tab */
type NotebookTab uint64

const (
	NOTEBOOK_TAB_FIRST NotebookTab = iota
	NOTEBOOK_TAB_LAST
)

/* Arg flags */
type ArgFlags uint64

const (
	ARC_READABLE       ArgFlags = ArgFlags(PARAM_READABLE)
	ARC_WRITABLE       ArgFlags = ArgFlags(PARAM_WRITABLE)
	ARC_CONSTRUCT      ArgFlags = ArgFlags(PARAM_CONSTRUCT)
	ARC_CONSTRUCT_ONLY ArgFlags = ArgFlags(PARAM_CONSTRUCT_ONLY)
	ARC_CHILD_ARG      ArgFlags = 1 << 4
)

/* Private flags */
type PrivateFlags uint64

const (
	PRIVATE_USER_STYLE      PrivateFlags = 1 << 0
	PRIVATE_RESIZE_PENDING  PrivateFlags = 1 << 2
	PRIVATE_HAS_POINTER     PrivateFlags = 1 << 3
	PRIVATE_SHADOWED        PrivateFlags = 1 << 4
	PRIVATE_HAS_SHAPE_MASK  PrivateFlags = 1 << 5
	PRIVATE_IN_REPARENT     PrivateFlags = 1 << 6
	PRIVATE_DIRECTION_SET   PrivateFlags = 1 << 7
	PRIVATE_DIRECTION_LTR   PrivateFlags = 1 << 8
	PRIVATE_ANCHORED        PrivateFlags = 1 << 9
	PRIVATE_CHILD_VISIBLE   PrivateFlags = 1 << 10
	PRIVATE_REDRAW_ON_ALLOC PrivateFlags = 1 << 11
	PRIVATE_ALLOC_NEEDED    PrivateFlags = 1 << 12
	PRIVATE_REQUEST_NEEDED  PrivateFlags = 1 << 13
)

/* Progress bar style */
type ProgressBarStyle uint64

const (
	PROGRESS_CONTINUOUS ProgressBarStyle = iota
	PROGRESS_DISCRETE
)

/* Progress bar orientation */
type ProgressBarOrientation uint64

const (
	PROGRESS_LEFT_TO_RIGHT ProgressBarOrientation = iota
	PROGRESS_RIGHT_TO_LEFT
	PROGRESS_BOTTOM_TO_TOP
	PROGRESS_TOP_TO_BOTTOM
)

/* Node color */
type RBNodeColor uint64

const (
	RBNODE_BLACK RBNodeColor = 1 << 0
	RBNODE_RED   RBNodeColor = 1 << iota
	RBNODE_IS_PARENT
	RBNODE_IS_SELECTED
	RBNODE_IS_PRELIT
	RBNODE_IS_SEMI_COLLAPSED
	RBNODE_IS_SEMI_EXPANDED
	RBNODE_INVALID
	RBNODE_COLUMN_INVALID
	RBNODE_DESCENDANTS_INVALID
	RBNODE_NON_COLORS RBNodeColor = RBNodeColor(RBNODE_IS_PARENT)
)

/* Rc flags */
type RcFlags uint64

const (
	RC_FG RcFlags = 1 << 0
	RC_BG RcFlags = 1 << iota
	RC_TEXT
	RC_BASE
)

/* Rc token type */
type RcTokenType uint64

const (
	RC_TOKEN_INVALID RcTokenType = RcTokenType(TOKEN_LAST)
	RC_TOKEN_INCLUDE RcTokenType = iota
	RC_TOKEN_NORMAL
	RC_TOKEN_ACTIVE
	RC_TOKEN_PRELIGHT
	RC_TOKEN_SELECTED
	RC_TOKEN_INSENSITIVE
	RC_TOKEN_FG
	RC_TOKEN_BG
	RC_TOKEN_TEXT
	RC_TOKEN_BASE
	RC_TOKEN_XTHICKNESS
	RC_TOKEN_YTHICKNESS
	RC_TOKEN_FONT
	RC_TOKEN_FONTSET
	RC_TOKEN_FONT_NAME
	RC_TOKEN_BC_PIXMAP
	RC_TOKEN_PIXMAP_PATH
	RC_TOKEN_STYLE
	RC_TOKEN_BINDING
	RC_TOKEN_BIND
	RC_TOKEN_WIDGET
	RC_TOKEN_WIDGET_CLASS
	RC_TOKEN_CLASS
	RC_TOKEN_LOWEST
	RC_TOKEN_CTK
	RC_TOKEN_APPLICATION
	RC_TOKEN_THEME
	RC_TOKEN_RC
	RC_TOKEN_HIGHEST
	RC_TOKEN_ENGINE
	RC_TOKEN_MODULE_PATH
	RC_TOKEN_IM_MODULE_PATH
	RC_TOKEN_IM_MODULE_FILE
	RC_TOKEN_STOCK
	RC_TOKEN_LTR
	RC_TOKEN_RTL
	RC_TOKEN_COLOR
	RC_TOKEN_UNBIND
	RC_TOKEN_LAST
)

/* Recent sort type */
type RecentSortType uint64

const (
	RECENT_SORT_NONE RecentSortType = 0
	RECENT_SORT_MRU  RecentSortType = iota
	RECENT_SORT_LRU
	RECENT_SORT_CUSTOM
)

/* Recent chooser error */
type RecentChooserError uint64

const (
	RECENT_CHOOSER_ERROR_NOT_FOUND RecentChooserError = iota
	RECENT_CHOOSER_ERROR_INVALID_URI
)

/* Recent chooser prop */
type RecentChooserProp uint64

const (
	RECENT_CHOOSER_PROP_FIRST          RecentChooserProp = 0
	RECENT_CHOOSER_PROP_RECENT_MANAGER RecentChooserProp = iota
	RECENT_CHOOSER_PROP_SHOW_PRIVATE
	RECENT_CHOOSER_PROP_SHOW_NOT_FOUND
	RECENT_CHOOSER_PROP_SHOW_TIPS
	RECENT_CHOOSER_PROP_SHOW_ICONS
	RECENT_CHOOSER_PROP_SELECT_MULTIPLE
	RECENT_CHOOSER_PROP_LIMIT
	RECENT_CHOOSER_PROP_LOCAL_ONLY
	RECENT_CHOOSER_PROP_SORT_TYPE
	RECENT_CHOOSER_PROP_FILTER
	RECENT_CHOOSER_PROP_LAST
)

/* Recent filter flags */
type RecentFilterFlags uint64

const (
	RECENT_FILTER_URI          RecentFilterFlags = 1 << 0
	RECENT_FILTER_DISPLAY_NAME RecentFilterFlags = 1 << iota
	RECENT_FILTER_MIME_TYPE
	RECENT_FILTER_APPLICATION
	RECENT_FILTER_GROUP
	RECENT_FILTER_AGE
)

/* Recent manager error */
type RecentManagerError uint64

const (
	RECENT_MANAGER_ERROR_NOT_FOUND RecentManagerError = iota
	RECENT_MANAGER_ERROR_INVALID_URI
	RECENT_MANAGER_ERROR_INVALID_ENCODING
	RECENT_MANAGER_ERROR_NOT_REGISTERED
	RECENT_MANAGER_ERROR_READ
	RECENT_MANAGER_ERROR_WRITE
	RECENT_MANAGER_ERROR_UNKNOWN
)

/* Size group mode */
type SizeGroupMode uint64

const (
	SIZE_GROUP_NONE SizeGroupMode = iota
	SIZE_GROUP_HORIZONTAL
	SIZE_GROUP_VERTICAL
	SIZE_GROUP_BOTH
)

/* Spin button update policy */
type SpinButtonUpdatePolicy uint64

const (
	UPDATE_ALWAYS SpinButtonUpdatePolicy = iota
	UPDATE_IF_VALID
)

/* Spin type */
type SpinType uint64

const (
	SPIN_STEP_FORWARD SpinType = iota
	SPIN_STEP_BACKWARD
	SPIN_PAGE_FORWARD
	SPIN_PAGE_BACKWARD
	SPIN_HOME
	SPIN_END
	SPIN_USER_DEFINED
)

/* Text buffer target info */
type TextBufferTargetInfo int

const (
	TEXT_BUFFER_TARGET_INFO_BUFFER_CONTENTS TextBufferTargetInfo = -1
	TEXT_BUFFER_TARGET_INFO_RICH_TEXT       TextBufferTargetInfo = -2
	TEXT_BUFFER_TARGET_INFO_TEXT            TextBufferTargetInfo = -3
)

/* Text search flags */
type TextSearchFlags uint64

const (
	TEXT_SEARCH_VISIBLE_ONLY TextSearchFlags = 1 << 0
	TEXT_SEARCH_TEXT_ONLY    TextSearchFlags = 1 << iota
)

/* Text window type */
type TextWindowType uint64

const (
	TEXT_WINDOW_PRIVATE TextWindowType = iota
	TEXT_WINDOW_WIDGET
	TEXT_WINDOW_TEXT
	TEXT_WINDOW_LEFT
	TEXT_WINDOW_RIGHT
	TEXT_WINDOW_TOP
	TEXT_WINDOW_BOTTOM
)

/* Toolbar child type */
type ToolbarChildType uint64

const (
	TOOLBAR_CHILD_SPACE ToolbarChildType = iota
	TOOLBAR_CHILD_BUTTON
	TOOLBAR_CHILD_TOGGLEBUTTON
	TOOLBAR_CHILD_RADIOBUTTON
	TOOLBAR_CHILD_WIDGET
)

/* Toolbar space style */
type ToolbarSpaceStyle uint64

const (
	TOOLBAR_SPACE_EMPTY ToolbarSpaceStyle = iota
	TOOLBAR_SPACE_LINE
)

/* Tree view mode */
type TreeViewMode uint64

const (
	TREE_VIEW_LINE TreeViewMode = iota
	TREE_VIEW_ITEM
)

/* Tree model flags */
type TreeModelFlags uint64

const (
	TREE_MODEL_ITERS_PERSIST TreeModelFlags = 1 << 0
	TREE_MODEL_LIST_ONLY     TreeModelFlags = 1 << iota
)

/* Tree view flags */
type TreeViewFlags uint64

const (
	TREE_VIEW_IS_LIST        TreeViewFlags = 1 << 0
	TREE_VIEW_SHOW_EXPANDERS TreeViewFlags = 1 << iota
	TREE_VIEW_IN_COLUMN_RESIZE
	TREE_VIEW_ARROW_PRELIT
	TREE_VIEW_HEADERS_VISIBLE
	TREE_VIEW_DRAW_KEYFOCUS
	TREE_VIEW_MODEL_SETUP
	TREE_VIEW_IN_COLUMN_DRAG
)

/* Tree view drop position */
type TreeViewDropPosition uint64

const (
	TREE_VIEW_DROP_BEFORE TreeViewDropPosition = iota
	TREE_VIEW_DROP_AFTER
	TREE_VIEW_DROP_INTO_OR_BEFORE
	TREE_VIEW_DROP_INTO_OR_AFTER
)

/* Tree view column sizing */
type TreeViewColumnSizing uint64

const (
	TREE_VIEW_COLUMN_GROW_ONLY TreeViewColumnSizing = iota
	TREE_VIEW_COLUMN_AUTOSIZE
	TREE_VIEW_COLUMN_FIXED
)

/* Manager item type */
type UIManagerItemType uint64

const (
	UI_MANAGER_AUTO    UIManagerItemType = 0
	UI_MANAGER_MENUBAR UIManagerItemType = 1 << 0
	UI_MANAGER_MENU    UIManagerItemType = 1 << iota
	UI_MANAGER_TOOLBAR
	UI_MANAGER_PLACEHOLDER
	UI_MANAGER_POPUP
	UI_MANAGER_MENUITEM
	UI_MANAGER_TOOLITEM
	UI_MANAGER_SEPARATOR
	UI_MANAGER_ACCELERATOR
	UI_MANAGER_POPUP_WITH_ACCELS
)

/* Widget flags */
type WidgetFlags uint64

func (f WidgetFlags) HasBit(flag WidgetFlags) bool {
	return cbits.Has(uint64(f), uint64(flag))
}

const (
	NULL_WIDGET_FLAG WidgetFlags = 0
	TOPLEVEL         WidgetFlags = 1 << (iota + 2)
	NO_WINDOW
	REALIZED
	MAPPED
	VISIBLE
	SENSITIVE
	PARENT_SENSITIVE
	CAN_FOCUS
	HAS_FOCUS
	CAN_DEFAULT
	HAS_DEFAULT
	HAS_GRAB
	RC_STYLE
	COMPOSITE_CHILD
	NO_REPARENT
	APP_PAINTABLE
	RECEIVES_DEFAULT
	DOUBLE_BUFFERED
	NO_SHOW_ALL
	INVALID_WIDGET_FLAG
)

/* Widget help type */
type WidgetHelpType uint64

const (
	WIDGET_HELP_TOOLTIP WidgetHelpType = iota
	WIDGET_HELP_WHATS_THIS
)

/* Param flags */
type ParamFlags uint64

const (
	PARAM_READABLE  ParamFlags = 1 << 0
	PARAM_WRITABLE  ParamFlags = 1 << iota
	PARAM_READWRITE ParamFlags = ParamFlags(PARAM_READABLE | PARAM_WRITABLE)
	PARAM_CONSTRUCT
	PARAM_CONSTRUCT_ONLY
	PARAM_LAX_VALIDATION
	PARAM_STATIC_NAME
	PARAM_PRIVATE ParamFlags = ParamFlags(PARAM_STATIC_NAME)
	PARAM_STATIC_NICK
	PARAM_STATIC_BLURB
	PARAM_EXPLICIT_NOTIFY
)

/* Error type */
type ErrorType uint64

const (
	ERR_UNKNOWN ErrorType = iota
	ERR_UNEXP_EOF
	ERR_UNEXP_EOF_IN_STRING
	ERR_UNEXP_EOF_IN_COMMENT
	ERR_NON_DIGIT_IN_CONST
	ERR_DIGIT_RADIX
	ERR_FLOAT_RADIX
	ERR_FLOAT_MALFORMED
)

/* Token type */
type TokenType uint64

const (
	TOKEN_EOF   TokenType = 0
	TOKEN_NONE  TokenType = 256
	TOKEN_ERROR TokenType = iota
	TOKEN_CHAR
	TOKEN_BINARY
	TOKEN_OCTAL
	TOKEN_INT
	TOKEN_HEX
	TOKEN_FLOAT
	TOKEN_STRING
	TOKEN_SYMBOL
	TOKEN_IDENTIFIER
	TOKEN_IDENTIFIER_NULL
	TOKEN_COMMENT_SINGLE
	TOKEN_COMMENT_MULTI
	TOKEN_LAST
)

type ExtensionMode uint64

const (
	EXTENSION_EVENTS_NONE ExtensionMode = iota
	EXTENSION_EVENTS_ALL
	EXTENSION_EVENTS_CURSOR
)
