package main

import (
	"context"
	"testing"
	"time"

	"github.com/z-Wind/getNovel/noveler"
)

func TestConcurrentEngine_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		seeds []Request
	}
	tests := []struct {
		name string
		e    *ConcurrentEngine
		args args
		want chan *noveler.NovelChapterHTML
	}{
		// TODO: Add test cases.
		{"Test",
			&ConcurrentEngine{
				Scheduler:   &QueueScheduler{ctx: ctx},
				WorkerCount: 10,
				ctx:         ctx,
			},
			args{[]Request{Request{Order: 1, URL: "https://www.google.com/"}}},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataChan := tt.e.Run(tt.args.seeds...)
			for tt.e.numTasks != 0 {
				select {
				case data := <-dataChan:
					t.Logf("ConcurrentEngine.Run() = %+v, want %v", data, tt.want)
				case <-time.After(time.Second * 1):
					t.Fatal("Timeout")
				}
			}
		})
	}
}
