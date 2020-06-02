package dmodels

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type testTime struct {
	T Time `json:"t"`
}

func TestTime(t *testing.T) {
	tm := time.Unix(1565885014, 0)
	s1 := testTime{T: Time{tm}}
	b, err := json.Marshal(s1)
	if err != nil {
		t.Error(err)
		return
	}
	var s2 testTime
	err = json.Unmarshal(b, &s2)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Error("not equal", s1, s2)
	}
}
