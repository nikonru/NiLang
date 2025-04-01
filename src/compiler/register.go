package compiler

type register = string

const (
	AX = "AX"
	BX = "BX"
	CX = "CX" // flag for bot's memory being ready for reading
	DX = "DX" // bot's memory

	SD = "SD"
	MD = "MD"
	EN = "EN"
	AG = "AG"
)
