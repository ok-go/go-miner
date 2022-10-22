package board

import (
	"go-miner/pkg/point"
	"math/rand"
)

type Cell struct {
	MineCount int
	HasBomb   bool
	Opened    bool
	Selected  bool
	Marked    bool
}

type Board struct {
	W, H int
	data map[point.Point]*Cell
}

func NewBoard(w, h, mineCount int) *Board {
	if mineCount > w*h {
		panic("too much mines")
	}

	b := Board{w, h, make(map[point.Point]*Cell, w*h)}
	for i := 0; i < w*h; i++ {
		cell := &Cell{}
		b.data[point.P(i%w, i/w)] = cell
	}

	for i := 0; i < mineCount; i++ {
		p := point.P(rand.Intn(w), rand.Intn(h))
		if b.data[p].HasBomb {
			i--
			continue
		}
		b.data[p].HasBomb = true
		for _, n := range point.Ring(p, 1) {
			if cell, ok := b.data[n]; ok {
				cell.MineCount += 1
			}
		}
	}

	// open random space
	for gp, cell := range b.data {
		if !cell.HasBomb && cell.MineCount == 0 {
			b.RecursiveOpen(gp)
			break
		}
	}

	return &b
}

func (b *Board) RecursiveOpen(p point.Point) {
	if cell := b.At(p); cell != nil && !cell.Opened {
		cell.Opened = true
		if cell.MineCount == 0 && !cell.HasBomb {
			for _, n := range point.Ring(p, 1) {
				b.RecursiveOpen(n)
			}
		}
	}
}

func (b *Board) At(p point.Point) *Cell {
	if c, ok := b.data[p]; ok {
		return c
	}
	return nil
}

func (b *Board) Field() map[point.Point]*Cell {
	return b.data
}

func (b *Board) TryMarkNeighbours(p point.Point) (int, bool) {
	cell := b.At(p)
	if cell == nil {
		return 0, true
	}

	var nmarked []point.Point
	var marked = 0
	nps := point.Ring(p, 1)
	for _, np := range nps {
		if c := b.At(np); c != nil && !c.Opened {
			if c.Marked {
				marked += 1
			} else {
				nmarked = append(nmarked, np)
			}
		}
	}
	newMarked := 0
	for _, np := range nmarked {
		if c := b.At(np); c != nil {
			if marked+len(nmarked) == cell.MineCount {
				c.Marked = true
				newMarked += 1
			}
			if marked == cell.MineCount {
				b.RecursiveOpen(np)
				if c := b.At(np); c != nil && c.HasBomb {
					return newMarked, false
				}
			}
		}
	}
	return newMarked, true
}

func (b *Board) ShowAll() {
	for _, cell := range b.data {
		if !cell.HasBomb {
			cell.Opened = true
		} else {
			cell.Marked = true
		}
	}
}
