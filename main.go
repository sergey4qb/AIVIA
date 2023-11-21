package main

import (
	"AIVIA/infrastructure/web/binance_client"
	"AIVIA/service/price_processing"
	"github.com/aiviaio/go-binance/v2"
)

func main() {
	client := binance.NewClient("", "")
	web := binance_client.NewClient(client, "https://api.binance.com", "https://fapi.binance.com")
	info := price_processing.NewService(web)
	info.PriceProcessing()
}
