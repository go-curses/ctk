package ctk

import (
	"github.com/go-curses/cdk"
)

// Editable Hierarchy:
//	CInterface
//	  +- Editable
type Editable interface {
	/* Base Interface */

	// Selects a region of text. The characters that are selected are those
	// characters at positions from start_pos up to, but not including end_pos .
	// If end_pos is negative, then the the characters selected are those
	// characters from start_pos to the end of the text. Note that positions are
	// specified in characters, not bytes.
	// Parameters:
	// 	startPos	start of region
	// 	endPos	end of region
	SelectRegion(startPos int, endPos int)

	// Retrieves the selection bound of the editable. start_pos will be filled
	// with the start of the selection and end_pos with end. If no text was
	// selected both will be identical and FALSE will be returned. Note that
	// positions are specified in characters, not bytes.
	// Parameters:
	// 	startPos	location to store the starting position, or NULL.
	// 	endPos	location to store the end position, or NULL.
	// Returns:
	// 	TRUE if an area is selected, FALSE otherwise
	GetSelectionBounds() (startPos, endPos int, value bool)

	// Inserts new_text_length bytes of new_text into the contents of the widget,
	// at position position . Note that the position is in characters, not in
	// bytes. The function updates position to point after the newly inserted
	// text.
	// Parameters:
	// 	newText	the text to append
	// 	newTextLength	the length of the text in bytes, or -1
	// 	position	location of the position text will be inserted at.
	InsertText(newText string, position int)

	// Deletes a sequence of characters. The characters that are deleted are
	// those characters at positions from start_pos up to, but not including
	// end_pos . If end_pos is negative, then the the characters deleted are
	// those from start_pos to the end of the text. Note that the positions are
	// specified in characters, not bytes.
	// Parameters:
	// 	startPos	start position
	// 	endPos	end position
	DeleteText(startPos int, endPos int)

	// Retrieves a sequence of characters. The characters that are retrieved are
	// those characters at positions from start_pos up to, but not including
	// end_pos . If end_pos is negative, then the the characters retrieved are
	// those characters from start_pos to the end of the text. Note that
	// positions are specified in characters, not bytes.
	// Parameters:
	// 	startPos	start of text
	// 	endPos	end of text
	// Returns:
	// 	a pointer to the contents of the widget as a string. This
	// 	string is allocated by the Editable implementation and
	// 	should be freed by the caller.
	GetChars(startPos int, endPos int) (value string)

	// Removes the contents of the currently selected content in the editable and
	// puts it on the clipboard.
	CutClipboard()

	// Copies the contents of the currently selected content in the editable and
	// puts it on the clipboard.
	CopyClipboard()

	// Pastes the content of the clipboard to the current position of the cursor
	// in the editable.
	PasteClipboard()

	// Deletes the currently selected text of the editable. This call doesn't do
	// anything if there is no selected text.
	DeleteSelection()

	// Sets the cursor position in the editable to the given value. The cursor is
	// displayed before the character with the given (base 0) index in the
	// contents of the editable. The value must be less than or equal to the
	// number of characters in the editable. A value of -1 indicates that the
	// position should be set after the last character of the editable. Note that
	// position is in characters, not in bytes.
	// Parameters:
	// 	position	the position of the cursor
	SetPosition(position int)

	// Retrieves the current position of the cursor relative to the start of the
	// content of the editable. Note that this position is in characters, not in
	// bytes.
	// Returns:
	// 	the cursor position
	GetPosition() (value int)

	// Determines if the user can edit the text in the editable widget or not.
	// Parameters:
	// 	isEditable	TRUE if the user is allowed to edit the text
	// in the widget
	SetEditable(isEditable bool)

	// Retrieves whether editable is editable. See SetEditable.
	// Returns:
	// 	TRUE if editable is editable.
	GetEditable() (value bool)
}

// The ::changed signal is emitted at the end of a single user-visible
// operation on the contents of the Editable. E.g., a paste operation that
// replaces the contents of the selection will cause only one signal emission
// (even though it is implemented by first deleting the selection, then
// inserting the new content, and may cause multiple ::notify::text signals
// to be emitted).
const SignalChangedText cdk.Signal = "changed"

// This signal is emitted when text is deleted from the widget by the user.
// The default handler for this signal will normally be responsible for
// deleting the text, so by connecting to this signal and then stopping the
// signal with g_signal_stop_emission, it is possible to modify the range
// of deleted text, or prevent it from being deleted entirely. The start_pos
// and end_pos parameters are interpreted as for DeleteText.
// Listener function arguments:
// 	startPos int	the starting position
// 	endPos int	the end position
const SignalDeleteText cdk.Signal = "delete-text"

// This signal is emitted when text is inserted into the widget by the user.
// The default handler for this signal will normally be responsible for
// inserting the text, so by connecting to this signal and then stopping the
// signal with g_signal_stop_emission, it is possible to modify the
// inserted text, or prevent it from being inserted entirely.
// Listener function arguments:
// 	newText string	the new text to insert
// 	newTextLength int	the length of the new text, in bytes, or -1 if new_text is nul-terminated
// 	position interface{}	the position, in characters, at which to insert the new text. this is an in-out parameter.  After the signal emission is finished, it should point after the newly inserted text.
const SignalInsertText cdk.Signal = "insert-text"
