package util

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// URLHTMLToUTF8Encoding 將網頁編碼為 UTF8 並回傳 reader
func URLHTMLToUTF8Encoding(URL string) (io.Reader, string, bool, error) {
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

	r, name, certain, err := ToUTF8Encoding(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "ToUTF8Encoding")
		return nil, "", false, err
	}

	return r, name, certain, nil
}

// ToUTF8Encoding 將 reader 轉換為 UTF8
func ToUTF8Encoding(r io.Reader) (io.Reader, string, bool, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		err = errors.Wrap(err, "ioutil.ReadAll")
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

	// 確定連結斷開，若對方不斷開仍存活，可能造成 goroutine leakage
	// 連接的客戶端可以持有的最大空閒連接，預設 2
	// http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = -1
	// 設置DisableKeepAlives=true，則會請求的時候自動加上請求頭"Connection", "close"
	// 這樣在服務端響應完後就會立即關閉連接，否則連接將由客戶端關閉
	//http.DefaultTransport.(*http.Transport).DisableKeepAlives = true
	// 以上會造成連線慢，需斟酌使用

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client.Do")
	}

	return resp, nil
}
