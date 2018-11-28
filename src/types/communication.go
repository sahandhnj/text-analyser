package types

type Req struct {
	Text   string `json:"text"`
	Parsed string
}

type Resp struct {
	Tokens     []Token  `json:"tokens"`
	UniqeNames []string `json:"uniqueNames"`
	Text       string   `json:"text"`
}
