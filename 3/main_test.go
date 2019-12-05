package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindClosestIntersection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		input        io.Reader
		wantPoint    Point
		wantDistance float64
	}{
		{
			name: "small example",
			input: strings.NewReader(`R8,U5,L5,D3
U7,R6,D4,L4`),
			wantPoint:    Point{3, 3},
			wantDistance: 6,
		},
		//		{
		//			name: "happy path",
		//			input: strings.NewReader(`R75,D30,R83,U83,L12,D49,R71,U7,L72
		//U62,R66,U55,R34,D71,R55,D58,R83`),
		//			wantPoint:    Point{},
		//			wantDistance: 159,
		//		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			point, dist := FindClosestIntersection(tt.input)
			assert.Equal(t, tt.wantPoint, point)
			assert.Equal(t, tt.wantDistance, dist)
		})
	}
}
