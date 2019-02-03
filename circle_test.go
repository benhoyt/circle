package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestText(t *testing.T) {
	buf := &bytes.Buffer{}
	_ = drawCircleText(3, buf)
	expected := `
  ###  
 #   # 
#     #
#     #
#     #
 #   # 
  ###  
`
	if "\n"+buf.String() != expected {
		t.Errorf("not equal")
	}
}

type point struct {
	x int
	y int
}

func TestCompare(t *testing.T) {
	for r := 0; r < 100; r++ {
		pointsSqrt := make(map[point]bool)
		drawCircleSqrt(r, func(x, y int) {
			pointsSqrt[point{x, y}] = true
		})
		pointsInt := make(map[point]bool)
		drawCircleInt(r, func(x, y int) {
			pointsInt[point{x, y}] = true
		})
		if !reflect.DeepEqual(pointsSqrt, pointsInt) {
			t.Errorf("radius %d not equal", r)
		}
	}
}
func BenchmarkSqrt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		drawCircleSqrt(100, func(x, y int) {})
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		drawCircleInt(100, func(x, y int) {})
	}
}
