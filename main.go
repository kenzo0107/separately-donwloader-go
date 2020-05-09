package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	procs     int    = runtime.NumCPU()
	body             = make([]string, procs+1)
	targetURL string = "http://kenzo0107.github.io/"
	filename  string = "index.html"
	dst       string = "tmp"
)

func main() {
	if err := download(); err != nil {
		log.Fatal(err)
	}
}

func _main() {

}

func download() error {
	// ダウンロードする URL からファイルサイズ取得
	filesize, err := filesize(targetURL)
	if err != nil {
		return err
	}

	// ファイルサイズを並行処理数で割ったサイズ
	chunk := filesize / procs

	// ファイルサイズを並行処理数で割った余り
	lastRequestSize := filesize % procs

	// goroutine で並行して range access でダウンロード
	bc := context.Background()
	eg, ctx := errgroup.WithContext(bc)
	for i := 0; i < procs; i++ {
		i := i
		// bytes=<min>-<max>
		min := chunk * i
		max := chunk * (i + 1)

		if i == procs-1 {
			max += lastRequestSize
		}
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				fmt.Println("canceled")
				return nil
			default:
				return rangeAccessRequest(min, max, i, targetURL)
			}
		})
	}

	// eg.Go() でエラーで一番最初のエラーを返す
	if err := eg.Wait(); err != nil {
		return err
	}

	// Range Access でダウンロードした分割ファイルをつなぎ合わせる
	if err := bindwithFiles(filename, filesize); err != nil {
		return err
	}
	return nil
}

func filesize(targetURL string) (int, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return 0, err
	}
	res, err := http.Head(u.String())
	if err != nil {
		return 0, err
	}
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0])
	return length, nil
}

// Range Access ダウンロード
func rangeAccessRequest(min, max, i int, url string) error {
	req, _ := http.NewRequest("GET", url, nil)
	rangeBytes := fmt.Sprintf("bytes=%d-%d", min, max)
	req.Header.Add("Range", rangeBytes)

	httpClient := new(http.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	reader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body[i] = string(reader)
	fpath := filepath.Join(dst, strconv.Itoa(i))
	if err := ioutil.WriteFile(fpath, []byte(string(body[i])), 0x777); err != nil {
		return err
	}
	return nil
}

// Ragen Access 分割ダウンロードしたファイルを1つにまとめる
func bindwithFiles(filename string, filesize int) error {
	fh, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create a file in download location")
	}
	defer func() {
		err = fh.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	var f string
	for i := 0; i < procs; i++ {
		i := i
		f = filepath.Join(dst, strconv.Itoa(i))
		subfp, err := os.Open(f)
		if err != nil {
			return errors.Wrap(err, "failed to open "+f+" in download location")
		}
		io.Copy(fh, subfp)

		// Not use defer
		subfp.Close()

		// remove a file in download location for join
		if err := os.Remove(f); err != nil {
			return errors.Wrap(err, "failed to remove a file in download location")
		}
	}
	return nil
}
