// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package enums

import (
	"fmt"
	"strings"
)

type AccelFlags uint64

const (
	ACCEL_VISIBLE AccelFlags = 1 << 0
	ACCEL_LOCKED  AccelFlags = 1 << iota
	ACCEL_MASK    AccelFlags = 0
)

type AssistantPageType uint64

const (
	ASSISTANT_PAGE_CONTENT AssistantPageType = iota
	ASSISTANT_PAGE_INTRO
	ASSISTANT_PAGE_CONFIRM
	ASSISTANT_PAGE_SUMMARY
	ASSISTANT_PAGE_PROGRESS
)

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

type CalendarDisplayOptions uint64

const (
	CALENDAR_SHOW_HEADING   CalendarDisplayOptions = 1 << 0
	CALENDAR_SHOW_DAY_NAMES CalendarDisplayOptions = 1 << iota
	CALENDAR_NO_MONTH_CHANGE
	CALENDAR_SHOW_WEEK_NUMBERS
	CALENDAR_WEEK_START_MONDAY
	CALENDAR_SHOW_DETAILS
)

type CellRendererState uint64

const (
	CELL_RENDERER_SELECTED CellRendererState = 1 << 0
	CELL_RENDERER_PRELIT   CellRendererState = 1 << iota
	CELL_RENDERER_INSENSITIVE
	CELL_RENDERER_SORTED
	CELL_RENDERER_FOCUSED
)

type CellRendererMode uint64

const (
	CELL_RENDERER_MODE_INERT CellRendererMode = iota
	CELL_RENDERER_MODE_ACTIVATABLE
	CELL_RENDERER_MODE_EDITABLE
)

type CellRendererAccelMode uint64

const (
	CELL_RENDERER_ACCEL_MODE_CTK CellRendererAccelMode = iota
	CELL_RENDERER_ACCEL_MODE_OTHER
)

type CellType uint64

const (
	CELL_EMPTY CellType = iota
	CELL_TEXT
	CELL_PIXMAP
	CELL_PIXTEXT
	CELL_WIDGET
)

type CListDragPos uint64

const (
	CLIST_DRAC_NONE CListDragPos = iota
	CLIST_DRAC_BEFORE
	CLIST_DRAC_INTO
	CLIST_DRAC_AFTER
)

type ButtonAction uint64

const (
	BUTTON_IGNORED ButtonAction = 1 << iota
	BUTTON_SELECTS
	BUTTON_DRAGS
	BUTTON_EXPANDS
)

type CTreePos uint64

const (
	CTREE_POS_BEFORE CTreePos = iota
	CTREE_POS_AS_CHILD
	CTREE_POS_AFTER
)

type CTreeLineStyle uint64

const (
	CTREE_LINES_NONE CTreeLineStyle = iota
	CTREE_LINES_SOLID
	CTREE_LINES_DOTTED
	CTREE_LINES_TABBED
)

type CTreeExpanderStyle uint64

const (
	CTREE_EXPANDER_NONE CTreeExpanderStyle = iota
	CTREE_EXPANDER_SQUARE
	CTREE_EXPANDER_TRIANGLE
	CTREE_EXPANDER_CIRCULAR
)

type CTreeExpansionType uint64

const (
	CTREE_EXPANSION_EXPAND CTreeExpansionType = iota
	CTREE_EXPANSION_EXPAND_RECURSIVE
	CTREE_EXPANSION_COLLAPSE
	CTREE_EXPANSION_COLLAPSE_RECURSIVE
	CTREE_EXPANSION_TOGGLE
	CTREE_EXPANSION_TOGGLE_RECURSIVE
)

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

type DialogFlags uint64

const (
	DialogModal             DialogFlags = 1 << 0
	DialogDestroyWithParent DialogFlags = 1 << iota
	DialogNoSeparator
)

type EntryIconPosition uint64

const (
	ENTRY_ICON_PRIMARY EntryIconPosition = iota
	ENTRY_ICON_SECONDARY
)

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

type ArrowPlacement uint64

const (
	ARROWS_BOTH ArrowPlacement = iota
	ARROWS_START
	ARROWS_END
)

type ArrowType uint64

const (
	ArrowUp ArrowType = iota
	ArrowDown
	ArrowLeft
	ArrowRight
	ArrowNone
)

type AttachOptions uint64

