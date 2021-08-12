package coingecko

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	coingecko "github.com/superoo7/go-gecko/v3"
	"net/http"
	"time"
)

const (
	coinID = "persistence"
)

type CoinGecko struct {
	client *coingecko.Client
}

func NewGecko() *CoinGecko {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return &CoinGecko{
		client: coingecko.NewClient(httpClient),
	}
}

func (g CoinGecko) GetMarketData() (price, volume24h decimal.Decimal, err error) {
	data, err := g.client.CoinsID(coinID, false, true, true, false, false, false)
	if err != nil {
		return price, volume24h, fmt.Errorf("client.CoinsID: %s", err.Error())
	}
	if data.MarketData.MarketCap == nil {
		return price, volume24h, errors.New("MarketData.MarketCap is nil")
	}

	return decimal.NewFromFloat(data.MarketData.CurrentPrice["usd"]), decimal.NewFromFloat(data.MarketData.TotalVolume["usd"]), nil
}
