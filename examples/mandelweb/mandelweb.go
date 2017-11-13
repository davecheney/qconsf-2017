// mandelbrot example code adapted from Francesc Campoy's mandelbrot package.
// https://github.com/campoy/mandelbrot
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/bits"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	http.HandleFunc("/mandelbrot", mandelbrot)
	log.Println("listening on http://127.0.0.1:8080/")
	http.ListenAndServe(":8080", logRequest(http.DefaultServeMux))
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, req)
		log.Println(req.RemoteAddr, req.RequestURI, time.Since(start))
	})
}

func mandelbrot(w http.ResponseWriter, req *http.Request) {
	const height, width = 512, 512
	img := make([][]color.RGBA, height)
	for i := range img {
		img[i] = make([]color.RGBA, width)
	}
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
	png.Encode(w, Image(img))
}

type Image [][]color.RGBA

func (i Image) At(x, y int) color.Color { return i[x][y] }
func (i Image) ColorModel() color.Model { return color.RGBAModel }
func (i Image) Bounds() image.Rectangle { return image.Rect(0, 0, len(i), len(i[0])) }

func fillPixel(p *color.RGBA, x, y int) {
	const n = 1000
	const Limit = 2.0
	const Zoom = 4
	Zr, Zi, Tr, Ti := 0.0, 0.0, 0.0, 0.0
	Cr := (Zoom*float64(x)/float64(n) - 1.5)
	Ci := (Zoom*float64(y)/float64(n) - 1.0)

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
