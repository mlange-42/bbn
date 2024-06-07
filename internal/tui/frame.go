package tui

var (
	BorderNone = []rune{' ', ' '}
	BorderH    = []rune{'─', '═'}
	BorderV    = []rune{'│', '║'}
	BorderNE   = []rune{'┐', '╗'}
	BorderNW   = []rune{'┌', '╔'}
	BorderSE   = []rune{'┘', '╝'}
	BorderSW   = []rune{'└', '╚'}
)

const (
	Shade = '░'
	Full  = '█'
)

var Partial = []rune{
	'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█',
}
