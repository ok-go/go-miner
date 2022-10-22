package point

import "fmt"

type Point struct{ X, Y int }

func P(x, y int) Point {
	return Point{x, y}
}

func (p Point) String() string {
	return fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
}

func Ring(p Point, lvl int) []Point {
	if lvl <= 0 {
		return []Point{p}
	}
	var result []Point
	for y := p.Y - lvl; y < p.Y+lvl; y++ {
		result = append(result, Point{p.X - lvl, y})
	}
	for x := p.X - lvl; x < p.X+lvl; x++ {
		result = append(result, Point{x, p.Y + lvl})
	}
	for y := p.Y + lvl; y > p.Y-lvl; y-- {
		result = append(result, Point{p.X + lvl, y})
	}
	for x := p.X + lvl; x > p.X-lvl; x-- {
		result = append(result, Point{x, p.Y - lvl})
	}

	return result
}