const (
	EXPAND AttachOptions = 1 << 0
	SHRINK AttachOptions = 1 << iota
	FILL
)

type ButtonBoxStyle uint64

const (
	BUTTONBOX_DEFAULT_STYLE ButtonBoxStyle = iota
	BUTTONBOX_SPREAD                       // Spread the buttons evenly and centered away from the edges
	BUTTONBOX_EDGE                         // Spread the buttons evenly and centered from edge to edge
	BUTTONBOX_START                        // Group buttons at the start
	BUTTONBOX_END                          // Group buttons at the end
	BUTTONBOX_CENTER                       // Group buttons at the center
	BUTTONBOX_EXPAND                       // Buttons are expanded to evenly consume all space available
)

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

type DirectionType uint64

const (
	DIR_TAB_FORWARD DirectionType = iota
	DIR_TAB_BACKWARD
	DIR_UP
	DIR_DOWN
	DIR_LEFT
	DIR_RIGHT
)

type ExpanderStyle uint64

const (
	EXPANDER_COLLAPSED ExpanderStyle = iota
	EXPANDER_SEMI_COLLAPSED
	EXPANDER_SEMI_EXPANDED
	EXPANDER_EXPANDED
)

type SensitivityType uint64

const (
	SensitivityAuto SensitivityType = iota
	SensitivityOn
	SensitivityOff
)

type SideType uint64

const (
	SIDE_TOP SideType = iota
	SIDE_BOTTOM
	SIDE_LEFT
	SIDE_RIGHT
)

type TextDirection uint64

const (
	TextDirNone TextDirection = iota
	TextDirLtr
	TextDirRtl
)

type MatchType uint64

const (
	MATCH_ALL MatchType = iota
	MATCH_ALL_TAIL
	MATCH_HEAD
	MATCH_TAIL
	MATCH_EXACT
	MATCH_LAST
)

type MenuDirectionType uint64

const (
	MENU_DIR_PARENT MenuDirectionType = iota
	MENU_DIR_CHILD
	MENU_DIR_NEXT
	MENU_DIR_PREV
)

type MessageType uint64

const (
	MESSAGE_INFO MessageType = iota
	MESSAGE_WARNING
	MESSAGE_QUESTION
	MESSAGE_ERROR
	MESSAGE_OTHER
)

type MetricType uint64

const (
	PIXELS MetricType = iota
	INCHES
	CENTIMETERS
)

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

type ScrollStep uint64

const (
	SCROLL_STEPS ScrollStep = iota
	SCROLL_PAGES
	SCROLL_ENDS
	SCROLL_HORIZONTAL_STEPS
	SCROLL_HORIZONTAL_PAGES
	SCROLL_HORIZONTAL_ENDS
)

type CornerType uint64

const (
	CornerTopLeft CornerType = iota
	CornerBottomLeft
	CornerTopRight
	CornerBottomRight
)

type PackType uint64

const (
	PackStart PackType = iota
	PackEnd
)

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

type PathPriorityType uint64

const (
	PATH_PRIO_LOWEST      PathPriorityType = 0
	PATH_PRIO_CTK         PathPriorityType = 4
	PATH_PRIO_APPLICATION PathPriorityType = 8
	PATH_PRIO_THEME       PathPriorityType = 10
	PATH_PRIO_RC          PathPriorityType = 12
	PATH_PRIO_HIGHEST     PathPriorityType = 15
)

type PathType uint64

const (
	PATH_WIDGET PathType = iota
	PATH_WIDGET_CLASS
	PATH_CLASS
)

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
	return nil, fmt.Errorf("unknown value for PolicyType.FromString(%v)", value)
}

type PositionType uint64

const (
	POS_LEFT PositionType = iota
	POS_RIGHT
	POS_TOP
	POS_BOTTOM
)

type ReliefStyle uint64

const (
	RELIEF_NORMAL ReliefStyle = iota
	RELIEF_HALF
	RELIEF_NONE
)

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

type SelectionMode uint64

const (
	SELECTION_NONE SelectionMode = iota
	SELECTION_SINGLE
	SELECTION_BROWSE
	SELECTION_MULTIPLE
	SELECTION_EXTENDED SelectionMode = SelectionMode(SELECTION_MULTIPLE)
)

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
	StateNone   StateType = 0
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

