package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/sahandhnj/text-analyser/src/types"
)

const TOKEN = "Cjzvcc3sYcm12Eye2r4pQsS2pezphs"

func analyseRequest(w http.ResponseWriter, r *http.Request) {
	token, ok := r.URL.Query()["token"]

	if !ok || len(token[0]) < 1 || token[0] != TOKEN {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Valid token is required"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wring!"))
		return
	}

	var req types.Req
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wring!"))
		return
	}
	req.Parsed = strip.StripTags(req.Text)

	ners := nlp(req.Parsed, req.Lang)
	tokens := rakeIt(req.Parsed, req.Lang)

	candidates := make([]string, 0)
	for _, ner := range ners {
		//fmt.Printf("%s -> %s\n", ner.Word, ner.Type)
		candidates = append(candidates, ner.Word)
	}

	// for _, token := range tokens {
	// 	fmt.Printf("%s --> %f\n", token.Value, token.Score)
	// }

	resp := &types.Resp{
		Tokens:     tokens,
		UniqeNames: uniqueString(candidates),
		Text:       req.Parsed,
	}
	output, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wring!"))
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
	return
}

func uniqueString(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}
