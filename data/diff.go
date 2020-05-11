package data

type DiffReq struct {
	Type int `json:"type"`
	Text1 string `json:"text1"`
	Text2 string `json:"text2"`
}

type DiffResp struct {
	Html string `json:"html"`
}

const (
	TypeTextExactDiff = iota
	TypeTextDiffByLine
)