package binance_client

import (
	"AIVIA/dto"
	"encoding/json"
	"fmt"
	"github.com/aiviaio/go-binance/v2"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	client         *binance.Client
	spotBaseURL    string
	futuresBaseURL string
}

func NewClient(client *binance.Client, spotEndpoint string, futuresEndpoint string) *Client {
	return &Client{
		client:         client,
		spotBaseURL:    spotEndpoint,
		futuresBaseURL: futuresEndpoint,
	}
}

func (client *Client) GetExchangeInfo() (dto.ExchangeInfo, error) {
	var exchangeInfo dto.ExchangeInfo
	res, err := client.client.HTTPClient.Get(fmt.Sprintf("%s/api/v3/exchangeInfo", client.spotBaseURL))
	if err != nil {
		return dto.ExchangeInfo{}, err
	}
	if res.StatusCode != http.StatusOK {
		binErr := handleBinanceError(res)
		return dto.ExchangeInfo{}, binErr
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return dto.ExchangeInfo{}, err
	}

	if err = json.Unmarshal(body, &exchangeInfo); err != nil {
		return dto.ExchangeInfo{}, err
	}
	return exchangeInfo, nil
}

func (client *Client) GetTickerPrice(symbol string) (dto.Price, error) {
	var tickerPrice dto.Price
	url := fmt.Sprintf("%s/fapi/v2/ticker/price?symbol=%s", client.futuresBaseURL, strings.ToLower(symbol))
	res, err := client.client.HTTPClient.Get(url)
	if err != nil {
		return dto.Price{}, err
	}
	if res.StatusCode != http.StatusOK {
		binErr := handleBinanceError(res)
		return dto.Price{}, fmt.Errorf("%v for %s %s", ErrGetPrice, symbol, binErr)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return dto.Price{}, err
	}

	if err = json.Unmarshal(body, &tickerPrice); err != nil {
		return dto.Price{}, err
	}

	return tickerPrice, nil
}

func handleBinanceError(response *http.Response) error {
	var binanceError dto.BinanceError
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &binanceError); err != nil {
		return err
	}
	return fmt.Errorf("reason: %s", binanceError.Msg)
}
