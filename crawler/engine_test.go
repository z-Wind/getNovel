package crawler

import (
	"testing"
	"time"
)

func TestConcurrentEngine_Run(t *testing.T) {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	type args struct {
		seeds []Request
	}
	tests := []struct {
		name string
		e    *ConcurrentEngine
		args args
		want chan interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataChan := tt.e.Run(tt.args.seeds...)
			for tt.e.NumTasks != 0 {
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
