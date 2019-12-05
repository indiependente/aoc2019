package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindClosestIntersection_Manhattan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		input        io.Reader
		wantPoint    Point
		wantDistance float64
	}{
		{
			name: "small example - manhattan distance",
			input: strings.NewReader(`R8,U5,L5,D3
U7,R6,D4,L4`),
			wantPoint:    Point{3, 3},
			wantDistance: 6,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			point, _, dist := FindClosestIntersection(tt.input)
			assert.Equal(t, tt.wantPoint, point)
			assert.Equal(t, tt.wantDistance, dist)
		})
	}
}
func TestFindClosestIntersection_Steps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     io.Reader
		wantPoint Point
		wantSteps int
	}{
		{
			name: "small example - steps distance",
			input: strings.NewReader(`R8,U5,L5,D3
U7,R6,D4,L4`),
			wantPoint: Point{6, 8},
			wantSteps: 30,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, weighted, _ := FindClosestIntersection(tt.input)
			assert.Equal(t, tt.wantPoint, weighted.p)
			assert.Equal(t, tt.wantSteps, weighted.w)
		})
	}
}
