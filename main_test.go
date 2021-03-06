package main

import (
	"testing"

	"github.com/z-Wind/getNovel/noveler"
)

func Test_getNovel(t *testing.T) {
	t.Skip()

	type args struct {
		novel noveler.Noveler
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"完本神站", args{&noveler.WanbentxtNoveler{URL: "https://www.wanbentxt.com/18868/"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getNovel(tt.args.novel); (err != nil) != tt.wantErr {
				t.Errorf("getNovel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_chooseNoveler(t *testing.T) {
	type args struct {
		URLNovel string
	}
	tests := []struct {
		name    string
		args    args
		check   func(interface{}) bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{"完本神站 mobile", args{URLNovel: "https://m.wanbentxt.com/8895/"}, func(i interface{}) bool { _, ok := i.(noveler.WanbentxtNoveler); return ok }, false},
		{"完本神站", args{URLNovel: "https://www.wanbentxt.com/8895/"}, func(i interface{}) bool { _, ok := i.(noveler.WanbentxtNoveler); return ok }, false},
		{"小說狂人", args{URLNovel: "https://czbooks.net/n/u5a6m"}, func(i interface{}) bool { _, ok := i.(noveler.CzbooksNoveler); return ok }, false},
		{"黃金屋 簡體", args{URLNovel: "https://www.hjwzw.com/Book/Chapter/37176"}, func(i interface{}) bool { _, ok := i.(noveler.HjwzwNoveler); return ok }, false},
		{"黃金屋", args{URLNovel: "https://tw.hjwzw.com/Book/Chapter/37176"}, func(i interface{}) bool { _, ok := i.(noveler.HjwzwNoveler); return ok }, false},
		{"UU看書網", args{URLNovel: "https://www.uukanshu.com/b/81074/"}, func(i interface{}) bool { _, ok := i.(noveler.UUkanshuNoveler); return ok }, false},
		{"UU看書網 TW", args{URLNovel: "https://tw.uukanshu.com/b/81005/"}, func(i interface{}) bool { _, ok := i.(noveler.UUkanshuNoveler); return ok }, false},
		{"飄天文學", args{URLNovel: "https://www.ptwxz.com/html/9/9795/index.html"}, func(i interface{}) bool { _, ok := i.(noveler.PtwxzNoveler); return ok }, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chooseNoveler(tt.args.URLNovel)
			if (err != nil) != tt.wantErr {
				t.Errorf("chooseNoveler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check(got) {
				t.Errorf("chooseNoveler() = %v,  wrong type", got)
			}
		})
	}
}
