package main

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/yeqown/infrastructure/pkg/fontutil"
	"github.com/yeqown/log"
	pkgfont "golang.org/x/image/font"
)

// Text .
type Text struct {
	X             int
	Y             int
	Size          int
	FontFamily    string // font family
	Content       string // content
	AutoCalculate bool   // auto calculate the text setting related to Img and Background
}

// NewDefaultText .
// TODO:
func NewDefaultText(content string) *Text {
	return &Text{
		X:             0,  // auto calculated
		Y:             0,  // auto calculated
		FontFamily:    "", // first font famliy in curretn system, FIXME: set default font family for system
		Size:          0,  // auto calculated
		Content:       content,
		AutoCalculate: true,
	}
}

// NewText .
func NewText(x, y, size int, family, content string) *Text {
	return &Text{
		X:             x,
		Y:             y,
		FontFamily:    family,
		Size:          size,
		Content:       content,
		AutoCalculate: false,
	}
}

// TODO: font family support
// TODO: auto phgraph ?
func (t *Text) calculateTextOpt(H int) {
	t.FontFamily = fontutil.AssemFontPath("JetBrainsMono-Regular.ttf")
	// t.FontFamily = "C:\\Users\\yeqown\\AppData\\Local\\Microsoft\\Windows\\Fonts\\JetBrainsMono-Regular.ttf"
	t.Size = H / 4
	t.X = (H - t.Size) / 2
	// t.Y = (img.W - (t.Size)*len(t.Content)) / 2
	t.Y = 100
	log.Infof("text = %v", t)
}

func (t *Text) draw(dst *image.RGBA) (err error) {
	var (
		dpi      float64 = 200
		size     float64 = 20
		fontByts []byte
	)
	// parse font file
	fontByts, err = ioutil.ReadFile(t.FontFamily)
	if err != nil {
		log.Error(err)
		return
	}
	font, err := freetype.ParseFont(fontByts)
	if err != nil {
		log.Error(err)
		return
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(font)
	ctx.SetFontSize(size)
	ctx.SetClip(dst.Bounds())
	ctx.SetDst(dst)
	ctx.SetSrc(image.Black)
	ctx.SetHinting(pkgfont.HintingNone)
	// or
	// ctx.SetHinting(pkgfont.HintingFull)

	// Draw the text.
	// TODO: set init position for text
	pt := freetype.Pt(200, 10+int(ctx.PointToFixed(size)>>6))
	ctx.DrawString(t.Content, pt)
	return nil
}

var (
	defaultColors = map[string]color.RGBA{
		"white": color.RGBA{0, 0, 0, 0},
		"black": color.RGBA{255, 255, 255, 0},
		"gray":  color.RGBA{},
		"blue":  color.RGBA{},
		"pink":  color.RGBA{174, 56, 121, 1},
	}
)

// NewBackground .
func NewBackground(col string) *Background {
	rgb, ok := defaultColors[col]
	if !ok {
		// true: could not found color by name
		rgb = defaultColors["white"]
	}

	return &Background{
		color: rgb,
	}
}

// Background . image or pure color
type Background struct {
	color color.RGBA
}

// TODO: finish this part
func (bg *Background) draw(dst *image.RGBA) error {
	col := image.NewUniform(bg.color)
	draw.Draw(dst, dst.Bounds(), col, image.ZP, draw.Src)
	return nil
}
