package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	rake "github.com/Obaied/RAKE.Go"
	corenlp "github.com/hironobu-s/go-corenlp"
	"github.com/hironobu-s/go-corenlp/connector"
	"github.com/sahandhnj/text-analyser/types"

	"gopkg.in/jdkato/prose.v2"
)

const Address = ":3005"

func main() {
	http.HandleFunc("/analyse", analyseRequest)

	fmt.Println("Listening on " + Address)
	err := http.ListenAndServe(Address, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Req struct {
	Text string `json:"text"`
}

type Resp struct {
	Tokens []types.Token `json:"tokens"`
	Ners   []string      `json:"ners"`
}

func analyseRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var msg Req
	err = json.Unmarshal(body, &msg)
	if err != nil {
		panic(err)
	}

	ners := nlp(msg.Text)
	tokens := rakeIt(msg.Text)

	candidates := make([]string, 0)
	for _, ner := range ners {
		//fmt.Printf("%s -> %s\n", ner.Word, ner.Type)
		candidates = append(candidates, ner.Word)
	}

	// for _, token := range tokens {
	// 	fmt.Printf("%s --> %f\n", token.Value, token.Score)
	// }

	resp := &Resp{
		Tokens: tokens,
		Ners:   candidates,
	}
	output, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

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

func analyseSentence(text string) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatal(err)
	}

	for _, tok := range doc.Tokens() {
		fmt.Println(tok.Text, tok.Tag)
	}
}
func analyse(text string) {
	//https://medium.com/@errata.ai/prodigy-prose-radically-efficient-machine-teaching-in-go-93389bf2d772
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatal(err)
	}

	for _, ent := range doc.Entities() {
		fmt.Println(ent.Text, ent.Label)
	}
}

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
