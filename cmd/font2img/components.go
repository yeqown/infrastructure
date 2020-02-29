package main

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/yeqown/infrastructure/pkg/fontutil"
	"github.com/yeqown/log"
)

var (
	defaultColors = map[string]color.RGBA{
		"black": color.RGBA{0, 0, 0, 0},
		"white": color.RGBA{255, 255, 255, 0},
		"gray":  color.RGBA{},
		"blue":  color.RGBA{},
		"pink":  color.RGBA{174, 56, 121, 1},
	}
)

func getDefaultColorList() []string {
	colors := make([]string, 0, len(defaultColors))
	for k := range defaultColors {
		colors = append(colors, k)
	}
	return colors
}

// text .
type text struct {
	X             int
	Y             int
	Size          int
	FontFamily    string     // font family
	Content       string     // content
	AutoCalculate bool       // auto calculate the text setting related to Img and background
	DPI           int        // pixels count per block
	color         color.RGBA // color
}

// newtext .
// FIXME: conents be changed upper case of first character
func newtext(x, y, size int, family, color, content string) *text {
	log.Info(x, y, size, "family=", family, "content=", content)

	col, ok := defaultColors[color]
	if !ok {
		col = defaultColors["black"]
	}

	return &text{
		X:             x,
		Y:             y,
		FontFamily:    fontutil.AssemFontPath(family),
		Size:          size,
		Content:       content,
		AutoCalculate: false,
		DPI:           size, // FIXME: how to set DPI
		color:         col,
	}
}

// TODO:
func (t *text) autoCalculate(bgH int) {
	if !t.AutoCalculate {
		return
	}
	t.Size = bgH / 4
	t.X = (bgH - t.Size) / 2
	t.Y = 100 // TODO: auto calc
}

func (t *text) draw(dst *image.RGBA) (err error) {
	log.Infof("text=%v", *t)

	var (
		fontByts []byte
	)
	fontByts, err = ioutil.ReadFile(t.FontFamily)
	if err != nil {
		log.Error(err)
		return
	}
	// parse font file
	font, err := freetype.ParseFont(fontByts)
	if err != nil {
		log.Error(err)
		return
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(float64(t.DPI))
	ctx.SetFont(font)
	ctx.SetFontSize(float64(t.Size))
	ctx.SetClip(dst.Bounds())
	ctx.SetDst(dst)                       // background image setting here
	ctx.SetSrc(image.NewUniform(t.color)) // font color setting
	ctx.DrawString(t.Content, freetype.Pt(t.X, t.Y+t.Size))

	return nil
}

// newbackground .
func newbackground(col string, w, h int) *background {
	rgb, ok := defaultColors[col]
	if !ok {
		// true: could not found color by name
		log.Infof("could not load color=%s, then set default background color=white", col)
		rgb = defaultColors["white"]
	}

	// W, H with defualt value
	if w == 0 {
		w = 1600
	}

	if h == 0 {
		h = 900
	}

	return &background{
		color: rgb,
		H:     h,
		W:     w,
	}
}

// background . image or pure color
type background struct {
	color color.RGBA
	W     int // width
	H     int // height
}

// TODO: finish this part
func (bg *background) draw(dst *image.RGBA) error {
	col := image.NewUniform(bg.color)
	draw.Draw(dst, dst.Bounds(), col, image.ZP, draw.Src)
	return nil
}
