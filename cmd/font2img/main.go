package main

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/yeqown/log"
)

// the questions to ask
var qsPhase1 = []*survey.Question{
	{
		Name:      "content",
		Prompt:    &survey.Input{Message: "Input your content?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a background color:",
			Options: []string{"white", "black"},
			Default: "black",
		},
	},
	// {
	// 	Name:   "age",
	// 	Prompt: &survey.Input{Message: "How old are you?"},
	// },
}

type qsPhase1Answer struct {
	Content string `survey:"content"`
	Color   string `survey:"color"`
}

// var qsPhase2 = []*survey.Question{
// }

func main() {
	// // the answers will be written to this struct
	// answers := struct {
	// 	Name          string // survey will match the question and field names
	// 	FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
	// 	Age           int    // if the types don't match, survey will convert it
	// }{}

	var qsPhase1Ans = qsPhase1Answer{}

	// perform the questions
	err := survey.Ask(qsPhase1, &qsPhase1Ans)
	if err != nil {
		log.Error(err)
		return
	}

	img := NewImg(NewBackground(qsPhase1Ans.Color), NewDefaultText(qsPhase1Ans.Content))
	if err := img.Save(); err != nil {
		log.Error(err)
	}
}
