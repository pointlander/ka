// Copyright 2025 The KA Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/pointlander/compress"
)

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
	n := func(u [8 * 8]byte) [8 * 8]byte {
		n := [8 * 8]byte{}
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				n[j*8+i] = u[i*8+j]
			}
		}
		return n
	}
	for {
		var buffer bytes.Buffer
		compress.Mark1Compress1(u[:], &buffer)
		best := buffer.Len()
		uu := n(u)
		buffer = bytes.Buffer{}
		compress.Mark1Compress1(uu[:], &buffer)
		best += buffer.Len()
		for {
			v := u
			a, b := rng.Intn(len(v)), rng.Intn(len(v))
			v[a], v[b] = v[b], v[a]
			vv := n(v)
			var buffer bytes.Buffer
			compress.Mark1Compress1(v[:], &buffer)
			var buffer2 bytes.Buffer
			compress.Mark1Compress1(vv[:], &buffer2)
			if buffer.Len()+buffer2.Len() <= best {
				u = v
				break
			}
		}
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				b := u[i*8+j]
				if b == 0 {
					fmt.Printf("0")
				} else {
					fmt.Printf("1")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
