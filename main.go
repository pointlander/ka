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
	"io"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/pointlander/compress"
)

const (
	// Size is the size of the universe
	Size = 9
	// Scale is the scale factor for rendering
	Scale = 25
	// Iterations is the number of iterations
	Iterations = 1024
)

func Mark1Compress1(input []byte, output io.Writer) {
	data, channel := make([]byte, len(input)), make(chan []byte, 1)
	copy(data, input)
	channel <- data
	close(channel)
	compress.BijectiveBurrowsWheelerCoder(channel).MoveToFrontCoder().FilteredAdaptiveBitCoder().CodeBit(output)
}

// Coord is a coordinate
type Coord struct {
	X, Y int
	D    int
	R    float64
}

// States are the next states
var States = [8]Coord{
	{1, 1, 0, 0},
	{1, 0, 0, 0},
	{1, -1, 0, 0},
	{0, -1, 0, 0},
	{-1, -1, 0, 0},
	{-1, 0, 0, 0},
	{-1, 1, 0, 0},
	{0, 1, 0, 0},
}

// Circle is a circle
var Circle []Coord

func init() {
	for x := 0; x < Size; x++ {
		for y := 0; y < Size; y++ {
			dx, dy := 5-x, 5-y
			d := math.Sqrt(float64(dx*dx + dy*dy))
			if d <= 4 {
				Circle = append(Circle, Coord{
					X: x - 4,
					Y: y - 4,
					R: d,
				})
			}
		}
	}
	sort.Slice(Circle, func(i, j int) bool {
		if Circle[i].R < Circle[j].R {
			return true
		} else if Circle[i].R == Circle[j].R {
			if Circle[i].X < Circle[j].X {
				return true
			} else if Circle[i].X == Circle[j].X {
				return Circle[i].Y < Circle[j].Y
			}
		}
		return false
	})
}

// K computes the kolmogorov complexity
func K(u [Size * Size]byte, x0, y0 int) int {
	trace := []byte{}
	for _, v := range Circle {
		x, y := v.X+x0, v.Y+y0
		if x < 0 {
			x = Size + x
		}
		if y < 0 {
			y = Size + y
		}
		if y >= Size {
			y = y - Size
		}
		if x >= Size {
			x = x - Size
		}
		trace = append(trace, u[Size*y+x])
	}

	var buffer bytes.Buffer
	Mark1Compress1(trace, &buffer)
	return buffer.Len()
}

func main() {
	rng := rand.New(rand.NewSource(1))
	u := [Size * Size]byte{}
	for i := 0; i < 3; i++ {
		u[rng.Intn(Size*Size)] = 255
	}
	images := &gif.GIF{}
	var palette = []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}
	for step := 0; step < Iterations; step++ {
		best, coords := 00., make([]Coord, 0, 8)
		for y := 0; y < Size; y++ {
			for x := 0; x < Size; x++ {
				if u[y*Size+x] > 0 {
					best += float64(K(u, x, y))
					coords = append(coords, Coord{
						X: x,
						Y: y,
						D: rng.Intn(len(States)),
					})
				}
			}
		}
		for {
			fitness := 0.0
			v := u
			for _, coord := range coords {
				x, y := coord.X+States[coord.D].X, coord.Y+States[coord.D].Y
				if x < 0 {
					x = Size + x
				}
				if y < 0 {
					y = Size + y
				}
				if y >= Size {
					y = y - Size
				}
				if x >= Size {
					x = x - Size
				}
				v[coord.Y*Size+coord.X], v[y*Size+x] = v[y*Size+x], v[coord.Y*Size+coord.X]
			}
			for _, coord := range coords {
				x, y := coord.X+States[coord.D].X, coord.Y+States[coord.D].Y
				fitness += float64(K(v, x, y))
			}
			if fitness <= best {
				u = v
				break
			}
			for i := range coords {
				coords[i].D = rng.Intn(len(States))
			}
		}
		verse := image.NewPaletted(image.Rect(0, 0, Size*Scale, Size*Scale), palette)
		for i := 0; i < Size; i++ {
			for j := 0; j < Size; j++ {
				b := u[i*Size+j]
				if b != 0 {
					xx, yy := j*Scale, i*Scale
					for x := 0; x < Scale; x++ {
						for y := 0; y < Scale; y++ {
							dx, dy := Scale/2-float64(x), Scale/2-float64(y)
							d := 2 * math.Sqrt(dx*dx+dy*dy) / Scale
							if d < 1 {
								verse.Set(xx+x, yy+y, color.RGBA{0xff, 0xff, 0xff, 0xff})
							}
						}
					}
				}
			}
		}
		for x := 0; x < int(float64(step)*Size*Scale/float64(Iterations)); x++ {
			for y := Size*Scale - 10; y < Size*Scale; y++ {
				verse.Set(x, y, color.RGBA{0, 0, 0xff, 0xff})
			}
		}
		images.Image = append(images.Image, verse)
		images.Delay = append(images.Delay, 20)
		fmt.Println(step, best)
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
