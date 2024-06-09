package tui

var (
	BorderH  = []rune{'─', '═'}
	BorderV  = []rune{'│', '║'}
	BorderNE = []rune{'┐', '╗'}
	BorderNW = []rune{'┌', '╔'}
	BorderSE = []rune{'┘', '╝'}
	BorderSW = []rune{'└', '╚'}
)

const (
	Empty          = ' '
	Shade          = '░'
	Full           = '█'
	SelectionStart = '['
	SelectionEnd   = ']'
	EvidenceStart  = '+'
	EvidenceEnd    = '+'
	ArrowUp        = '^'
	ArrowDown      = 'v'
	ArrowLeft      = '<'
	ArrowRight     = '>'
)

var Partial = []rune{
	'▏', '▎', '▍', '▌', '▋', '▊', '▉',
}
