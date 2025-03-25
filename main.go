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
	"image/png"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/pointlander/compress"
	"github.com/pointlander/ka/vector"
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
type Circle []Coord

func NewCircle(size int) Circle {
	circle := make([]Coord, 0, 8)
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			dx, dy := (size/2+1)-x, (size/2+1)-y
			d := math.Sqrt(float64(dx*dx + dy*dy))
			if d <= 4 {
				circle = append(circle, Coord{
					X: [2]int{x - size/2, y - size/2},
					R: d,
				})
			}
		}
	}
	sort.Slice(circle, func(i, j int) bool {
		if circle[i].R < circle[j].R {
			return true
		} else if circle[i].R == circle[j].R {
			if circle[i].X[0] < circle[j].X[0] {
				return true
			} else if circle[i].X[0] == circle[j].X[0] {
				return circle[i].X[1] < circle[j].X[1]
			}
		}
		return false
	})
	return circle
}

// K computes the kolmogorov complexity
func (c Circle) K(size int, u []byte, x0, y0 int) int {
	trace := []byte{}
	for _, v := range c {
		x, y := v.X[0]+x0, v.X[1]+y0
		for x < 0 {
			x = size + x
		}
		for y < 0 {
			y = size + y
		}
		for x >= size {
			x = x - Size
		}
		for y >= size {
			y = y - Size
		}
		trace = append(trace, u[size*y+x])
	}

	var buffer bytes.Buffer
	Mark1Compress1(trace, &buffer)
	return buffer.Len()
}

