package main

import (
	"context"
	"time"

	corenlp "github.com/hironobu-s/go-corenlp"
	"github.com/hironobu-s/go-corenlp/connector"
	"github.com/sahandhnj/text-analyser/types"
)

func getContext(text string) []string {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	c := connector.NewHTTPClient(ctx, "http://127.0.0.1:9000/")
	c.Annotators = []string{"pos"}

	doc, err := corenlp.Annotate(c, text)
	if err != nil {
		panic(err)
	}

	contexts := make([]string, 0)

	for _, sentence := range doc.Sentences {
		var word string = ""

		for _, token := range sentence.Tokens {
			if word != "" && !capturePos(token.Pos) {
				contexts = append(contexts, word)
				word = ""
			}

			if capturePos(token.Pos) {
				if word != "" {
					word = word + " " + token.Word
				} else {
					word = token.Word
				}
			}
		}

		if word != "" {
			contexts = append(contexts, word)
		}
	}

	return contexts
}

func nlp(text string) []types.NerToken {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	c := connector.NewHTTPClient(ctx, "http://127.0.0.1:9000/")
	c.Annotators = []string{"ner"}

	doc, err := corenlp.Annotate(c, text)
	if err != nil {
		panic(err)
	}

	ners := make([]types.NerToken, 0)

	for _, sentence := range doc.Sentences {
		var word string = ""
		var nertype string = ""
		for _, token := range sentence.Tokens {
			if word != "" && nertype != token.Ner {
				ners = append(ners, types.NerToken{
					Word: word,
					Type: nertype,
				})

				word = ""
				nertype = ""
			}

			if !skipNers(token.Ner) {
				nertype = token.Ner
				if word != "" {
					word = word + " " + token.Word
				} else {
					word = token.Word
				}
			}
		}
	}

	return ners
}

func skipNers(category string) bool {
	switch category {
	case
		"O",
		"DATE",
		"MONEY",
		"ORDINAL",
		"TIME",
		"NUMBER":
		return true
	}
	return false
}

func capturePos(category string) bool {
	switch category {
	case
		"NN",
		"NNS",
		"JJ":
		return true
	}
	return false
}
