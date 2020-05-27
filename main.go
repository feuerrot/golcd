package main

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/gonutz/framebuffer"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

func calcorigin(d *font.Drawer, s string, fb image.Image) fixed.Point26_6 {
	d.Dot = fixed.P(0, 0)
	bounds, _ := d.BoundString(s)
	fbsize := fb.Bounds()
	fbSize := fixed.P(fbsize.Max.X, fbsize.Max.Y)

	sSize := fixed.Point26_6{
		X: bounds.Max.X,
		Y: bounds.Min.Y,
	}
	tmp := fbSize.Sub(sSize)
	origin := tmp.Div(fixed.I(2))

	return origin
}

func getColor(offset int) color.Color {
	index := time.Now().Hour()*15 + offset
	return colorful.Hsv(float64(index), 1, 1)
}

func main() {
	fb, err := framebuffer.Open("/dev/fb1")
	if err != nil {
		panic(err)
	}
	defer fb.Close()

	buffer := image.NewRGBA(fb.Bounds())
	bg := image.NewUniform(getColor(0))
	fg := image.Black
	draw.Draw(buffer, buffer.Bounds(), bg, image.ZP, draw.Src)

	ttfont, err := truetype.Parse(gomono.TTF)
	ttfoption := &truetype.Options{
		Size:    64,
		Hinting: font.HintingFull,
	}

	d := &font.Drawer{
		Dst:  buffer,
		Src:  fg,
		Face: truetype.NewFace(ttfont, ttfoption),
	}

	var origin = calcorigin(d, "23:59:59", fb)

	var t time.Time
	var s string

	for {
		t = time.Now()
		s = t.Format("15:04:05")

		bg.C = getColor(0)
		draw.Draw(buffer, buffer.Bounds(), bg, image.ZP, draw.Src)
		d.Dot = origin
		d.DrawString(s)
		sleep := time.Until(time.Now().Truncate(time.Second).Add(time.Second))
		time.Sleep(sleep)
		draw.Draw(fb, fb.Bounds(), buffer, image.ZP, draw.Src)
	}
}
