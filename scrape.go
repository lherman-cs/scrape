package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func download(dir, url string) {
	file := path.Base(url)
	path := filepath.Join(dir, file)

	resp, err := http.Get(url)
	must(err)

	writer, err := os.Create(path)
	must(err)
	defer writer.Close()

	reader := bufio.NewReader(resp.Body)
	reader.WriteTo(writer)
}

func downloadWeb(url string) {
	dir := path.Base(url)
	os.Mkdir(dir, 0775)

	resp, err := http.Get(url)
	must(err)

	body, err := ioutil.ReadAll(resp.Body)
	must(err)

	s := string(body)
	exp := regexp.MustCompile("<a href=\"(.+?\\..+?|Makefile)\">")
	results := exp.FindAllStringSubmatch(s, -1)

	var wg sync.WaitGroup
	for _, result := range results {
		wg.Add(1)
		go func(url string) {
			download(dir, url)
			wg.Done()
		}(url + "/" + result[1])
	}
	wg.Wait()
}

func downloadWebs(urls ...string) {
	bar := pb.StartNew(len(urls))
	for _, url := range urls {
		fmt.Println("Scraping", url)
		bar.Increment()
		downloadWeb(url)
	}
	bar.FinishPrint("Finish!")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: scrape [url...]")
	}
	downloadWebs(os.Args[1:]...)
}
