// Copyright 2014 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"strings"
	"syscall/js"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	count        int
	gophersImage *ebiten.Image
	op           = &ebiten.DrawImageOptions{}
	col          color.Color
	cam          *ebiten.Image
)

func update(screen *ebiten.Image) error {
	count++

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Center the image on the screen.
	//w, h := gophersImage.Size()
	//op.GeoM.Translate(float64(screenWidth-w)/2, float64(screenHeight-h)/2)

	// Rotate the hue.
	op.ColorM.RotateHue(float64(count%360) * 2 * math.Pi / 360)

	if cam != nil {
		screen.DrawImage(cam, op)
		return nil
	}

	if col != nil {
		screen.Fill(col)
	} else {
		screen.DrawImage(gophersImage, op)
	}
	return nil
}

func main() {
	// Decode image from a byte slice instead of a file so that
	// this example works in any working directory.
	// If you want to use a file, there are some options:
	// 1) Use os.Open and pass the file to the image decoder.
	//    This is a very regular way, but doesn't work on browsers.
	// 2) Use ebitenutil.OpenFile and pass the file to the image decoder.
	//    This works even on browsers.
	// 3) Use ebitenutil.NewImageFromFile to create an ebiten.Image directly from a file.
	//    This also works on browsers.
	img, _, err := image.Decode(bytes.NewReader(images.Gophers_jpg))
	if err != nil {
		log.Fatal(err)
	}
	gophersImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	sendFrame := func(i []js.Value) {
		if len(i) <= 0 {
			return
		}

		// data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAUAAAADwCAYAAABxLb1rAAAG+UlEQVR4Xu3UAQ0AMAwCQeZf9JbMxl
		prefix := "data:image/png;base64,"
		data := i[0].String()
		if !strings.HasPrefix(data, prefix) {
			println("wrong prefix")
			return
		}

		b, err := base64.StdEncoding.DecodeString(data[len(prefix):])
		if err != nil {
			println("error b64 decode")
			return
		}

		r := bytes.NewReader(b)
		im, _, err := image.Decode(r)
		if err != nil || im == nil {
			println("error decode")
			return
		}

		cam, err = ebiten.NewImageFromImage(im, ebiten.FilterDefault)

		//if rand.Intn(2) == 0 {
		//	col = color.White
		//} else {
		//	col = color.Black
		//}

		//gophersImage.Fill(color.White)
		//js.Global().Set("output", js.ValueOf(i[0].Int()+i[1].Int()))
		//println(js.ValueOf(i[0].Int() + i[1].Int()).String())
	}

	js.Global().Set("sendFrame", js.NewCallback(sendFrame))

	ebiten.SetMaxTPS(30)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Hue (Ebiten Demo)"); err != nil {
		log.Fatal(err)
	}
}
