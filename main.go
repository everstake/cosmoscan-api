package main

import (
	"github.com/everstake/cosmoscan-api/api"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services"
	"github.com/everstake/cosmoscan-api/services/modules"
	"github.com/everstake/cosmoscan-api/services/parser/hub3"
	"github.com/everstake/cosmoscan-api/services/scheduler"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := os.Setenv("TZ", "UTC")
	if err != nil {
		log.Fatal("os.Setenv (TZ): %s", err.Error())
	}

	cfg := config.GetConfig()
	d, err := dao.NewDAO(cfg)
	if err != nil {
		log.Fatal("dao.NewDAO: %s", err.Error())
	}

	s, err := services.NewServices(d, cfg)
	if err != nil {
		log.Fatal("services.NewServices: %s", err.Error())
	}

	prs := hub3.NewParser(cfg, d)

	apiServer := api.NewAPI(cfg, s, d)

	sch := scheduler.NewScheduler()

	sch.AddProcessWithInterval(s.UpdateValidatorsMap, time.Minute*10)
	sch.AddProcessWithInterval(s.UpdateProposals, time.Minute*15)
	sch.AddProcessWithInterval(s.UpdateValidators, time.Minute*15)
	sch.EveryDayAt(s.MakeUpdateBalances, 1, 0)
	sch.EveryDayAt(s.MakeStats, 2, 0)

	go s.KeepHistoricalState()

	g := modules.NewGroup(apiServer, sch, prs)
	g.Run()

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	<-interrupt
	g.Stop()

	os.Exit(0)
}
