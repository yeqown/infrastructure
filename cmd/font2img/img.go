package main

import (
	"bufio"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/yeqown/log"
)

type format string

const (
	unknownFormat format = "unknown"
	jpegFormat           = "jpeg"
	pngFormat            = "png"
)

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

	// init rgba image
	img.rgba = image.NewRGBA(image.Rect(0, 0, img.W, img.H))

	// calculate the parameters
	if img.Text.AutoCalculate {
		// true: open the switch of text position
		// TODO: pass in background params and text options
		img.Text.calculateTextOpt(img.H)
	}
}

// FIXME: may i need the method to schedule the draw methods ?
func (img *Img) process() {
	img.Background.draw(img.rgba)
	img.Text.draw(img.rgba)
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
