package dto

type Price struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}
