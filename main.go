package main

import (
	"image"
	"image/draw"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/gonutz/framebuffer"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

func centerdraw(d *font.Drawer, s string, fb image.Image) {
	d.Dot = fixed.P(0, 0)
	bounds, _ := d.BoundString(s)
	fbsize := fb.Bounds()
	fbSize := fixed.P(fbsize.Max.X, fbsize.Max.Y)

	sSize := fixed.Point26_6{bounds.Max.X, bounds.Min.Y}
	tmp := fbSize.Sub(sSize)
	d.Dot = tmp.Div(fixed.I(2))

	d.DrawString(s)
}

func calcorigin(d *font.Drawer, s string, fb image.Image) fixed.Point26_6 {
	d.Dot = fixed.P(0, 0)
	bounds, _ := d.BoundString(s)
	fbsize := fb.Bounds()
	fbSize := fixed.P(fbsize.Max.X, fbsize.Max.Y)

	sSize := fixed.Point26_6{bounds.Max.X, bounds.Min.Y}
	tmp := fbSize.Sub(sSize)
	origin := tmp.Div(fixed.I(2))

	return origin
}

func main() {
	fb, err := framebuffer.Open("/dev/fb1")
	if err != nil {
		panic(err)
	}
	defer fb.Close()

	img := image.NewRGBA(fb.Bounds())
	draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)

	ttfont, err := truetype.Parse(gomono.TTF)
	ttfoption := &truetype.Options{
		Size:    64,
		Hinting: font.HintingFull,
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: truetype.NewFace(ttfont, ttfoption),
	}

	var origin = calcorigin(d, "23:59:59", img)

	var t time.Time
	var s string

	for {
		t = time.Now()
		s = t.Format("15:04:05")

		draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
		d.Dot = origin
		d.DrawString(s)
		sleep := time.Until(time.Now().Truncate(time.Second).Add(time.Second))
		time.Sleep(sleep)
		draw.Draw(fb, fb.Bounds(), img, image.ZP, draw.Src)
	}
}
