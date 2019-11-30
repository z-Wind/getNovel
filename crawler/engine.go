package main

import (
	"context"
	"fmt"

	"github.com/z-Wind/getNovel/noveler"
	"github.com/z-Wind/getNovel/util"
)

// ConcurrentEngine 負責處理對外與建立 worker
type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ctx         context.Context
	numTasks    int
}

// Run 開始運作
func (e *ConcurrentEngine) Run(seeds ...Request) chan *noveler.NovelChapterHTML {
	parseResultChan := make(chan ParseResult)
	dataChan := make(chan *noveler.NovelChapterHTML)

	e.Scheduler.Run()
	e.numTasks = len(seeds)
	fmt.Printf("tasks: %d\n", e.numTasks)

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(parseResultChan, e.Scheduler)
	}

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	go func() {
		// 用 queue 先存起來，防止阻塞
		var dataQ []*noveler.NovelChapterHTML

		for {
			var activeData *noveler.NovelChapterHTML
			// channel 初值為 nil，並不會觸發 select，除非賦於值
			var activeDataChan chan<- *noveler.NovelChapterHTML
			if len(dataQ) > 0 {
				activeData = dataQ[0]
				activeDataChan = dataChan
			}

			select {
			case activeDataChan <- activeData:
				dataQ = dataQ[1:]
			case parseResult := <-parseResultChan:
				if parseResult.Item != nil {
					fmt.Printf("Get %d: %s\n", parseResult.Item.Order, parseResult.Item.URL)
					dataQ = append(dataQ, parseResult.Item)
				}
				e.numTasks -= parseResult.doneN
				fmt.Printf("tasks: %d\n", e.numTasks)

				// 排入新增的 requests
				for _, request := range parseResult.Requests {
					e.Scheduler.Submit(request)
				}
			case <-e.ctx.Done():
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
			case <-e.ctx.Done():
				fmt.Printf("ConcurrentEngine.createWorker.Done\n")
				return
			}
		}
	}()
}

func worker(req Request) ParseResult {
	// Request the HTML page
	// Create a new context with a deadline
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, noveler.Timeout)
	defer cancel()

	resp, err := util.HTTPGetwithContext(ctx, req.URL)
	if err != nil {
		fmt.Printf("ParseResult: HTTPGetwithContext(%s): %s\n", req.URL, err)
		return ParseResult{
			Item:     nil,
			Requests: []Request{req},
			doneN:    0,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("ParseResult: HTTPGetwithContext(%s): status code error: %d %s\n", req.URL, resp.StatusCode, resp.Status)
		return ParseResult{
			Item:     nil,
			Requests: []Request{req},
			doneN:    0,
		}
	}

	r, name, certain, err := util.ToUTF8Encoding(resp.Body)
	if err != nil {
		fmt.Printf("ParseResult: ToUTF8Encoding: name:%s, certain:%v err:%s\n", name, certain, err)
		return ParseResult{
			Item:     nil,
			Requests: []Request{req},
			doneN:    0,
		}
	}

	return ParseResult{
		Item: &noveler.NovelChapterHTML{
			NovelChapter: &noveler.NovelChapter{
				Order: req.Order,
				URL:   req.URL,
			},
			HTML: r},
		Requests: []Request{},
		doneN:    1,
	}
}
