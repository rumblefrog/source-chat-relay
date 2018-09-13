package net

type OPCODE int

const (
	MessageFrame OPCODE = iota
	PingFrame    OPCODE = iota
)
