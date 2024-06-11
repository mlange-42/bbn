package tui

type Color uint8

const (
	White Color = iota
	Yellow
	Red
	Blue
	Green
)

var ColorTags = []string{
	"[-]",
	"[yellow]",
	"[red]",
	"[blue]",
	"[green]",
}
