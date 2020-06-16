package services
// todo delete
//
//import (
//	"fmt"
//	"github.com/everstake/cosmoscan-api/dao/filters"
//	"github.com/everstake/cosmoscan-api/dmodels"
//	"github.com/everstake/cosmoscan-api/log"
//	"github.com/everstake/cosmoscan-api/smodels"
//	"time"
//)
//
//type multiValue struct {
//	Value1d  string
//	Value7d  string
//	Value30d string
//	Value90d string
//}
//
//func (v *multiValue) setValue(i int, s string) {
//	switch i {
//	case 0:
//		v.Value1d = s
//	case 1:
//		v.Value7d = s
//	case 2:
//		v.Value30d = s
//	case 3:
//		v.Value90d = s
//	}
//}
//
//func (s *ServiceFacade) KeepRangeStates() {
//	ranges := []time.Duration{time.Hour * 24, time.Hour * 24 * 7, time.Hour * 24 * 30, time.Hour * 24 * 90}
//	states := []struct {
//		title         string
//		cacheDuration time.Duration
//		fetcher       func() (multiValue, error)
//	}{
//		{
//			title:         dmodels.RangeStateNumberDelegators,
//			cacheDuration: time.Minute * 5,
//			fetcher: func() (value multiValue, err error) {
//				t := time.Now()
//				for i, r := range ranges {
//					total, err := s.dao.GetDelegatorsTotal(filters.TimeRange{From: dmodels.NewTime(t.Add(-r))})
//					if err != nil {
//						return value, fmt.Errorf("dao.GetDelegatorsTotal: %s", err.Error())
//					}
//					value.setValue(i, fmt.Sprintf("%d", total))
//				}
//				return value, nil
//			},
//		},
//		{
//			title:         dmodels.RangeStateNumberMultiDelegators,
//			cacheDuration: time.Minute * 5,
//			fetcher: func() (value multiValue, err error) {
//				t := time.Now()
//				for i, r := range ranges {
//					total, err := s.dao.GetMultiDelegatorsTotal(filters.TimeRange{From: dmodels.NewTime(t.Add(-r))})
//					if err != nil {
//						return value, fmt.Errorf("dao.GetMultiDelegatorsTotal: %s", err.Error())
//					}
//					value.setValue(i, fmt.Sprintf("%d", total))
//				}
//				return value, nil
//			},
//		},
//		{
//			title:         dmodels.RangeStateTransfersVolume,
//			cacheDuration: time.Minute * 5,
//			fetcher: func() (value multiValue, err error) {
//				t := time.Now()
//				for i, r := range ranges {
//					total, err := s.dao.GetTransferVolume(filters.TimeRange{From: dmodels.NewTime(t.Add(-r))})
//					if err != nil {
//						return value, fmt.Errorf("dao.GetTransferVolume: %s", err.Error())
//					}
//					value.setValue(i, total.String())
//				}
//				return value, nil
//			},
//		},
//		{
//			title:         dmodels.RangeStateFeeVolume,
//			cacheDuration: time.Minute * 5,
//			fetcher: func() (value multiValue, err error) {
//				t := time.Now()
//				for i, r := range ranges {
//					total, err := s.dao.GetTransactionsFeeVolume(filters.TimeRange{From: dmodels.NewTime(t.Add(-r))})
//					if err != nil {
//						return value, fmt.Errorf("dao.GetTransactionsFeeVolume: %s", err.Error())
//					}
//					value.setValue(i, total.String())
//				}
//				return value, nil
//			},
//		},
//		{
//			title:         dmodels.RangeStateHighestFee,
//			cacheDuration: time.Minute * 5,
//			fetcher: func() (value multiValue, err error) {
//				t := time.Now()
//				for i, r := range ranges {
//					total, err := s.dao.GetTransactionsHighestFee(filters.TimeRange{From: dmodels.NewTime(t.Add(-r))})
//					if err != nil {
//						return value, fmt.Errorf("dao.GetTransactionsHighestFee: %s", err.Error())
//					}
//					value.setValue(i, total.String())
//				}
//				return value, nil
//			},
//		},
//		{
//			title:         dmodels.RangeStateUndelegationVolume,
//			cacheDuration: time.Minute * 5,
//			fetcher: func() (value multiValue, err error) {
//				t := time.Now()
//				for i, r := range ranges {
//					total, err := s.dao.GetUndelegationsVolume(filters.TimeRange{From: dmodels.NewTime(t.Add(-r))})
//					if err != nil {
//						return value, fmt.Errorf("dao.GetUndelegationsVolume: %s", err.Error())
//					}
//					value.setValue(i, total.String())
//				}
//				return value, nil
//			},
//		},
//		{
//			title:         dmodels.RangeStateBlockDelay,
//			cacheDuration: time.Minute * 5,
//			fetcher: func() (value multiValue, err error) {
//				t := time.Now()
//				for i, r := range ranges {
//					total, err := s.dao.GetAvgBlocksDelay(filters.TimeRange{From: dmodels.NewTime(t.Add(-r))})
//					if err != nil {
//						return value, fmt.Errorf("dao.GetAvgBlocksDelay: %s", err.Error())
//					}
//					value.setValue(i, fmt.Sprintf("%f", total))
//				}
//				return value, nil
//			},
//		},
//	}
//
//	items, err := s.dao.GetRangeStates(nil)
//	if err != nil {
//		log.Error("KeepRangeStates: dao.GetRangeStates: %s", err.Error())
//		return
//	}
//	mp := make(map[string]dmodels.RangeState)
//	for _, item := range items {
//		mp[item.Title] = item
//	}
//
//	for {
//		for _, state := range states {
//			item, found := mp[state.title]
//			if found {
//				if time.Now().Sub(item.UpdatedAt) < state.cacheDuration {
//					continue
//				}
//			}
//			values, err := state.fetcher()
//			if err != nil {
//				log.Error("KeepRangeStates: %s", err.Error())
//				continue
//			}
//			model := dmodels.RangeState{
//				Title:     state.title,
//				Value1d:   values.Value1d,
//				Value7d:   values.Value7d,
//				Value30d:  values.Value30d,
//				Value90d:  values.Value90d,
//				UpdatedAt: time.Now(),
//			}
//			if found {
//				err = s.dao.UpdateRangeState(model)
//				if err != nil {
//					log.Error("KeepRangeStates: UpdateRangeState %s", err.Error())
//				}
//			} else {
//				err = s.dao.CreateRangeState(model)
//				if err != nil {
//					log.Error("KeepRangeStates: CreateRangeState (%s) %s", state.title, err.Error())
//				}
//			}
//			if err != nil {
//				mp[state.title] = model
//			}
//		}
//
//		<-time.After(time.Second * 5)
//	}
//}
//
//func (s *ServiceFacade) GetNetworkStates() (states map[string]smodels.RangeState, err error) {
//	models, err := s.dao.GetRangeStates([]string{
//		dmodels.RangeStateNumberDelegators,
//		dmodels.RangeStateNumberMultiDelegators,
//		dmodels.RangeStateTransfersVolume,
//		dmodels.RangeStateFeeVolume,
//		dmodels.RangeStateHighestFee,
//		dmodels.RangeStateUndelegationVolume,
//		dmodels.RangeStateBlockDelay,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("dao.GetRangeStates: %s", err.Error())
//	}
//	states = make(map[string]smodels.RangeState)
//	for _, model := range models {
//		states[model.Title] = smodels.RangeState{
//			D1:  model.Value1d,
//			D7:  model.Value7d,
//			D30: model.Value30d,
//			D90: model.Value90d,
//		}
//	}
//	return states, nil
//}
