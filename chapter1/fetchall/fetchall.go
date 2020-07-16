package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // start a goroutine
	}

	var messages []string
	for range os.Args[1:] {
		messages = append(messages, <-ch) // receive from channel ch
	}

	result := fmt.Sprintf("%.2fs elapsed time\n", time.Since(start).Seconds())
	messages = append(messages, result)
	saveResult(messages)

}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // Returning error through channel
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // Don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%2.fs %7d %s", secs, nbytes, url)
}

func saveResult(messages []string) {
	file, err := os.Create("results.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	file.Write([]byte(strings.Join(messages, "\n")))
}
