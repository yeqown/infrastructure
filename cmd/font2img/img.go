package main

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/yeqown/log"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type format string

const (
	unknownFormat format = "unknown"
	jpegFormat           = "jpeg"
	pngFormat            = "png"
)

// // Drawer .contains Draw method
// type Drawer interface {
// 	Draw() error
// }

// NewBackground .
func NewBackground(col string) *Background {
	var rgb color.RGBA
	switch col {
	case "white":
		rgb = color.RGBA{255, 255, 255, 0}
	case "black":
		rgb = color.RGBA{0, 0, 0, 0}
	}

	return &Background{
		color: rgb,
	}
}

// Background . image or pure color
type Background struct {
	color color.RGBA
}

// NewImg .
func NewImg(bg *Background, t *Text) *Img {
	img := Img{
		Background: bg,
		Text:       t,
		W:          600,
		H:          400,
		outputOpt: outputOption{
			filename: "output.jpg",
			format:   jpegFormat,
		},
	}

	img.init()
	return &img
}

type outputOption struct {
	filename string // save filename
	format   format
}

// Img .
type Img struct {
	Background *Background  // background setting
	Text       *Text        // font and font setting
	W          int          // width
	H          int          // height
	outputOpt  outputOption // output options
	rgba       *image.RGBA  // image
}

// create a image buffer to draw image in memory
func (img *Img) init() {
	// W, H should be paramters to pass in or calculated by background
	if img.W == 0 {
		img.W = 600
	}

	if img.H == 0 {
		img.H = 400
	}

	img.rgba = image.NewRGBA(image.Rect(0, 0, img.W, img.H))

	// calculate the parameters
	if img.Text.AutoCalculate {
		img.calculateTextOpt()
	}
}

// TODO: font family support
// TODO: auto phgraph ?
func (img *Img) calculateTextOpt() {
	img.Text.FontFamily = "asdajslk"
	img.Text.Size = img.H / 4
	img.Text.X = (img.H - img.Text.Size) / 2
	// img.Text.Y = (img.W - (img.Text.Size)*len(img.Text.Content)) / 2
	img.Text.Y = 100
	log.Infof("text = %v", img.Text)
}

// TODO: set background options
func (img *Img) drawBackground() {
	col := image.NewUniform(img.Background.color)
	draw.Draw(img.rgba, img.rgba.Bounds(), col, image.ZP, draw.Src)
}

// TODO: set font options
func (img *Img) drawText() {
	col := color.RGBA{200, 100, 0, 255}

	d := font.Drawer{
		Dst:  img.rgba,
		Src:  image.NewUniform(col),
		Face: nil,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(img.Text.X * 64), Y: fixed.Int26_6(img.Text.Y * 64)},
	}

	d.DrawString(img.Text.Content)
}

// FIXME: may i need the method to schedule the draw methods ?
func (img *Img) process() {
	img.drawBackground()
	img.drawText()
}

// Save . output the image to disk
func (img *Img) Save() (err error) {
	img.process()

	// open file
	fd, err := os.OpenFile(img.outputOpt.filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("could not open filename=%s, err=%v", img.outputOpt.filename, err)
		return err
	}
	defer fd.Close()

	// create buf
	buf := bufio.NewWriter(fd)

	// encode image to different format
	switch img.outputOpt.format {
	case jpegFormat:
		jpeg.Encode(buf, img.rgba, nil)
	case pngFormat:
		png.Encode(buf, img.rgba)
	default:
		err := errors.New("unknown format")
		log.Errorf("unknown format type=%s, err=%v", img.outputOpt.format, err)
		return err
	}

	// flush memory into fd
	return buf.Flush()
}
