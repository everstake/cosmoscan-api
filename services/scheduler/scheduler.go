package scheduler

import (
	"context"
	"github.com/everstake/cosmoscan-api/log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	intervalRunType   runType = "interval"
	periodRunType     runType = "period"
	everyDayRunType   runType = "every_day"
	everyMonthRunType runType = "every_month"
)

type (
	runType string
	Process func()
	task    struct {
		runType  runType
		process  Process
		duration time.Duration
		atTime   atTime
	}
	atTime struct {
		day     int
		hours   int
		minutes int
	}
	Scheduler struct {
		tskCh      chan task
		wg         *sync.WaitGroup
		ctx        context.Context
		cancel     context.CancelFunc
		mu         *sync.RWMutex
		tasks      []task
		alreadyRun bool
	}
)

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		ctx:    ctx,
		cancel: cancel,
		tskCh:  make(chan task),
		wg:     &sync.WaitGroup{},
		mu:     &sync.RWMutex{},
		tasks:  make([]task, 0),
	}
}

func (sch *Scheduler) AddProcessWithInterval(process Process, interval time.Duration) {
	tsk := task{
		runType:  intervalRunType,
		process:  process,
		duration: interval,
	}
	sch.addTask(tsk)
}

func (sch *Scheduler) AddProcessWithPeriod(process Process, period time.Duration) {
	tsk := task{
		runType:  periodRunType,
		process:  process,
		duration: period,
	}
	sch.addTask(tsk)
}

func (sch *Scheduler) EveryDayAt(process Process, hours int, minutes int) {
	tsk := task{
		runType: everyDayRunType,
		process: process,
		atTime: atTime{
			hours:   hours,
			minutes: minutes,
		},
	}
	sch.addTask(tsk)
}

func (sch *Scheduler) EveryMonthAt(process Process, day int, hours int, minutes int) {
	tsk := task{
		runType: everyMonthRunType,
		process: process,
		atTime: atTime{
			day:     day,
			hours:   hours,
			minutes: minutes,
		},
	}
	sch.addTask(tsk)
}

func (sch *Scheduler) Run() error {
	sch.markAsAlreadyRun()
	for _, t := range sch.tasks {
		sch.runTask(t)
	}
	for {
		select {
		case <-sch.ctx.Done():
			return nil
		case t := <-sch.tskCh:
			sch.runTask(t)
		}
	}
}

func (sch *Scheduler) Stop() error {
	sch.cancel()
	sch.wg.Wait()
	return nil
}

func (sch *Scheduler) Title() string {
	return "Scheduler"
}

func (sch *Scheduler) runTask(t task) {
	switch t.runType {
	case intervalRunType:
		go func() {
			runByInterval(sch.ctx, t.process, t.duration)
			sch.wg.Done()
		}()
	case periodRunType:
		go func() {
			runByPeriod(sch.ctx, t.process, t.duration)
			sch.wg.Done()
		}()
	case everyDayRunType:
		go func() {
			runEveryDayAt(sch.ctx, t.process, t.atTime)
			sch.wg.Done()
		}()
	case everyMonthRunType:
		go func() {
			runEveryMonthAt(sch.ctx, t.process, t.atTime)
			sch.wg.Done()
		}()
	}
	log.Debug("Scheduler run process %s", t.process.GetName())
}

func (sch *Scheduler) addTask(tsk task) {
	sch.wg.Add(1)
	if !sch.isAlreadyRun() {
		sch.mu.Lock()
		sch.tasks = append(sch.tasks, tsk)
		sch.mu.Unlock()
		return
	}
	sch.tskCh <- tsk
}

func (sch *Scheduler) isAlreadyRun() bool {
	sch.mu.RLock()
	defer sch.mu.RUnlock()
	return sch.alreadyRun
}

func (sch *Scheduler) markAsAlreadyRun() {
	sch.mu.Lock()
	sch.alreadyRun = true
	sch.mu.Unlock()
}

func runByInterval(ctx context.Context, process Process, interval time.Duration) {
	if interval == 0 {
		log.Error("Scheduler: interval is zero, process %s", process.GetName())
		return
	}
	for {
		process()
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			continue
		}
	}
}

func runByPeriod(ctx context.Context, process Process, period time.Duration) {
	if period == 0 {
		log.Error("Scheduler: period is zero, process %s", process.GetName())
		return
	}
	periodCh := time.After(period)
	for {
		periodCh = time.After(period)
		process()
		select {
		case <-ctx.Done():
			return
		case <-periodCh:
			continue
		}
	}
}

func runEveryDayAt(ctx context.Context, process Process, a atTime) {
	for {
		now := time.Now()
		year, month, day := now.Date()
		today := time.Date(year, month, day, a.hours, a.minutes, 0, 0, time.Local)
		var duration time.Duration
		if today.After(now) {
			duration = today.Sub(now)
		} else {
			tomorrow := today.Add(time.Hour * 24)
			duration = tomorrow.Sub(now)
		}
		next := time.After(duration)
		select {
		case <-ctx.Done():
			return
		case <-next:
			process()
		}
	}
}

func runEveryMonthAt(ctx context.Context, process Process, a atTime) {
	for {
		now := time.Now()
		year, month, _ := now.Date()
		timeInCurrentMonth := time.Date(year, month, a.day, a.hours, a.minutes, 0, 0, time.Local)
		var duration time.Duration
		if timeInCurrentMonth.After(now) {
			duration = timeInCurrentMonth.Sub(now)
		} else {
			nextMonth := timeInCurrentMonth.AddDate(0, 1, 0)
			duration = nextMonth.Sub(now)
		}
		next := time.After(duration)
		select {
		case <-ctx.Done():
			return
		case <-next:
			process()
		}
	}
}

func (p Process) GetName() string {
	path := runtime.FuncForPC(reflect.ValueOf(p).Pointer()).Name()
	if path == "" {
		return path
	}
	parts := strings.Split(path, ".")
	if len(path) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
