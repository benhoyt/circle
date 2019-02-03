// Draw circles using the Bresenham Circle Algorithm

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

const defaultRadius = 3

func main() {
	radius := flag.Int("r", defaultRadius, "radius of circle")
	httpAddr := flag.String("http", "", "run HTTP server on `address` (e.g., :8080) to serve circle images\nexample URL: http://localhost:8080?r=42")
	flag.Parse()

	if *httpAddr != "" {
		http.HandleFunc("/", httpHandler)
		log.Printf("listening on %s", *httpAddr)
		log.Fatal(http.ListenAndServe(*httpAddr, nil))
	} else {
		writer := bufio.NewWriter(os.Stdout)
		_ = drawCircleText(*radius, writer)
		_ = writer.Flush()
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	radius := defaultRadius
	radiusStr := r.FormValue("r")
	if radiusStr != "" {
		var err error
		radius, err = strconv.Atoi(radiusStr)
		if err != nil || radius < 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `radius must be a positive integer, e.g., ?r=42`)
			return
		}
	}
	log.Printf("drawing circle of radius %d", radius)
	err := drawCircleImage(radius, w)
	if err != nil {
		fmt.Fprintf(w, "error encoding image: %s", err)
		return
	}
}

// Draw circle of radius r and write to writer as PNG.
func drawCircleImage(r int, writer io.Writer) error {
	// Encode image as PNG with white pixels for circle
	size := r*2 + 1
	im := image.NewNRGBA(image.Rect(0, 0, size, size))
	drawCircleInt(r, func(x, y int) {
		im.SetNRGBA(x+r, r-y, color.NRGBA{255, 255, 255, 255})
	})
	return png.Encode(writer, im)
}

// Draw circle of radius r and write to writer as text.
func drawCircleText(r int, writer io.Writer) error {
	// Create 2-D "screen buffer" to hold pixels, but include room
	// for a newline after each line. For example, with r=3 showing,
	// the top-right quadrant drawn (newlines are denoted by 'N'):
	// ...##..N
	// ...|.#.N
	// ...|..#N
	// ---+--#N
	// ...|...N
	// ...|...N
	// ...|...N
	size := r*2 + 1
	buf := bytes.Repeat([]byte{' '}, (size+1)*size)
	for i := size; i < len(buf); i += size + 1 {
		buf[i] = '\n'
	}
	drawCircleInt(r, func(x, y int) {
		buf[(size+1)*(r-y)+x+r] = '#'
	})
	_, err := writer.Write(buf)
	return err
}

// Draw circle of radius r using given putPixel function; use simple
// square root method.
func drawCircleSqrt(r int, putPixel func(x, y int)) {
	y := r
	rsq := r * r
	for x := 0; x <= y; x++ {
		// Just calculate y = sqrt(r^2 - x^2)
		y := int(math.Round(math.Sqrt(float64(rsq - x*x))))
		putPixel(x, y)
		putPixel(y, x)
		putPixel(-x, y)
		putPixel(-y, x)
		putPixel(x, -y)
		putPixel(y, -x)
		putPixel(-x, -y)
		putPixel(-y, -x)
	}
}

// Draw circle of radius r using given putPixel function; use
// Bresenham-ish method with no sqrt and only integer math.
func drawCircleInt(r int, putPixel func(x, y int)) {
	x := 0
	y := r
	xsq := 0
	rsq := r * r
	ysq := rsq
	// Loop x from 0 to the line x==y. Start y at r and each time
	// around the loop either keep it the same or decrement it.
	for x <= y {
		putPixel(x, y)
		putPixel(y, x)
		putPixel(-x, y)
		putPixel(-y, x)
		putPixel(x, -y)
		putPixel(y, -x)
		putPixel(-x, -y)
		putPixel(-y, -x)

		// New x^2 = (x+1)^2 = x^2 + 2x + 1
		xsq = xsq + 2*x + 1
		x++
		// Potential new y^2 = (y-1)^2 = y^2 - 2y + 1
		y1sq := ysq - 2*y + 1
		// Choose y or y-1, whichever gives smallest error
		a := xsq + ysq
		b := xsq + y1sq
		if a-rsq >= rsq-b {
			y--
			ysq = y1sq
		}
	}
}