// U is a universe simulator
func U(filename string, size int, next func(v []byte, rng *rand.Rand) []byte) {
	rng := rand.New(rand.NewSource(1))
	u := make([]byte, size*size)
	for i := 0; i < 9; i++ {
		u[rng.Intn(len(u))] = 255
	}

	images := &gif.GIF{}
	var palette = []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{0xff, 0xff, 0xff, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}

	for step := 0; step < Iterations; step++ {
		u = next(u, rng)
		verse := image.NewPaletted(image.Rect(0, 0, size*Scale, size*Scale), palette)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				b := u[i*size+j]
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
		for x := 0; x < int(float64(step)*float64(size)*Scale/float64(Iterations)); x++ {
			for y := size*Scale - 10; y < size*Scale; y++ {
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
	circle := NewCircle(Size)
	U("capitalistic.gif", Size, func(u []byte, rng *rand.Rand) []byte {
		for {
			ax, ay, bx, by := rng.Intn(Size), rng.Intn(Size), rng.Intn(Size), rng.Intn(Size)
			v := make([]byte, len(u))
			copy(v, u)
			v[ay*Size+ax], v[by*Size+bx] = v[by*Size+bx], v[ay*Size+ax]
			aBefore := circle.K(Size, u, ax, ay)
			bBefore := circle.K(Size, u, bx, by)
			aAfter := circle.K(Size, v, ax, ay)
			bAfter := circle.K(Size, v, bx, by)
			if aAfter < aBefore && bAfter < bBefore {
				return v
			}
		}
	})
}

// Communistic model
func Communistic() {
	circle := NewCircle(Size)
	U("communistic.gif", Size, func(u []byte, rng *rand.Rand) []byte {
		for j := 0; j < 256; j++ {
			ax, ay, bx, by := rng.Intn(Size), rng.Intn(Size), rng.Intn(Size), rng.Intn(Size)
			v := make([]byte, len(u))
			copy(v, u)
			v[ay*Size+ax], v[by*Size+bx] = v[by*Size+bx], v[ay*Size+ax]
			before := 0
			for y := 0; y < Size; y++ {
				for x := 0; x < Size; x++ {
					before += circle.K(Size, u, x, y)
				}
			}
			after := 0
			for y := 0; y < Size; y++ {
				for x := 0; x < Size; x++ {
					after += circle.K(Size, v, x, y)
				}
			}
			if after < before {
				return v
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

	rng := rand.New(rand.NewSource(1))
	iris := Load()
	fmt.Println(len(iris))

	const (
		Width  = 512
		Height = 512
	)
	start := image.NewRGBA(image.Rect(0, 0, Width, Height))
	stop := image.NewRGBA(image.Rect(0, 0, Width, Height))
	var palette = []color.Color{
		color.RGBA{0xff, 0, 0, 0xff},
		color.RGBA{0, 0xff, 0, 0xff},
		color.RGBA{0, 0, 0xff, 0xff},
	}

	done := make(chan bool, 8)
	process := func(i, j int) {
		project := NewMatrix(4, 2)
		for r := 0; r < project.Rows; r++ {
			for c := 0; c < project.Cols; c++ {
				project.Data = append(project.Data, rng.Float32())
			}
		}
		for r := 0; r < project.Rows; r++ {
			row := project.Data[r*project.Cols : (r+1)*project.Cols]
			norm := sqrt(vector.Dot(row, row))
			for k := range row {
				row[k] /= norm
			}
		}
		projections := make([]Matrix, 0, 8)
		for _, flower := range iris {
			point := NewMatrix(4, 1)
			for _, x := range flower.Measures {
				point.Data = append(point.Data, float32(x))
			}
			projection := project.MulT(point)
			projections = append(projections, projection)
		}
		max := [2]float32{}
		for _, projection := range projections {
			for k := range max {
				if projection.Data[k] > max[k] {
					max[k] = projection.Data[k]
				}
			}
		}
		u := make([]byte, 256*256)
		type Point struct {
			X     [2]int
			Color byte
		}
		white, black := make([]Point, 0, 8), make([]Point, 0, 8)
		for i, projection := range projections {
			color := byte(255)
			if iris[i].Label == "Iris-versicolor" {
				color = 254
			} else if iris[i].Label == "Iris-virginica" {
				color = 253
			}
			x, y := int(255*projection.Data[0]/max[0]), int(255*projection.Data[1]/max[1])
			u[y*256+x] = 255
			white = append(white, Point{
				X:     [2]int{x, y},
				Color: color,
			})
		}
		for y := 0; y < 256; y++ {
			for x := 0; x < 256; x++ {
				if u[y*256+x] == 0 {
					black = append(black, Point{
						X:     [2]int{x, y},
						Color: 0,
					})

				}
			}
		}
		for _, point := range white {
			start.Set(point.X[0]+i*256, point.X[1]+j*256, palette[point.Color-253])
		}
		for _, point := range black {
			start.Set(point.X[0]+i*256, point.X[1]+j*256, color.RGBA{0, 0, 0, 0xFF})
		}
		circle := NewCircle(256)
		for s := 0; s < 1024; s++ {
			for {
				a, b := rng.Intn(len(white)), rng.Intn(len(black))
				v := make([]byte, len(u))
				copy(v, u)
				ax, ay, bx, by := white[a].X[0], white[a].X[1], black[b].X[0], black[b].X[1]
				v[ay*256+ax], v[by*256+bx] = v[by*256+bx], v[ay*256+ax]
				aBefore := circle.K(256, u, ax, ay)
				bBefore := circle.K(256, u, bx, by)
				aAfter := circle.K(256, v, ax, ay)
				bAfter := circle.K(256, v, bx, by)
				if aAfter < aBefore && bAfter < bBefore {
					white[a].X, black[b].X = black[b].X, white[a].X
					u = v
					break
				}
			}
			fmt.Println(s)
		}
		for _, point := range white {
			stop.Set(point.X[0]+i*256, point.X[1]+j*256, palette[point.Color-253])
		}
		for _, point := range black {
			stop.Set(point.X[0]+i*256, point.X[1]+j*256, color.RGBA{0, 0, 0, 0xFF})
		}
		done <- true
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			go process(i, j)
		}
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			<-done
		}
	}

	{
		out, err := os.Create("iris_start.png")
		if err != nil {
			panic(err)
		}
		defer out.Close()

		err = png.Encode(out, start)
		if err != nil {
			panic(err)
		}
	}
	{
		out, err := os.Create("iris_stop.png")
		if err != nil {
			panic(err)
		}
		defer out.Close()

		err = png.Encode(out, stop)
		if err != nil {
			panic(err)
		}
	}
}
