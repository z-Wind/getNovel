package util

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	// "log"
	"os"
	"testing"

	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// URLHTMLToUTF8Encoding 將網頁編碼為 UTF8 並回傳 reader
func URLHTMLToUTF8Encoding(URL string) (io.Reader, string, bool, error) {
	var body io.Reader

	if testing.Testing() {
		urlAfter, found := strings.CutPrefix(URL, "https://")
		if !found {
			panic("url not https")
		}
		filename := "../test_dataset/" + urlAfter
		if strings.HasSuffix(urlAfter, "/") {
			filename += "index.html"
		} else {
			fstat, err := os.Stat(filename)
			if err != nil {
				return nil, "", false, err
			}
			if fstat.IsDir() {
				filename += "/index.html"
			}
		}
		// aa, _ := os.Getwd()
		// log.Println(filename, aa)
		file, err := os.Open(filename)
		defer func() {
			_ = file.Close()
		}()
		if err != nil {
			return nil, "", false, err
		}
		body = file
	} else {
		// Create a new context with a deadline
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		resp, err := HTTPGetwithContext(ctx, URL)
		if err != nil {
			err = errors.Wrap(err, "HTTPGetwithContext")
			return nil, "", false, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("response status code: %d", resp.StatusCode)
			return nil, "", false, err
		}
		body = resp.Body
	}

	r, name, certain, err := ToUTF8Encoding(body)
	if err != nil {
		err = errors.Wrap(err, "ToUTF8Encoding")
		return nil, "", false, err
	}

	return r, name, certain, nil
}

// ToUTF8Encoding 將 reader 轉換為 UTF8
func ToUTF8Encoding(r io.Reader) (io.Reader, string, bool, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		err = errors.Wrap(err, "io.ReadAll")
		return nil, "", false, err
	}

	e, name, certain, err := DetermineEncodingFromReader(bytes.NewReader(b))
	if err != nil {
		err = errors.Wrap(err, "DetermineEncodingFromReader")
		return nil, "", false, err
	}

	t := transform.NewReader(bytes.NewReader(b), e.NewDecoder())
	return t, name, certain, nil
}

// DetermineEncodingFromReader 偵測 reader 的編碼
func DetermineEncodingFromReader(r io.Reader) (encoding.Encoding, string, bool, error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		fmt.Printf("r: %s : %s", r, err)
		err = errors.Wrap(err, "bufio.NewReader")
		return nil, "", false, err
	}

	e, name, certain := charset.DetermineEncoding(b, "")
	return e, name, certain, nil
}

// HTTPGetwithContext 將 http.Get 加入 context
func HTTPGetwithContext(ctx context.Context, URL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest")
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:57.0) Gecko/20100101 Firefox/57.0")
	req = req.WithContext(ctx)
	// adding connection:close header hoping to get rid
	// of too many files open error. Found this in http://craigwickesser.com/2015/01/golang-http-to-many-open-files/
	// 連線會變慢，需增加 worker 數目
	req.Header.Add("Connection", "close")

	// 確定連結斷開，若對方不斷開仍存活，可能造成 goroutine leakage
	// 連接的客戶端可以持有的最大空閒連接，預設 2
	// http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = -1
	// 設置DisableKeepAlives=true，則會請求的時候自動加上請求頭"Connection", "close"
	// 這樣在服務端響應完後就會立即關閉連接，否則連接將由客戶端關閉
	//http.DefaultTransport.(*http.Transport).DisableKeepAlives = true
	// 以上會造成連線慢，需斟酌使用

	// client := http.DefaultClient
	client := &http.Client{
		Timeout: time.Second * 30,
		// 跳過 https 認證
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	return resp, nil
}

// ToAbsoluteURL 將 url 轉為絕對路徑
func ToAbsoluteURL(baseURI, uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrap(err, "url.Parse")
	}

	base, err := url.Parse(baseURI)
	if err != nil {
		return "", errors.Wrap(err, "url.Parse")
	}

	return base.ResolveReference(u).String(), nil
}

// FormatText 格式化內文
func FormatText(text string) string {
	text = strings.TrimSpace(text)
	// 每行之間空出一行
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, "\n\n")
	// 將句首的空格(半形 or 全形)取代為兩個半形空格
	text = regexp.MustCompile(`(?m)^[ 　 ]*`).ReplaceAllString(text, "  ")
	// 將只有空格的行，移除其空格
	text = regexp.MustCompile(`(?m)^[ 　 ]+$`).ReplaceAllString(text, "")

	return text
}

// MergeTitle 合併標題
func MergeTitle(text, chapterTitle string) string {

	return fmt.Sprintf("%s\n\n%s\n\n\n\n\n", chapterTitle, text)
}
