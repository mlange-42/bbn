package tui

type Color uint8

const (
	Black Color = iota
	Maroon
	Green
	Olive
	Navy
	Purple
	Teal
	Silver
	Gray
	Red
	Lime
	Yellow
	Blue
	Fuchsia
	Aqua
	White
	ColorsEnd
)

var ColorTags = []string{
	"[black]",
	"[maroon]",
	"[green]",
	"[olive]",
	"[navy]",
	"[purple]",
	"[teal]",
	"[silver]",
	"[gray]",
	"[red]",
	"[lime]",
	"[yellow]",
	"[blue]",
	"[fuchsia]",
	"[aqua]",
	"[-]",
}

var ColorNames = []string{
	"black",
	"maroon",
	"green",
	"olive",
	"navy",
	"purple",
	"teal",
	"silver",
	"gray",
	"red",
	"lime",
	"yellow",
	"blue",
	"fuchsia",
	"aqua",
	"white",
}

var NamedColors = map[string]Color{
	"black":   Black,
	"maroon":  Maroon,
	"green":   Green,
	"olive":   Olive,
	"navy":    Navy,
	"purple":  Purple,
	"teal":    Teal,
	"silver":  Silver,
	"gray":    Gray,
	"red":     Red,
	"lime":    Lime,
	"yellow":  Yellow,
	"blue":    Blue,
	"fuchsia": Fuchsia,
	"aqua":    Aqua,
	"white":   White,
}
