package net

// OPCODE -
type OPCODE int

// MessageFrame - Message payloads
// PingFrame - Simple & minimal ping request
const (
	MessageFrame OPCODE = iota
	PingFrame    OPCODE = iota
)
