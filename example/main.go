package main

import (
	"github.com/magsoft-se/wg"
	"image/color"
)

const WIDTH = 800
const HEIGHT = 600
const PORT = 80

const sqrwidth = 20
const sqrheight = 20
const step = 3
const NUM_POINTS = 100

type Point struct {
	x int
	y int
}

var points [NUM_POINTS]Point
var gx, gy int

func GameLoop() {
	wg.ClearImage(color.Black)

	if wg.GetKey(37) {
		if gx > 0 {
			gx -= step
		}
	}
	if wg.GetKey(39) {
		if gx < WIDTH-1-sqrwidth {
			gx += step
		}
	}
	if wg.GetKey(38) {
		if gy > 0 {
			gy -= step
		}
	}
	if wg.GetKey(40) {
		if gy < HEIGHT-1-sqrheight {
			gy += step
		}
	}

	if wg.GetMLBtn() {
		for i := NUM_POINTS - 1; i > 0; i-- {
			points[i] = points[i-1]
		}
		points[0].x = wg.GetMX()
		points[0].y = wg.GetMY()
	}

	for i := 0; i < NUM_POINTS-1; i++ {
		wg.GetImage().Set(points[i].x, points[i].y, color.RGBA{255, 255, 127, 255})
	}

	for i := 0; i < sqrwidth; i++ {
		for j := 0; j < sqrheight; j++ {
			wg.GetImage().Set(gx+i, gy+j, color.RGBA{127, 127, 255, 255})
		}
	}
}

func main() {
	// Place index.html (from eg GOPATH/src/github.com/magsoft-se/wg) in the same folder as this file (main.go)
	// Make sure width, height and port passed to wg.Start are the same as in index.html
	// The function ref passed, eg GameLoop will be called periodically every frame in a goroutine
	// First start this program. Then point a browser to localhost and you should se a blue square controllable
	// with the arrow keys, and be able to draw a dotted line with the mouse.
	wg.Start(WIDTH, HEIGHT, PORT, GameLoop)
}
