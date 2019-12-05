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
type Segment struct {
	a, b Point
}
type Wire struct {
	segments []Segment
	symbol   rune
}
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
func (g Grid) MarkRow(row, from, to int, r rune) []Point {
	var intersections []Point
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
				intersections = append(intersections, Point{row, i})
				continue
			}
		}
		g[row][i] = r
	}
	return intersections
}
func (g Grid) MarkCol(col, from, to int, r rune) []Point {
	var intersections []Point
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
				intersections = append(intersections, Point{i, col})
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
	closestPoint, dist := FindClosestIntersection(f)
	fmt.Printf("Closest point: [%d %d]\n", closestPoint.x, closestPoint.y)
	fmt.Printf("Manhattan distance = %f\n", dist)
}

func FindClosestIntersection(r io.Reader) (Point, float64) {
	var (
		closestPoint = Point{
			x: int(math.Inf(0)),
			y: int(math.Inf(0)),
		}
		origin = Point{
			x: 0,
			y: 0,
		}
		shortestDistance = manhattanDistance(closestPoint, origin)
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
				closestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, intersections)
				y += cells
			case 'L':
				fmt.Printf("%c: %c %d => [%d, %d]\n", wire, d, cells, x, y-cells)
				intersections := g.MarkRow(x, y, y-cells, wire)
				closestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, intersections)
				y -= cells
			case 'U':
				fmt.Printf("%c: %c %d => [%d, %d]\n", wire, d, cells, x+cells, y)
				intersections := g.MarkCol(y, x, x+cells, wire)
				closestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, intersections)
				x += cells
			case 'D':
				fmt.Printf("%c: %c %d => [%d, %d]\n", wire, d, cells, x-cells, y)
				intersections := g.MarkCol(y, x, x-cells, wire)
				closestPoint, shortestDistance = updateClosest(closestPoint, shortestDistance, intersections)
				x -= cells
			}
		}
	}
	return closestPoint, shortestDistance
}

func updateClosest(closestPoint Point, shortestDistance float64, intersections []Point) (Point, float64) {
	var (
		origin = Point{0, 0}
	)
	if len(intersections) > 0 {
		fmt.Println(intersections)
	}
	for _, isect := range intersections {
		if isect == origin {
			continue
		}
		dist := manhattanDistance(isect, origin)
		if dist < shortestDistance {
			closestPoint = isect
			shortestDistance = dist
			fmt.Printf("intersection point: [%d %d] => %f\n", isect.x, isect.y, dist)
		}
	}
	return closestPoint, shortestDistance
}

func randomRune() rune {
	rand.Seed(time.Now().UnixNano())
	return rune(rand.Intn(26) + 65)
}
