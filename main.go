// Copyright 2025 The Grid Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	//"gonum.org/v1/plot"
	//"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/vg"
	//"gonum.org/v1/plot/vg/draw"
	"github.com/pointlander/compress"
)

func AreRatiosEqual(a, b, c, d int) bool {
	if b == 0 || d == 0 {
		return false
	}
	return a*d == b*c
}

var (
	// FlagGenerate generate cas
	FlagGenerate = flag.Bool("generate", false, "generate cas")
)

// Generate generate cas
func Generate() {
	type Ratio struct {
		Ratio float64
		One   int
		Zero  int
		Found bool
	}
	ratio := make([]Ratio, 256)
	const size = 8 * 1024
	for rule := range 256 {
		img := image.NewGray(image.Rect(0, 0, size, size/2))
		//points := make(plotter.XYs, 0, 8)
		grid := make([]byte, size)
		grid[size/2] = 1
		for iteration := range size / 2 {
			for key, value := range grid {
				if value > 0 {
					value = 0
				} else {
					value = 255
				}
				img.SetGray(key, iteration, color.Gray{Y: byte(value)})
			}
			next := make([]byte, len(grid))
			for cell := 1; cell < len(grid)-1; cell++ {
				state := grid[cell-1]*4 + grid[cell]*2 + grid[cell+1]*1
				next[cell] = byte((rule >> state) & 1)
			}
			grid = next
			/*one, zero := 0, 0
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
			points = append(points, plotter.XY{X: float64(iteration), Y: r})*/
		}

		output, err := os.Create(fmt.Sprintf("plots/ca%d.png", rule))
		if err != nil {
			panic(err)
		}
		defer output.Close()

		err = png.Encode(output, img)
		if err != nil {
			panic(err)
		}

		/*p := plot.New()

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
		}*/
		var buffer bytes.Buffer
		compress.Mark1Compress16(grid, &buffer)
		zero := buffer.Len()
		one := size
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

func main() {
	flag.Parse()

	if *FlagGenerate {
		Generate()
		return
	}

	rule := 110
	rng := rand.New(rand.NewSource(1))
	target := make([]byte, 8+2)
	t := byte(0)
	for i := range 8 {
		bit := byte(rng.Intn(2))
		target[i+1] = bit
		t |= bit << i
	}
	for range 4 {
		fmt.Printf("target %v\n", target)
		next := [][]byte{}
		for guess := range 256 {
			g := make([]byte, 8+2)
			for i := range 8 {
				g[i+1] = byte((guess >> i) & 1)
			}
			infer := make([]byte, 8+2)
			for cell := 1; cell < len(g)-1; cell++ {
				state := g[cell-1]*4 + g[cell]*2 + g[cell+1]*1
				infer[cell] = byte((rule >> state) & 1)
			}
			equals := true
			for key, value := range target {
				if value != infer[key] {
					equals = false
					break
				}
			}
			if equals {
				fmt.Printf("%v\n", g)
				next = append(next, g)
			}
		}
		if len(next) == 0 {
			break
		}
		target = next[0]
	}
}
