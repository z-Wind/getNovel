package main

import (
	"context"
	"fmt"
)

// QueueScheduler 分配 request 給 worker
type QueueScheduler struct {
	workerChan  chan chan Request
	requestChan chan Request
	ctx         context.Context
}

// Submit 提交任務
func (s *QueueScheduler) Submit(r Request) {
	select {
	case s.requestChan <- r:
	case <-s.ctx.Done():
		fmt.Printf("QueueScheduler.Submit.Done\n")
	}
}

// WorkerReady 將空閒的 worker 排進序列
func (s *QueueScheduler) WorkerReady(w chan Request) {
	select {
	case s.workerChan <- w:
	case <-s.ctx.Done():
		fmt.Printf("QueueScheduler.WorkerReady.Done\n")
	}
}

// Run 執行調配
func (s *QueueScheduler) Run() {
	s.requestChan = make(chan Request)
	s.workerChan = make(chan chan Request)

	go func() {
		// 用 queue 先存起來，防止阻塞
		var requestQ []Request
		var workerQ []chan Request

		for {
			var activeRequest Request
			// channel 初值為 nil，並不會觸發 select，除非賦於值
			var activeWorker chan<- Request
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}

			select {
			case activeWorker <- activeRequest:
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
			case <-s.ctx.Done():
				fmt.Printf("QueueScheduler.Run.Done\n")
				return
			}
		}
	}()
}
