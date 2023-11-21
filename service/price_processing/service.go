package price_processing

import (
	"AIVIA/infrastructure/web/binance_client"
	"fmt"
	"log"
	"sync"
)

type Service struct {
	client *binance_client.Client
}

func NewService(client *binance_client.Client) *Service {
	return &Service{client: client}
}

func (s *Service) PriceProcessing() {
	info, err := s.client.GetExchangeInfo()
	if err != nil {
		log.Println(err)
	}
	var pairs []string
	for i := range info.Symbols[:5] {
		pair := info.Symbols[i].BaseAsset + info.Symbols[i].QuoteAsset
		pairs = append(pairs, pair)
	}
	priceChannel := make(chan map[string]string, len(pairs))
	var wg sync.WaitGroup

	for _, pair := range pairs {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			price, err := s.client.GetTickerPrice(p)
			if err != nil {
				log.Println(err)
			}
			priceChannel <- map[string]string{p: price.Price}
		}(pair)
	}

	wg.Wait()
	close(priceChannel)

	for priceData := range priceChannel {
		for symbol, price := range priceData {
			fmt.Println(symbol, price)
		}
	}
}
