package modules

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/log"
	"os"
	"sync"
	"time"
)

var gracefulTimeout = time.Second * 5

type Module interface {
	Run() error
	Stop() error
	Title() string
}

type Group struct {
	modules []Module
}

type errResp struct {
	err    error
	module string
}

func NewGroup(module ...Module) *Group {
	return &Group{
		modules: module,
	}
}

func (g *Group) Run() {
	Errors := make(chan errResp, len(g.modules))
	for _, m := range g.modules {
		go func(m Module) {
			err := m.Run()
			errResp := errResp{
				err:    err,
				module: m.Title(),
			}
			Errors <- errResp
		}(m)
	}
	// handle errors
	go func() {
		for {
			err := <-Errors
			if err.err != nil {
				log.Error("Module [%s] return error: %s", err.module, err.err)
				g.Stop()
				os.Exit(0)
			}
			log.Info("Module [%s] finish work", err.module)
		}
	}()
}

func (g *Group) Stop() {
	wg := &sync.WaitGroup{}
	wg.Add(len(g.modules))
	for _, m := range g.modules {
		go func(m Module) {
			err := stopModule(m)
			if err != nil {
				log.Error("Module [%s] stopped with error: %s", m.Title(), err.Error())
			}
			wg.Done()
		}(m)
	}
	wg.Wait()
	log.Info("All modules was stopped")
}

func stopModule(m Module) error {
	if m == nil {
		return nil
	}
	result := make(chan error)
	go func() {
		result <- m.Stop()
	}()
	select {
	case err := <-result:
		return err
	case <-time.After(gracefulTimeout):
		return fmt.Errorf("stoped by timeout")
	}
}