// func (s StateType) HasBit(state StateType) bool {
// 	return cbits.Has(uint64(s), uint64(state))
// }
//
// func (s StateType) String() (label string) {
// 	label = ""
// 	update := func(state StateType, name string) {
// 		if s.HasBit(state) {
// 			if len(label) > 0 {
// 				label += " | "
// 			}
// 			label += name
// 		}
// 	}
// 	update(StateNormal, "normal")
// 	update(StateActive, "active")
// 	update(StatePrelight, "prelight")
// 	update(StateInsensitive, "insensitive")
// 	update(StateSelected, "selected")
// 	if label == "" {
// 		label = "unknown"
// 	}
// 	return
// }

type SubmenuDirection uint64

const (
	DIRECTION_LEFT SubmenuDirection = iota
	DIRECTION_RIGHT
)

type SubmenuPlacement uint64

const (
	TOP_BOTTOM SubmenuPlacement = iota
	LEFT_RIGHT
)

type ToolbarStyle uint64

const (
	TOOLBAR_ICONS ToolbarStyle = iota
	TOOLBAR_TEXT
	TOOLBAR_BOTH
	TOOLBAR_BOTH_HORIZ
)

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
	return nil, fmt.Errorf("unknown value for UpdateType.FromString(%v)", value)
}

type Visibility uint64

const (
	VISIBILITY_NONE Visibility = iota
	VISIBILITY_PARTIAL
	VISIBILITY_FULL
)

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
	// A utility window
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

type WindowPosition uint64

const (
	WinPosNone WindowPosition = iota
	WinPosCenter
	WinPosMouse
	WinPosCenterAlways
	WinPosCenterOnParent
)

type SortType uint64

const (
	SORT_ASCENDING SortType = iota
	SORT_DESCENDING
)

type IMPreeditStyle uint64

const (
	IM_PREEDIT_NOTHING IMPreeditStyle = iota
	IM_PREEDIT_CALLBACK
	IM_PREEDIT_NONE
)

type IMStatusStyle uint64

const (
	IM_STATUS_NOTHING IMStatusStyle = iota
	IM_STATUS_CALLBACK
	IM_STATUS_NONE
)

type PackDirection uint64

const (
	PACK_DIRECTION_LTR PackDirection = iota
	PACK_DIRECTION_RTL
	PACK_DIRECTION_TTB
	PACK_DIRECTION_BTT
)

type PrintPages uint64

const (
	PRINT_PAGES_ALL PrintPages = iota
	PRINT_PAGES_CURRENT
	PRINT_PAGES_RANGES
	PRINT_PAGES_SELECTION
)

type PageSet uint64

const (
	PAGE_SET_ALL PageSet = iota
	PAGE_SET_EVEN
	PAGE_SET_ODD
)

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

type Unit uint64

const (
	UNIT_PIXEL Unit = iota
	UNIT_POINTS
	UNIT_INCH
	UNIT_MM
)

type TreeViewGridLines uint64

const (
	TREE_VIEW_GRID_LINES_NONE TreeViewGridLines = iota
	TREE_VIEW_GRID_LINES_HORIZONTAL
	TREE_VIEW_GRID_LINES_VERTICAL
	TREE_VIEW_GRID_LINES_BOTH
)

type FileChooserAction uint64

const (
	FILE_CHOOSER_ACTION_OPEN FileChooserAction = iota
	FILE_CHOOSER_ACTION_SAVE
	FILE_CHOOSER_ACTION_SELECT_FOLDER
	FILE_CHOOSER_ACTION_CREATE_FOLDER
)

type FileChooserConfirmation uint64

const (
	FILE_CHOOSER_CONFIRMATION_CONFIRM FileChooserConfirmation = iota
	FILE_CHOOSER_CONFIRMATION_ACCEPT_FILENAME
	FILE_CHOOSER_CONFIRMATION_SELECT_AGAIN
)

type FileChooserError uint64

const (
	FILE_CHOOSER_ERROR_NONEXISTENT FileChooserError = iota
	FILE_CHOOSER_ERROR_BAD_FILENAME
	FILE_CHOOSER_ERROR_ALREADY_EXISTS
	FILE_CHOOSER_ERROR_INCOMPLETE_HOSTNAME
)

type LoadState uint64

const (
	LOAD_EMPTY LoadState = iota
	LOAD_PRELOAD
	LOAD_LOADING
	LOAD_FINISHED
)

