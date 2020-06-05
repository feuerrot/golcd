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

var origin fixed.Point26_6

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
	index := (time.Now().Add(time.Second).Minute()*6 + offset) % 360
	return colorful.Hsv(float64(index), 1, 1)
}

func drawTime(tick time.Time, d *font.Drawer, bg image.Image, buffer draw.Image) {
	// We need to draw the next second, as we do a sleep until the start of the next second
	s := tick.Format("15:04:05")

	draw.Draw(buffer, buffer.Bounds(), bg, image.ZP, draw.Src)
	d.Dot = origin
	d.DrawString(s)
}

func sleepUntilNextSecond(tick time.Time) {
	//fmt.Printf("now: %v, trunc: %v, next: %v\n", time.Now(), time.Now().Truncate(time.Second), tick)
	time.Sleep(time.Until(tick))
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

	origin = calcorigin(d, "23:59:59", fb)

	for {
		nextTick := time.Now().Truncate(time.Second).Add(time.Second)
		bg.C = getColor(time.Now().Add(time.Second).Truncate(time.Second).Second() / 10)
		drawTime(nextTick, d, bg, buffer)

		sleepUntilNextSecond(nextTick)

		// copy buffer
		draw.Draw(fb, fb.Bounds(), buffer, image.ZP, draw.Src)
	}
}
