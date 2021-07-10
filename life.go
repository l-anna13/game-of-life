package main

import (
	"bytes"
	"fmt"
	"time"
)

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[x][y] = b
}

// Alive reports whether the specified cell is alive.
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	s := f.s[x][y]
	return s
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	// if lees then 2 or more than 3 alive make current "dead" cell
	if alive < 2 || alive > 3 {
		return false
	} else if !f.s[x][y] && alive == 3 { // if current is dead and exactly 3 around are alive then current becomes alive
		return true
	} else if alive == 2 || alive == 3 { // if 2 or 3 are alive current is alive
		return true
	} else { // else dead cell
		return false
	}
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b *Field
	w, h int
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for i := 0; i < (w * h); i++ {
		for j := 0; j < w; j++ {
			for k := 0; k < h; k++ {
				// "Glider" as a starting point in a middle of a 25*25 universe
				if j == 12 && k == 11 {
					a.Set(j, k, true)
				} else if j == 13 && k == 12 {
					a.Set(j, k, true)
				} else if (j == 11 && k == 13) ||
					(j == 12 && k == 13) ||
					(j == 13 && k == 13) {
					a.Set(j, k, true)
				} else {
					a.Set(j, k, false)
				}
			}
		}
	}
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			s := l.a.Next(x, y)
			l.b.Set(x, y, s)
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}

// String returns the game board as a string.
func (l *Life) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			b := byte('-')
			if l.a.Alive(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func main() {
	l := NewLife(25, 25)
	for i := 0; i < 5; i++ {
		fmt.Print("\x0c", l) // Clear screen and print field.
		l.Step()
		time.Sleep(time.Second / 30)
	}
}
