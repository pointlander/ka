// Copyright 2025 The KA Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
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
	X [2]int
	R float64
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
					X: [2]int{x - 4, y - 4},
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

// K computes the kolmogorov complexity
func K(u [Size * Size]byte, x0, y0 int) int {
	trace := []byte{}
	for _, v := range Circle {
		x, y := v.X[0]+x0, v.X[1]+y0
		for x < 0 {
			x = Size + x
		}
		for y < 0 {
			y = Size + y
		}
		for x >= Size {
			x = x - Size
		}
		for y >= Size {
			y = y - Size
		}
		trace = append(trace, u[Size*y+x])
	}

	var buffer bytes.Buffer
	Mark1Compress1(trace, &buffer)
	return buffer.Len()
}

// U is a universe simulator
func U(filename string, next func(v [Size * Size]byte, rng *rand.Rand) [Size * Size]byte) {
	rng := rand.New(rand.NewSource(1))
	u := [Size * Size]byte{}
	for i := 0; i < 9; i++ {
		u[rng.Intn(Size*Size)] = 255
	}

	images := &gif.GIF{}
	var palette = []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}

	for step := 0; step < Iterations; step++ {
		u = next(u, rng)
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
		fmt.Println(step)
	}

	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	err = gif.EncodeAll(out, images)
	if err != nil {
		panic(err)
	}

}

// Capitalistic mode
func Capitalistic() {
	U("capitalistic.gif", func(u [Size * Size]byte, rng *rand.Rand) [Size * Size]byte {
		for {
			ax, ay, bx, by := rng.Intn(Size), rng.Intn(Size), rng.Intn(Size), rng.Intn(Size)
			v := u
			v[ay*Size+ax], v[by*Size+bx] = v[by*Size+bx], v[ay*Size+ax]
			aBefore := K(u, ax, ay)
			bBefore := K(u, bx, by)
			aAfter := K(v, ax, ay)
			bAfter := K(v, bx, by)
			if aAfter < aBefore && bAfter < bBefore {
				u = v
				break
			}
		}
		return u
	})
}

// Communistic model
func Communistic() {
	U("communistic.gif", func(u [Size * Size]byte, rng *rand.Rand) [Size * Size]byte {
		for j := 0; j < 256; j++ {
			ax, ay, bx, by := rng.Intn(Size), rng.Intn(Size), rng.Intn(Size), rng.Intn(Size)
			v := u
			v[ay*Size+ax], v[by*Size+bx] = v[by*Size+bx], v[ay*Size+ax]
			before := 0
			for y := 0; y < Size; y++ {
				for x := 0; x < Size; x++ {
					before += K(u, x, y)
				}
			}
			after := 0
			for y := 0; y < Size; y++ {
				for x := 0; x < Size; x++ {
					after += K(v, x, y)
				}
			}
			if after < before {
				u = v
				break
			}
		}
		return u
	})
}

var (
	// FlagCapitalistic capitalistic model
	FlagCapitalistic = flag.Bool("capitalistic", false, "capitalistic mode")
	// FlagCommunistic communistic model
	FlagCommunistic = flag.Bool("communistic", false, "communistic model")
)

func main() {
	flag.Parse()

	if *FlagCapitalistic {
		Capitalistic()
		return
	}

	if *FlagCommunistic {
		Communistic()
		return
	}

	iris := Load()
	fmt.Println(len(iris))
}
