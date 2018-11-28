package main

import (
	rake "github.com/Obaied/RAKE.Go"
	"github.com/sahandhnj/text-analyser/types"
)

func rakeIt(text string) []types.Token {
	candidates := rake.RunRake(text)
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
