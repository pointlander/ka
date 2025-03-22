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
	Iterations = 256
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
	X      [2]float64
	Avg    [2]float64
	Stddev [2]float64
	S      [2]float64
	R      float64
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
					X: [2]float64{float64(x - 4), float64(y - 4)},
					R: d,
				})
			}
		}
	}
	sort.Slice(Circle, func(i, j int) bool {
		if Circle[i].R < Circle[j].R {
			return true
		} else if Circle[i].R == Circle[j].R {
			if Circle[i].X[0] < Circle[j].X[0] {
				return true
			} else if Circle[i].X[0] == Circle[j].X[0] {
				return Circle[i].X[1] < Circle[j].X[1]
			}
		}
		return false
	})
}

func R(x float64) int {
	return int(math.Round(x))
}

// K computes the kolmogorov complexity
func K(qv [Size * Size]byte, u [Size * Size]byte, x0, y0 float64) int {
	trace := []byte{}
	for _, v := range Circle {
		x, y := v.X[0]+x0, v.X[1]+y0
		for R(x) < 0 {
			x = Size + x
		}
		for R(y) < 0 {
			y = Size + y
		}
		for R(x) >= Size {
			x = x - Size
		}
		for R(y) >= Size {
			y = y - Size
		}
		if qv[Size*R(y)+R(x)] > 0 || u[Size*R(y)+R(x)] > 0 {
			trace = append(trace, 255)
		} else {
			trace = append(trace, 0)
		}
	}

	var buffer bytes.Buffer
	Mark1Compress1(trace, &buffer)
	return buffer.Len()
}

func main() {
	rng := rand.New(rand.NewSource(1))
	u := [Size * Size]byte{}
	qv := [Size * Size]byte{}
	for i := 0; i < 9; i++ {
		u[rng.Intn(Size*Size)] = 255
	}
	images := &gif.GIF{}
	var palette = []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}
	best, coords := 0.0, make([]Coord, 0, 8)
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			if u[y*Size+x] > 0 {
				best += float64(K(qv, u, float64(x), float64(y)))
				coords = append(coords, Coord{
					X:      [2]float64{float64(x), float64(y)},
					Avg:    [2]float64{0, 0},
					Stddev: [2]float64{float64(4), float64(4)},
				})
			}
		}
	}
	for step := 0; step < Iterations; step++ {
		for j := 0; j < 256; j++ {
			qv := [Size * Size]byte{}
			for i := 0; i < 3; i++ {
				qv[rng.Intn(Size*Size)] = 255
			}
			fitness := 0.0
			v := u
			for i, coord := range coords {
				x, y := coord.X[0]+rng.NormFloat64()*coord.Stddev[0]+coord.Avg[0],
					coord.X[1]+rng.NormFloat64()*coord.Stddev[1]+coord.Avg[1]
				for R(x) < 0 {
					x = Size + x
				}
				for R(y) < 0 {
					y = Size + y
				}
				for R(y) >= Size {
					y = y - Size
				}
				for R(x) >= Size {
					x = x - Size
				}
				coords[i].S[0], coords[i].S[1] = x, y
				v[R(coord.X[1])*Size+R(coord.X[0])], v[R(y)*Size+R(x)] =
					v[R(y)*Size+R(x)], v[R(coord.X[1])*Size+R(coord.X[0])]
			}
			for _, coord := range coords {
				fitness += float64(K(qv, v, coord.S[0], coord.S[1]))
			}
			if fitness <= best {
				for i := range coords {
					coords[i].X = coords[i].S
				}
				u = v
				best = fitness
				break
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
