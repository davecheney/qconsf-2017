// mandelbrot example code adapted from Francesc Campoy's mandelbrot package.
// https://github.com/campoy/mandelbrot
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/bits"
	"os"
	"sync"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	const (
		height = 1024
		width  = 1024
		output = "mandelbrot.png"
	)

	// open a new file
	f, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}

	// create the image
	img := make([][]color.RGBA, height)
	for i := range img {
		img[i] = make([]color.RGBA, width)
	}

	// compute mandelbrot
	var wg sync.WaitGroup
	wg.Add(height)
	for i := range img {
		go func(i int) {
			defer wg.Done()
			for j := range img[i] {
				fillPixel(&img[i][j], i, j)
			}
		}(i)
	}
	wg.Wait()

	// and encoding it
	if err := png.Encode(f, Image(img)); err != nil {
		log.Fatal(err)
	}
}

type Image [][]color.RGBA

func (i Image) At(x, y int) color.Color { return i[x][y] }
func (i Image) ColorModel() color.Model { return color.RGBAModel }
func (i Image) Bounds() image.Rectangle { return image.Rect(0, 0, len(i), len(i[0])) }

func fillPixel(p *color.RGBA, x, y int) {
	const n = 1000
	const Limit = 2.0
	Zr, Zi, Tr, Ti := 0.0, 0.0, 0.0, 0.0
	Cr := (2*float64(x)/float64(n) - 1.5)
	Ci := (2*float64(y)/float64(n) - 1.0)

	for i := 0; i < n && (Tr+Ti <= Limit*Limit); i++ {
		Zi = 2*Zr*Zi + Ci
		Zr = Tr - Ti + Cr
		Tr = Zr * Zr
		Ti = Zi * Zi
	}
	paint(p, Tr, Ti)
}

func paint(p *color.RGBA, x, y float64) {
	n := byte(x * y * 2)
	p.R, p.G, p.B, p.A = bits.RotateLeft8(n, 2), bits.RotateLeft8(n, 2), bits.RotateLeft8(n, 6), 255
}
