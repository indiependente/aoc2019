package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	x, y int
}
type WeightedIntersection struct {
	p Point
	w int
}
type Segment struct {
	a, b Point
}
type Wire struct {
	segments []Segment
	symbol   rune
}

// TODO: switch from Grid to Segment based
type Grid map[int]map[int]rune

func (g Grid) Get(x, y int) rune {
	r, _ := g[x][y]
	return r
}
func (g Grid) Marked(x, y int) bool {
	_, ok := g[x][y]
	return ok
}
func (g Grid) Mark(x, y int, r rune) {
	g[x][y] = r
}
func (g Grid) MarkRow(row, from, to int, r rune) []WeightedIntersection {
	var intersections []WeightedIntersection
	if from > to {
		from, to = to, from
	}
	_, ok := g[row]
	if !ok {
		g[row] = make(map[int]rune)
	}
	for i := from; i <= to; i++ {
		v, ok := g[row][i]
		if ok {
			if v != r {
				fmt.Printf("intersection [%d %d] steps = %d\n", row, i, int(math.Abs(float64(from-i))))
				intersections = append(intersections, WeightedIntersection{
					p: Point{row, i},
					w: int(math.Abs(float64(from - i))),
				})
				continue
			}
		}
		g[row][i] = r
	}
	return intersections
}
func (g Grid) MarkCol(col, from, to int, r rune) []WeightedIntersection {
	var intersections []WeightedIntersection
	if from > to {
		from, to = to, from
	}
	for i := from; i <= to; i++ {
		_, ok := g[i]
		if !ok {
			g[i] = make(map[int]rune)
		}
		v, ok := g[i][col]
		if ok {
			if v != r {
				fmt.Printf("intersection [%d %d] steps = %d\n", i, col, int(math.Abs(float64(from-i))))
				intersections = append(intersections, WeightedIntersection{
					p: Point{i, col},
					w: int(math.Abs(float64(from - i))),
				})
				continue
			}
		}
		g[i][col] = r
	}
	return intersections
}

func manhattanDistance(a Point, b Point) float64 {
	return math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y))
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	closestPoint, cheapest, dist := FindClosestIntersection(f)
	fmt.Printf("Closest point: [%d %d]\nManhattan distance = %f\n", closestPoint.x, closestPoint.y, dist)
	fmt.Printf("Cheapest point (steps): [%d %d]\nNo. Steps = %d\n", cheapest.p.x, cheapest.p.y, cheapest.w)
}

func FindClosestIntersection(r io.Reader) (Point, WeightedIntersection, float64) {
	var (
		closestPoint = Point{
			x: math.MaxInt32,
			y: math.MaxInt32,
		}
		origin = Point{
			x: 0,
			y: 0,
		}
		shortestDistance = manhattanDistance(closestPoint, origin)
		cheapestPoint    = WeightedIntersection{
			p: Point{
				x: math.MaxInt32,
				y: math.MaxInt32,
			},
			w: math.MaxInt32,
		}
	)
	g := make(Grid)
	s := bufio.NewScanner(r)
	for s.Scan() {
		wire := randomRune()
		line := s.Text()
		directions := strings.Split(line, ",")
		x, y := 0, 0
		for _, s := range directions {
			d := s[0]
			cells, err := strconv.Atoi(s[1:])
			if err != nil {
				log.Fatal(err)
			}
			switch d {
			case 'R':
				fmt.Printf("%c: %c %d => [%d, %d]\n", wire, d, cells, x, y+cells)
				intersections := g.MarkRow(x, y, y+cells, wire)
				closestPoint, cheapestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, cheapestPoint, intersections)
				y += cells
			case 'L':
				fmt.Printf("%c: %c %d => [%d, %d]\n", wire, d, cells, x, y-cells)
				intersections := g.MarkRow(x, y, y-cells, wire)
				closestPoint, cheapestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, cheapestPoint, intersections)
				y -= cells
			case 'U':
				fmt.Printf("%c: %c %d => [%d, %d]\n", wire, d, cells, x+cells, y)
				intersections := g.MarkCol(y, x, x+cells, wire)
				closestPoint, cheapestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, cheapestPoint, intersections)
				x += cells
			case 'D':
				fmt.Printf("%c: %c %d => [%d, %d]\n", wire, d, cells, x-cells, y)
				intersections := g.MarkCol(y, x, x-cells, wire)
				closestPoint, cheapestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, cheapestPoint, intersections)
				x -= cells
			}
		}
	}
	return closestPoint, cheapestPoint, shortestDistance
}

func updateClosest(closestPoint Point, shortestDistance float64, cheapestPoint WeightedIntersection, intersections []WeightedIntersection) (Point, WeightedIntersection, float64) {
	var (
		origin = Point{0, 0}
	)
	if len(intersections) > 0 {
		fmt.Println(intersections)
	}
	for _, isect := range intersections {
		if isect.p == origin {
			continue
		}
		dist := manhattanDistance(isect.p, origin)
		if dist < shortestDistance {
			closestPoint = isect.p
			shortestDistance = dist
		}
		if isect.w < cheapestPoint.w {
			cheapestPoint = isect
			fmt.Printf("cheap intersection point: [%d %d] => %d\n", isect.p.x, isect.p.y, isect.w)
		}
	}
	return closestPoint, cheapestPoint, shortestDistance
}

func randomRune() rune {
	rand.Seed(time.Now().UnixNano())
	return rune(rand.Intn(26) + 65)
}
