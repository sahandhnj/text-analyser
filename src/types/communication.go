package types

import (
	rake "github.com/sahandhnj/text-analyser/src/rakeimpl"
)

type Req struct {
	Text   string    `json:"text"`
	Lang   rake.LANG `json:"Lang"`
	Parsed string
}

type Resp struct {
	Tokens     []Token  `json:"tokens"`
	UniqeNames []string `json:"uniqueNames"`
	Text       string   `json:"text"`
}
