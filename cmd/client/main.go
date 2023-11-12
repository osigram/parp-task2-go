package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const TestSize = 20000
const StringSize = 20
const Format = "{\"Message\": \"%v\"}"

const Path = "http://localhost:57309" // For go server
//const Path = "http://localhost:5011" // For c# server

//var client = http.DefaultClient

var client = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		WriteBufferSize:     830,
		ReadBufferSize:      830,
	},
}

func mustGenerateTestSlice() [][]byte {
	result := make([][]byte, TestSize)
	someByteSlice := make([]byte, StringSize)
	for i := 0; i < TestSize; i++ {
		for j := 0; j < StringSize; j++ {
			someByteSlice[j] = byte(97 + rand.Intn(24))
		}

		formattedString := fmt.Sprintf(Format, string(someByteSlice[:]))
		result[i] = []byte(formattedString)
	}

	return result
}

func makeRequest(test []byte, ch chan<- error) {
	req, err := http.NewRequest(
		"GET",
		Path,
		bytes.NewReader(test),
	)
	if err != nil {
		fmt.Println("err in creating request: ", err)
		ch <- err
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err in doing request: ", err)
		ch <- err
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		ch <- fmt.Errorf("err in reading body")
	}
	ch <- nil
}

func main() {
	testSlice := mustGenerateTestSlice()
	ch := make(chan error, TestSize)
	finalCh := make(chan int)

	go func() {
		counter := 0
		for i := 0; i < TestSize; i++ {
			if err, ok := <-ch; ok {
				if err == nil {
					counter++
				}
			}
		}
		finalCh <- counter
	}()

	start := time.Now()
	for _, test := range testSlice {
		go makeRequest(test, ch)
	}
	fmt.Println(<-finalCh, "/", TestSize, "successful calls")

	duration := time.Now().Sub(start)
	fmt.Println("That's all. Time:", duration)
}
