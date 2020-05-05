package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	procs     int    = runtime.NumCPU()
	body             = make([]string, procs+1)
	targetURL string = "https://ubuntu.com/download/desktop/thank-you"
	filename  string = "hoge"
)

func filesize() int {
	res, _ := http.Head(targetURL)
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0])
	return length
}

func main() {
	download()
}

func download() error {
	filesize := filesize()
	chunk := filesize / procs
	lastRequestSize := filesize % procs

	// goroutine で並行して range access でダウンロード
	wg := new(sync.WaitGroup)
	for i := 0; i < procs; i++ {
		wg.Add(1)
		min := chunk * i
		max := chunk * (i + 1)

		if i == procs-1 {
			max += lastRequestSize
		}

		go rangeAccessRequest(wg, min, max, i, targetURL)
	}
	wg.Wait()

	return bindwithFiles(filename, filesize)
}

func rangeAccessRequest(wg *sync.WaitGroup, min, max, i int, url string) {
	defer wg.Done()

	req, _ := http.NewRequest("GET", url, nil)
	rangeBytes := fmt.Sprintf("bytes=%d-%d", min, max)
	req.Header.Add("Range", rangeBytes)

	httpClient := new(http.Client)
	resp, errRes := httpClient.Do(req)
	if errRes != nil {
		log.Fatal(errRes)
	}

	defer resp.Body.Close()
	reader, _ := ioutil.ReadAll(resp.Body)
	body[i] = string(reader)
	ioutil.WriteFile(strconv.Itoa(i), []byte(string(body[i])), 0x777)
}

func bindwithFiles(filename string, filesize int) error {
	fh, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "failed to create a file in download location")
	}
	defer fh.Close()

	bar := pb.New64(int64(filesize))
	bar.Start()

	var f string
	for i := 0; i < procs; i++ {
		f = strconv.Itoa(i)
		subfp, err := os.Open(f)
		if err != nil {
			return errors.Wrap(err, "failed to open "+f+" in download location")
		}

		proxy := bar.NewProxyReader(subfp)
		io.Copy(fh, proxy)

		// Not use defer
		subfp.Close()

		// remove a file in download location for join
		if err := os.Remove(f); err != nil {
			return errors.Wrap(err, "failed to remove a file in download location")
		}
	}
	return nil
}