type ReloadState uint64

const (
	RELOAD_EMPTY ReloadState = iota
	RELOAD_HAS_FOLDER
)

type LocationMode uint64

const (
	LOCATION_MODE_PATH_BAR LocationMode = iota
	LOCATION_MODE_FILENAME_ENTRY
)

type OperationMode uint64

const (
	OPERATION_MODE_BROWSE OperationMode = iota
	OPERATION_MODE_SEARCH
	OPERATION_MODE_RECENT
)

type StartupMode uint64

const (
	STARTUP_MODE_RECENT StartupMode = iota
	STARTUP_MODE_CWD
)

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

type FileFilterFlags uint64

const (
	FILE_FILTER_FILENAME FileFilterFlags = 1 << 0
	FILE_FILTER_URI      FileFilterFlags = 1 << iota
	FILE_FILTER_DISPLAY_NAME
	FILE_FILTER_MIME_TYPE
)

type IconThemeError uint64

const (
	ICON_THEME_NOT_FOUND IconThemeError = iota
	ICON_THEME_FAILED
)

type ButtonsType uint64

const (
	BUTTONS_NONE ButtonsType = iota
	BUTTONS_OK
	BUTTONS_CLOSE
	BUTTONS_CANCEL
	BUTTONS_YES_NO
	BUTTONS_OK_CANCEL
)

type NotebookTab uint64

const (
	NOTEBOOK_TAB_FIRST NotebookTab = iota
	NOTEBOOK_TAB_LAST
)

type ArgFlags uint64

const (
	ARC_READABLE       ArgFlags = ArgFlags(PARAM_READABLE)
	ARC_WRITABLE       ArgFlags = ArgFlags(PARAM_WRITABLE)
	ARC_CONSTRUCT      ArgFlags = ArgFlags(PARAM_CONSTRUCT)
	ARC_CONSTRUCT_ONLY ArgFlags = ArgFlags(PARAM_CONSTRUCT_ONLY)
	ARC_CHILD_ARG      ArgFlags = 1 << 4
)

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

type ProgressBarStyle uint64

const (
	PROGRESS_CONTINUOUS ProgressBarStyle = iota
	PROGRESS_DISCRETE
)

type ProgressBarOrientation uint64

const (
	PROGRESS_LEFT_TO_RIGHT ProgressBarOrientation = iota
	PROGRESS_RIGHT_TO_LEFT
	PROGRESS_BOTTOM_TO_TOP
	PROGRESS_TOP_TO_BOTTOM
)

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

type RcFlags uint64

const (
	RC_FG RcFlags = 1 << 0
	RC_BG RcFlags = 1 << iota
	RC_TEXT
	RC_BASE
)

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

type RecentSortType uint64

const (
	RECENT_SORT_NONE RecentSortType = 0
	RECENT_SORT_MRU  RecentSortType = iota
	RECENT_SORT_LRU
	RECENT_SORT_CUSTOM
)

type RecentChooserError uint64

const (
	RECENT_CHOOSER_ERROR_NOT_FOUND RecentChooserError = iota
	RECENT_CHOOSER_ERROR_INVALID_URI
)

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

type RecentFilterFlags uint64

const (
	RECENT_FILTER_URI          RecentFilterFlags = 1 << 0
	RECENT_FILTER_DISPLAY_NAME RecentFilterFlags = 1 << iota
	RECENT_FILTER_MIME_TYPE
	RECENT_FILTER_APPLICATION
	RECENT_FILTER_GROUP
	RECENT_FILTER_AGE
)

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

type SizeGroupMode uint64

const (
	SIZE_GROUP_NONE SizeGroupMode = iota
	SIZE_GROUP_HORIZONTAL
	SIZE_GROUP_VERTICAL
	SIZE_GROUP_BOTH
)

type SpinButtonUpdatePolicy uint64

const (
	UPDATE_ALWAYS SpinButtonUpdatePolicy = iota
	UPDATE_IF_VALID
)

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

type TextBufferTargetInfo int

const (
	TEXT_BUFFER_TARGET_INFO_BUFFER_CONTENTS TextBufferTargetInfo = -1
	TEXT_BUFFER_TARGET_INFO_RICH_TEXT       TextBufferTargetInfo = -2
	TEXT_BUFFER_TARGET_INFO_TEXT            TextBufferTargetInfo = -3
)

