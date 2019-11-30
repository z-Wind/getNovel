package crawler

import (
	"context"
	"fmt"
)

// ConcurrentEngine 負責處理對外與建立 worker
type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	Ctx         context.Context
	NumTasks    int
}

// Run 開始運作
func (e *ConcurrentEngine) Run(seeds ...Request) chan interface{} {
	parseResultChan := make(chan ParseResult)
	dataChan := make(chan interface{})

	e.Scheduler.Run()
	e.NumTasks = len(seeds)
	fmt.Printf("tasks: %d\n", e.NumTasks)

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(parseResultChan, e.Scheduler)
	}

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	go func() {
		// 用 queue 先存起來，防止阻塞
		var dataQ []interface{}

		for {
			var activeData interface{}
			// channel 初值為 nil，並不會觸發 select，除非賦於值
			var activeDataChan chan<- interface{}
			if len(dataQ) > 0 {
				activeData = dataQ[0]
				activeDataChan = dataChan
			}

			select {
			case activeDataChan <- activeData:
				dataQ = dataQ[1:]
			case parseResult := <-parseResultChan:
				if parseResult.Item != nil {
					// fmt.Printf("Get %+v\n", parseResult.Item)
					dataQ = append(dataQ, parseResult.Item)
				}
				e.NumTasks -= parseResult.DoneN
				fmt.Printf("tasks: %d\n", e.NumTasks)

				// 排入新增的 requests
				for _, request := range parseResult.Requests {
					e.Scheduler.Submit(request)
				}
			case <-e.Ctx.Done():
				fmt.Printf("ConcurrentEngine.Run.Done\n")
				return
			}
		}
	}()

	return dataChan
}

func (e *ConcurrentEngine) createWorker(parseResultChan chan<- ParseResult, s Scheduler) {
	requestChan := make(chan Request)

	go func() {
		// 用 queue 先存起來，防止阻塞
		var parseResultQ []ParseResult

		s.WorkerReady(requestChan)

		for {
			var activeResult ParseResult
			// channel 初值為 nil，並不會觸發 select，除非賦於值
			var activeResultChan chan<- ParseResult
			if len(parseResultQ) > 0 {
				activeResult = parseResultQ[0]
				activeResultChan = parseResultChan
			}

			select {
			case activeResultChan <- activeResult:
				parseResultQ = parseResultQ[1:]
			case request := <-requestChan:
				result := worker(request)
				parseResultQ = append(parseResultQ, result)
				s.WorkerReady(requestChan)
			case <-e.Ctx.Done():
				fmt.Printf("ConcurrentEngine.createWorker.Done\n")
				return
			}
		}
	}()
}

func worker(req Request) ParseResult {
	parseResult, err := req.ParseFunc(req)
	if err != nil {
		fmt.Printf("ParseResult: req.ParseFunc: err:%s\n", err)
		return ParseResult{
			Item:     nil,
			Requests: []Request{req},
			DoneN:    0,
		}
	}

	return parseResult
}
