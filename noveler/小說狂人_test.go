package noveler

import (
	"io"
	"net/http"
	"testing"

	"github.com/z-Wind/getNovel/util"
)

func TestCzbooksNoveler_GetChapterURLs(t *testing.T) {
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		want    CzbooksNoveler
		wantErr bool
	}{
		// TODO: Add test cases.
		{"瘟疫医生", &CzbooksNoveler{URL: "https://czbooks.net/n/u5a6m"}, CzbooksNoveler{title: "《原來我是妖二代》", author: "賣報小郎君"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.GetChapterURLs()
			if (err != nil) != tt.wantErr {
				t.Errorf("CzbooksNoveler.getChapterURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 || tt.n.title != tt.want.title || tt.n.author != tt.want.author {
				t.Errorf("CzbooksNoveler = %v, want %v", tt.n, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}

func TestCzbooksNoveler_GetText(t *testing.T) {
	type args struct {
		html io.Reader
	}
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"test", &CzbooksNoveler{URL: "https://czbooks.net/n/u5a6m/uj6h"}, args{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := http.Get(tt.n.URL)
			r, _, _, _ := util.ToUTF8Encoding(resp.Body)
			resp.Body.Close()
			tt.args.html = r

			got, err := tt.n.GetText(tt.args.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("CzbooksNoveler.getText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("CzbooksNoveler.getText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCzbooksNoveler_MergeContent(t *testing.T) {
	t.Skip()

	type args struct {
		fromPath string
		toPath   string
	}
	tests := []struct {
		name    string
		n       *CzbooksNoveler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", &CzbooksNoveler{title: "《原來我是妖二代》", author: "賣報小郎君", numPages: 10}, args{fromPath: "./temp", toPath: "./finish"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.MergeContent(tt.args.fromPath, tt.args.toPath); (err != nil) != tt.wantErr {
				t.Errorf("CzbooksNoveler.MergeContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
