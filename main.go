package main

import (
	"bytes"
	"encoding/base64"
	"image"
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
	cam          *ebiten.Image
	opt          = &ebiten.DrawImageOptions{}
)

func update(screen *ebiten.Image) error {
	count++

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// Rotate the hue.
	opt.ColorM.RotateHue(float64(count%360) * 2 * math.Pi / 360)

	if cam != nil {
		screen.DrawImage(cam, opt)
		return nil
	}

	screen.DrawImage(gophersImage, opt)

	return nil
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(images.Gophers_jpg))
	if err != nil {
		log.Fatal(err)
	}
	gophersImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)

	sendFrame := func(i []js.Value) {
		if len(i) <= 0 {
			return
		}

		// data:image/png;base64,iVBORw0KGgo...
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
	}

	js.Global().Set("sendFrame", js.NewCallback(sendFrame))

	ebiten.SetMaxTPS(30)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Peepers"); err != nil {
		log.Fatal(err)
	}
}
