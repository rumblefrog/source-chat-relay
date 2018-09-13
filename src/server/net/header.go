package net

// Header - Part of every payload (length may be 0 or ommited)
type Header struct {
	OPCode OPCODE
	Length uint16
}
