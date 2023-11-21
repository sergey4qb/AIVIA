package dto

type BinanceError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
