package main

import (
	"fmt"
	"log"

	prose "gopkg.in/jdkato/prose.v2"
)

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
