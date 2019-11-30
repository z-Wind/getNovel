package crawler

// Request 需執行的任務
type Request struct {
	Item      interface{}
	ParseFunc func(Request) (ParseResult, error)
}

// ParseResult worker 回傳的執行結果
type ParseResult struct {
	Item     interface{}
	Requests []Request
	// 已執行完的任務數，用來扣除用
	DoneN int
}

// Scheduler 調配工作
type Scheduler interface {
	Submit(Request)
	WorkerReady(chan Request)
	Run()
}
