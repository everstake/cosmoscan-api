package services

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"time"
)

type nodeSize struct {
	Data struct {
		Result []struct {
			Values [][]decimal.Decimal `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func (s *ServiceFacade) GetSizeOfNode() (size float64, err error) {
	to := time.Now()
	from := to.Add(-time.Hour * 24)
	step := time.Hour / time.Second
	params := fmt.Sprintf("&start=%d&end=%d&step=%d", from.Unix(), to.Unix(), step)
	url := "https://eosmon.everstake.one/api/datasources/proxy/1/api/v1/query_range?query=cosmos_size_of_db%20%7B%7D" + params
	resp, err := http.Get(url)
	if err != nil {
		return size, fmt.Errorf("http.Get: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return size, fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	var nSize nodeSize
	err = json.Unmarshal(data, &nSize)
	if err != nil {
		return size, fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	if len(nSize.Data.Result) == 0 {
		return size, fmt.Errorf("invalid response (1)")
	}
	if len(nSize.Data.Result[0].Values) == 0 {
		return size, fmt.Errorf("invalid response  (2)")
	}
	if len(nSize.Data.Result[0].Values[len(nSize.Data.Result[0].Values)-1]) != 2 {
		return size, fmt.Errorf("invalid response  (3)")
	}
	size, _ = nSize.Data.Result[0].Values[len(nSize.Data.Result[0].Values)-1][1].Float64()
	return size, nil
}
