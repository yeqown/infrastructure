package main

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

// getSysTextList get font list from curretn system
// TODO:
func getSysTextList() []string {
	return nil
}
