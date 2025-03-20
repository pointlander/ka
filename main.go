// Copyright 2025 The KA Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"math"
	"math/rand"
	"os"

	"github.com/pointlander/compress"
)

var states = [8][2]int{
	{1, 1},
	{1, 0},
	{1, -1},
	{0, -1},
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, 1},
}

func press(u [8 * 8]byte, x0, y0 int) int {
	trace := []byte{u[8*y0+x0]}
	add := func(x, y int) {
		if x < 0 || y < 0 || x >= 8 || y >= 8 {
			trace = append(trace, 0)
			return
		}
		trace = append(trace, u[8*y+x])
	}
	for r := 1; r < 8; r++ {
		x, y, dx, dy := r-1, 0, 1, 1
		err := dx - (r * 2)

		for x >= y {
			add(x0+x, y0+y)
			add(x0+y, y0+x)
			add(x0-y, y0+x)
			add(x0-x, y0+y)
			add(x0-x, y0-y)
			add(x0-y, y0-x)
			add(x0+y, y0-x)
			add(x0+x, y0-y)

			if err <= 0 {
				y++
				err += dy
				dy += 2
			}
			if err > 0 {
				x--
				dx += 2
				err += dx - (r * 2)
			}
		}
	}
	var buffer bytes.Buffer
	compress.Mark1Compress1(trace, &buffer)
	return buffer.Len()
}

func main() {
	rng := rand.New(rand.NewSource(1))
	u := [8 * 8]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 255, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 255, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 255, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	images := &gif.GIF{}
	var palette = []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}
	const (
		Size  = 8
		Scale = 25
	)
	type Coord struct {
		X, Y int
		D    int
	}
	iterations := 128
	for step := 0; step < iterations; step++ {
		best, coords := 0, make([]Coord, 0, 8)
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if u[y*8+x] > 0 {
					best += press(u, x, y)
					coords = append(coords, Coord{
						X: x,
						Y: y,
						D: rng.Intn(len(states)),
					})
				}
			}
		}
		for {
			fitness := 0
			v := u
			for _, coord := range coords {
				x, y := coord.X+states[coord.D][0], coord.Y+states[coord.D][1]
				if x < 0 {
					x = 8 + x
				}
				if y < 0 {
					y = 8 + y
				}
				if y >= 8 {
					y = y - 8
				}
				if x >= 8 {
					x = x - 8
				}
				v[coord.Y*8+coord.X], v[y*8+x] = v[y*8+x], v[coord.Y*8+coord.X]
			}
			for _, coord := range coords {
				x, y := coord.X+states[coord.D][0], coord.Y+states[coord.D][1]
				if x < 0 {
					x = 8 + x
				}
				if y < 0 {
					y = 8 + y
				}
				if y >= 8 {
					y = y - 8
				}
				if x >= 8 {
					x = x - 8
				}
				fitness += press(v, x, y)
			}
			if fitness <= best {
				u = v
				break
			}
			for i := range coords {
				coords[i].D = rng.Intn(len(states))
			}
		}
		verse := image.NewPaletted(image.Rect(0, 0, Size*Scale, Size*Scale), palette)
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				b := u[i*8+j]
				if b != 0 {
					xx, yy := j*Scale, i*Scale
					for x := 0; x < Scale; x++ {
						for y := 0; y < Scale; y++ {
							var dx, dy float32 = Scale/2 - float32(x), Scale/2 - float32(y)
							d := 2 * float32(math.Sqrt(float64(dx*dx+dy*dy))) / Scale
							if d < 1 {
								verse.Set(xx+x, yy+y, color.RGBA{0xff, 0xff, 0xff, 0xff})
							}
						}
					}
				}
			}
		}
		for x := 0; x < int(float64(step)*Size*Scale/float64(iterations)); x++ {
			for y := Size*Scale - 10; y < Size*Scale; y++ {
				verse.Set(x, y, color.RGBA{0, 0, 0xff, 0xff})
			}
		}
		images.Image = append(images.Image, verse)
		images.Delay = append(images.Delay, 10)
		fmt.Println(step)
	}

	out, err := os.Create("ka.gif")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	err = gif.EncodeAll(out, images)
	if err != nil {
		panic(err)
	}
}
