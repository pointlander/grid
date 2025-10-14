// Copyright 2025 The Grid Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func AreRatiosEqual(a, b, c, d int) bool {
	if b == 0 || d == 0 {
		return false
	}
	return a*d == b*c
}

func main() {
	type Ratio struct {
		Ratio float64
		One   int
		Zero  int
		Found bool
	}
	ratio := make([]Ratio, 256)
	const size = 8 * 1024
	for rule := range 256 {
		points := make(plotter.XYs, 0, 8)
		grid := make([]byte, size)
		grid[size/2] = 1
		for iteration := range size / 2 {
			next := make([]byte, len(grid))
			for cell := 1; cell < len(grid)-1; cell++ {
				state := grid[cell-1]*4 + grid[cell]*2 + grid[cell+1]*1
				next[cell] = byte((rule >> state) & 1)
			}
			grid = next
			one, zero := 0, 0
			for _, value := range grid {
				if value == 0 {
					zero++
				} else {
					one++
				}
			}
			r := 0.0
			if one != 0 {
				r = float64(zero) / float64(one)
			}
			points = append(points, plotter.XY{X: float64(iteration), Y: r})
		}

		p := plot.New()

		p.Title.Text = "iteration vs ratio"
		p.X.Label.Text = "iteration"
		p.Y.Label.Text = "ratio"

		scatter, err := plotter.NewScatter(points)
		if err != nil {
			panic(err)
		}
		scatter.GlyphStyle.Radius = vg.Length(1)
		scatter.GlyphStyle.Shape = draw.CircleGlyph{}
		p.Add(scatter)

		err = p.Save(8*vg.Inch, 8*vg.Inch, fmt.Sprintf("plots/%d.png", rule))
		if err != nil {
			panic(err)
		}

		one, zero := 0, 0
		for _, value := range grid {
			if value == 0 {
				zero++
			} else {
				one++
			}
		}
		ratio[rule] = Ratio{
			Ratio: float64(zero) / float64(one),
			One:   one,
			Zero:  zero,
		}
	}
	for i, r := range ratio {
		if ratio[i].Found {
			continue
		}
		fmt.Println(i, r)
		for key, value := range ratio {
			if key == i || math.IsInf(r.Ratio, 0) {
				continue
			}
			if AreRatiosEqual(r.Zero, r.One, value.Zero, value.One) {
				ratio[key].Found = true
				fmt.Println(" ", key, value)
			}
		}
	}
}
