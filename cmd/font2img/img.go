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

// ImgFormat .
type ImgFormat string

const (
	unknownFormat ImgFormat = "unknown"
	jpegFormat              = "jpeg"
	pngFormat               = "png"
)

// NewImg .
func NewImg(bg *background, t *text, opt outputOption) *Img {
	img := Img{
		bg:        bg,
		txt:       t,
		outputOpt: opt,
	}

	img.init()
	return &img
}

type outputOption struct {
	filename string // save filename
	format   ImgFormat
}

// Img .
type Img struct {
	bg        *background  // background setting
	txt       *text        // font and font setting
	outputOpt outputOption // output options
	rgba      *image.RGBA  // image
}

// create a image buffer to draw image in memory
func (img *Img) init() {
	// init rgba image
	img.rgba = image.NewRGBA(image.Rect(0, 0, img.bg.W, img.bg.H))

	// calculate the parameters of text position and options
	img.txt.autoCalculate(img.bg.W, img.bg.H)
}

// img process
func (img *Img) process() {
	img.bg.draw(img.rgba)
	img.txt.draw(img.rgba)
}

func (img *Img) outputFilename() string {
	return img.outputOpt.filename + "." + string(img.outputOpt.format)
}

// Save . output the image to disk
func (img *Img) Save() (err error) {
	img.process()

	// open file
	fd, err := os.OpenFile(img.outputFilename(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("could not open filename=%s, err=%v", img.outputFilename(), err)
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
