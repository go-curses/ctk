package enums

import (
	cbits "github.com/go-curses/cdk/lib/bits"
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
