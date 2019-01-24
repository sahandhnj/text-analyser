package main

import (
	rake "github.com/sahandhnj/text-analyser/src/rakeimpl"
	"github.com/sahandhnj/text-analyser/src/types"
)

func rakeIt(text string, lang rake.LANG) []types.Token {
	if lang == "" {
		lang = rake.LANG_EN
	}

	candidates := rake.RunRake(text, lang)
	tokens := make([]types.Token, 0)

	var limit float64
	top20 := float64(len(candidates)) * 0.2

	for i, candidate := range candidates {
		if float64(i) > top20 {
			limit = candidate.Value
			break
		}
	}

	for _, candidate := range candidates {
		//fmt.Printf("%s --> %f\n", candidate.Key, candidate.Value)
		if candidate.Value > limit {
			contexts := getContext(candidate.Key)
			for _, context := range contexts {
				tokens = append(tokens, types.Token{
					Value: context,
					Score: candidate.Value,
				})
			}
		}
	}

	return tokens
}
