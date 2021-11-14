package ctk

// An Alignable Widget is one that implements the SetAlignment and GetAlignment
// methods for adjusting the positioning of the Widget. The Misc and Alignment
// types are the primary ones implementing this interface.
type Alignable interface {
	SetAlignment(xAlign float64, yAlign float64)
	GetAlignment() (xAlign float64, yAlign float64)
}