type TextSearchFlags uint64

const (
	TEXT_SEARCH_VISIBLE_ONLY TextSearchFlags = 1 << 0
	TEXT_SEARCH_TEXT_ONLY    TextSearchFlags = 1 << iota
)

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

type ToolbarChildType uint64

const (
	TOOLBAR_CHILD_SPACE ToolbarChildType = iota
	TOOLBAR_CHILD_BUTTON
	TOOLBAR_CHILD_TOGGLEBUTTON
	TOOLBAR_CHILD_RADIOBUTTON
	TOOLBAR_CHILD_WIDGET
)

type ToolbarSpaceStyle uint64

const (
	TOOLBAR_SPACE_EMPTY ToolbarSpaceStyle = iota
	TOOLBAR_SPACE_LINE
)

type TreeViewMode uint64

const (
	TREE_VIEW_LINE TreeViewMode = iota
	TREE_VIEW_ITEM
)

type TreeModelFlags uint64

const (
	TREE_MODEL_ITERS_PERSIST TreeModelFlags = 1 << 0
	TREE_MODEL_LIST_ONLY     TreeModelFlags = 1 << iota
)

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

type TreeViewDropPosition uint64

const (
	TREE_VIEW_DROP_BEFORE TreeViewDropPosition = iota
	TREE_VIEW_DROP_AFTER
	TREE_VIEW_DROP_INTO_OR_BEFORE
	TREE_VIEW_DROP_INTO_OR_AFTER
)

type TreeViewColumnSizing uint64

const (
	TREE_VIEW_COLUMN_GROW_ONLY TreeViewColumnSizing = iota
	TREE_VIEW_COLUMN_AUTOSIZE
	TREE_VIEW_COLUMN_FIXED
)

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

// TODO: enforce WidgetFlags properly

type WidgetFlags uint64

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
	COMPOSITE_PARENT
	INVALID_WIDGET_FLAG
)

type WidgetHelpType uint64

const (
	WIDGET_HELP_TOOLTIP WidgetHelpType = iota
	WIDGET_HELP_WHATS_THIS
)

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

type GCallback = func()

type GClosure = func(argv ...interface{}) (handled bool)

//go:generate stringer -output enums_string.go -type AssistantPageType,BuilderError,CellRendererMode,CellRendererAccelMode,CellType,CListDragPos,CTreePos,CTreeLineStyle,CTreeExpanderStyle,CTreeExpansionType,EntryIconPosition,AnchorType,ArrowPlacement,ArrowType,ButtonBoxStyle,DeleteType,DirectionType,ExpanderStyle,SensitivityType,SideType,TextDirection,MatchType,MenuDirectionType,MessageType,MetricType,MovementStep,ScrollStep,CornerType,PackType,LayoutStyle,PathPriorityType,PathType,PolicyType,PositionType,ReliefStyle,ScrollType,SelectionMode,ShadowType,SubmenuDirection,SubmenuPlacement,ToolbarStyle,UpdateType,Visibility,WindowTypeHint,WindowEdge,Gravity,WindowPosition,SortType,IMPreeditStyle,IMStatusStyle,PackDirection,PrintPages,PageSet,NumberUpLayout,Unit,TreeViewGridLines,FileChooserAction,FileChooserConfirmation,FileChooserError,LoadState,ReloadState,LocationMode,OperationMode,StartupMode,FileChooserProp,IconThemeError,ButtonsType,NotebookTab,ArgFlags,ProgressBarStyle,ProgressBarOrientation,RcTokenType,RecentSortType,RecentChooserError,RecentChooserProp,RecentManagerError,SizeGroupMode,SpinButtonUpdatePolicy,SpinType,TextBufferTargetInfo,TextWindowType,ToolbarChildType,ToolbarSpaceStyle,TreeViewMode,TreeViewDropPosition,TreeViewColumnSizing,WidgetHelpType,ErrorType,TokenType,ExtensionMode
//go:generate bitmasker -output enums_bitmask.go -kebab -type AccelFlags,CalendarDisplayOptions,CellRendererState,ButtonAction,DebugFlag,DialogFlags,AttachOptions,StateType,FileFilterFlags,PrivateFlags,RBNodeColor,RcFlags,RecentFilterFlags,TextSearchFlags,TreeModelFlags,TreeViewFlags,UIManagerItemType,WidgetFlags,ParamFlags