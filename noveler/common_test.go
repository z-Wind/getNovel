package noveler

import (
	"testing"

	"github.com/z-Wind/getNovel/crawler"
)

func Test_getParseResult(t *testing.T) {
	type args struct {
		novel Noveler
		req   crawler.Request
		reqN  int
	}
	tests := []struct {
		name    string
		args    args
		want    crawler.ParseResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"完本神站",
			args{
				novel: &WanbentxtNoveler{},
				req: crawler.Request{
					Item: NovelChapter{
						URL:   "https://www.wanbentxt.com/8895/5687694.html",
						Order: "0001",
					}},
				reqN: 1,
			},
			crawler.ParseResult{},
			false,
		},
		{
			"小說狂人",
			args{novel: &CzbooksNoveler{},
				req: crawler.Request{
					Item: NovelChapter{
						URL:   "https://czbooks.net/n/u5a6m/uj6h",
						Order: "0001",
					}},
				reqN: 0,
			},
			crawler.ParseResult{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getParseResult(tt.args.novel, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("getParseResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Done || len(got.Requests) != tt.args.reqN {
				t.Errorf("getParseResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeContent(t *testing.T) {
	t.Skip()

	type args struct {
		novelName string
		fileNames []string
		fromPath  string
		toPath    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", args{
			novelName: "瘟疫医生-作者：机器人瓦力.txt",
			fileNames: []string{"1.txt", "2.txt"},
			fromPath:  "./temp",
			toPath:    "./finish",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mergeContent(tt.args.novelName, tt.args.fileNames, tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("mergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
