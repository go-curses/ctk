package ctk

type Alignable interface {
	SetAlignment(xAlign float64, yAlign float64)
	GetAlignment() (xAlign float64, yAlign float64)
}
