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

// T computes the transpose
func T(u [8 * 8]byte) [8 * 8]byte {
	n := [8 * 8]byte{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			n[j*8+i] = u[i*8+j]
		}
	}
	return n
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
	iterations := 4096
	for step := 0; step < iterations; step++ {
		var buffer bytes.Buffer
		compress.Mark1Compress1(u[:], &buffer)
		best := buffer.Len()
		uu := T(u)
		buffer = bytes.Buffer{}
		compress.Mark1Compress1(uu[:], &buffer)
		best += buffer.Len()
		for {
			v := u
			a, b := rng.Intn(len(v)), rng.Intn(len(v))
			v[a], v[b] = v[b], v[a]
			vv := T(v)
			var buffer bytes.Buffer
			compress.Mark1Compress1(v[:], &buffer)
			var buffer2 bytes.Buffer
			compress.Mark1Compress1(vv[:], &buffer2)
			if buffer.Len()+buffer2.Len() <= best {
				u = v
				break
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
		images.Delay = append(images.Delay, 2)
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
