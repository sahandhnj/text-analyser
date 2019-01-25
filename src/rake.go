package main

import (
	rake "github.com/sahandhnj/text-analyser/src/rakeimpl"
	"github.com/sahandhnj/text-analyser/src/types"
)

func rakeIt(text string, lang rake.LANG) []types.Token {
	if lang == "" {
		lang = rake.LANG_EN
	}

	// buf, err := ioutil.ReadFile("text") // just pass the file name
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// s := string(buf)

	candidates := rake.RunRake(text, lang)
	tokens := make([]types.Token, 0)

	var limit float64
	var limitFactor float64
	switch {
	case len(candidates) < 5:
		limitFactor = 1
	case len(candidates) < 10:
		limitFactor = 0.8
	case len(candidates) < 20:
		limitFactor = 0.5
	default:
		limitFactor = 0.2
	}

	top20 := float64(len(candidates)) * limitFactor
	for i, candidate := range candidates {
		if float64(i) > top20 {
			limit = candidate.Value
			break
		}
	}

	for _, candidate := range candidates {
		if candidate.Value > limit {
			// fmt.Printf("%s --> %f\n", candidate.Key, candidate.Value)
			if lang == rake.LANG_EN {
				contexts := getContext(candidate.Key)
				for _, context := range contexts {
					tokens = append(tokens, types.Token{
						Value: context,
						Score: candidate.Value,
					})
				}
			}

			if lang == rake.LANG_NL {
				tokens = append(tokens, types.Token{
					Value: candidate.Key,
					Score: candidate.Value,
				})
			}

		}
	}

	return tokens
}
