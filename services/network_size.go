package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type nodeSize struct {
	Size float64 `json:"Size_of_dir_Gb"`
}

func (s *ServiceFacade) GetSizeOfNode() (size float64, err error) {
	// not public, available only for internal everstake services
	url := "http://s82.everstake.one:8060/monitoring"
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
	return nSize.Size, nil
}
