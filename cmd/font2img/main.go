package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/yeqown/infrastructure/pkg/fontutil"
	"github.com/yeqown/log"
)

func main() {
	var (
		bgQA      = bgQsAnswer{}
		contentQA = contentQsAnswer{}
		outputQA  = outputQsAnswer{}

		err error
	)

	// perform the questions
	if err = survey.Ask(bgQs, &bgQA); err != nil {
		log.Error(err)
		return
	}

	if err = survey.Ask(contentQs, &contentQA); err != nil {
		log.Error(err)
		return
	}

	if err = survey.Ask(outputQs, &outputQA); err != nil {
		log.Error(err)
		return
	}

	img := NewImg(
		newbackground(bgQA.Color, bgQA.W, bgQA.H),
		newtext(0, 0, contentQA.Size, contentQA.FontFamily, contentQA.Color, contentQA.Content),
		outputOption{filename: outputQA.Filename, format: ImgFormat(outputQA.Format)},
	)

	// generating img and save
	if err := img.Save(); err != nil {
		log.Error(err)
	}
}

// the background setting to ask
var bgQs = []*survey.Question{
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a background color:",
			Options: getDefaultColorList(),
			Default: "white",
		},
	},
	{
		Name: "w",
		Prompt: &survey.Input{
			Message: "Input background width: (default: 1600)",
			Default: "1600",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "h",
		Prompt: &survey.Input{
			Message: "Input background height: (default: 300)",
			Default: "300",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

type bgQsAnswer struct {
	Color string `survey:"color"`
	W     int    `survey:"w"`
	H     int    `survey:"h"`
}

// the text content setting to ask
var contentQs = []*survey.Question{
	{
		Name:      "content",
		Prompt:    &survey.Input{Message: "Input your content:"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose the text color:",
			Options: getDefaultColorList(),
			Default: "black",
		},
	},
	{
		Name: "fontfamily",
		Prompt: &survey.Select{
			Message: "Choose a font:",
			Options: fontutil.GetSysFontList(),
			Default: fontutil.GetSysDefaultFont(),
		},
	},
	{
		Name: "size",
		Prompt: &survey.Input{
			Message: "Input font size (px):",
			Default: "32",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

type contentQsAnswer struct {
	Content    string `survey:"content"`
	Color      string `survey:"color"`
	FontFamily string `survey:"fontfamily"`
	Size       int    `survey:"size"`
}

// the text output setting to ask
var outputQs = []*survey.Question{
	{
		Name: "filename",
		Prompt: &survey.Input{
			Message: "Input output filename:",
			Default: "untiltled",
		},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "format",
		Prompt: &survey.Select{
			Message: "Choose a output format:",
			Options: []string{"jpeg", "png"},
			Default: "jpeg",
		},
	},
}

type outputQsAnswer struct {
	Filename string `survey:"filename"`
	Format   string `survey:"format"`
}
