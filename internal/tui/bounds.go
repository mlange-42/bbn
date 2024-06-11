package tui

type Bounds struct {
	X int
	Y int
	W int
	H int
}

func NewBounds(x, y, w, h int) Bounds {
	return Bounds{
		X: x, Y: y, W: w, H: h,
	}
}

func (b *Bounds) Contains(x, y int) bool {
	return x >= b.X && x < b.X+b.W && y >= b.Y && y < b.Y+b.H
}

func (b *Bounds) Centroid() (int, int) {
	return b.X + b.W/2, b.Y + b.H/2
}

func (b *Bounds) Extend(other *Bounds) {
	if other.X < b.X {
		b.W += b.X - other.X
		b.X = other.X
	}
	if other.Y < b.Y {
		b.H += b.Y - other.Y
		b.Y = other.Y
	}
	if other.X+other.W > b.X+b.W {
		b.W = other.X + other.W - b.X
	}
	if other.Y+other.H > b.Y+b.H {
		b.H = other.Y + other.H - b.Y
	}
}
