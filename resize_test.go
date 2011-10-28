package hu

import (
	"testing"
	"os"
	"bytes"
	"io"
	"image"
	"image/jpeg"
	"fmt"
)

func TestResizeParse(t *testing.T) {
	f, err := os.Open("eikeon.com/static/test.jpg")
	check(err)

	// Grab the image data
	var buf bytes.Buffer
	io.Copy(&buf, f)
	i, _, err := image.Decode(&buf)
	check(err)

	// // Resize if too large, for more efficient moustachioing.
	// // We aim for less than 1200 pixels in any dimension; if the
	// // picture is larger than that, we squeeze it down to 600.
	// const max = 640*2
	// if b := i.Bounds(); b.Dx() > max || b.Dy() > max {
	// 	// If it's gigantic, it's more efficient to downsample first
	// 	// and then resize; resizing will smooth out the roughness.
	// 	if b.Dx() > 2*max || b.Dy() > 2*max {
	// 		w, h := max, max
	// 		if b.Dx() > b.Dy() {
	// 			h = b.Dy() * h / b.Dx()
	// 		} else {
	// 			w = b.Dx() * w / b.Dy()
	// 		}
	// 		i = Resample(i, i.Bounds(), w, h)
	// 		b = i.Bounds()
	// 	}
	// 	w, h := max/2, max/2
	// 	if b.Dx() > b.Dy() {
	// 		h = b.Dy() * h / b.Dx()
	// 	} else {
	// 		w = b.Dx() * w / b.Dy()
	// 	}
	// 	i = Resize(i, i.Bounds(), w, h)
	// }

	b := i.Bounds()
	//w, h := 1280, 1280
	w, h := 960, 960
	//w, h := 320, 320
	//w, h := 80, 80

	if b.Dx() > b.Dy() {
		h = b.Dy() * h / b.Dx()
	} else {
		w = b.Dx() * w / b.Dy()
	}
	i = Resize(i, i.Bounds(), w, h)

	// Encode as a new JPEG image.
	buf.Reset()
	err = jpeg.Encode(&buf, i, nil)
	check(err)

	out, err := os.Create("out.jpg")
	check(err)
	jpeg.Encode(out, i, nil)
	out.Close()
	fmt.Println("boo")
}

// check aborts the current execution if err is non-nil.
func check(err os.Error) {
	if err != nil {
		panic(err)
	}
}
