package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const ThCount = 6

func ExecutePipeline(freeFlowJobs ...job) {
	fmt.Println("Start ExecutePipeline")

	wg := &sync.WaitGroup{}
	in := make(chan interface{}, 100)

	for n, someJob := range freeFlowJobs {
		out := make(chan interface{}, 100)
		wg.Add(1)

		fmt.Println("in main ex", n, someJob, in, out)
		go func(in, out chan interface{}, someJob job, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)

			someJob(in, out)
		}(in, out, someJob, wg)
		in = out
	}
	wg.Wait()
}

var SingleHash = func(in, out chan interface{}) {
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for data := range in {
		wg.Add(1)
		strVal := strconv.Itoa(data.(int))
		go calculateSingleHash(strVal, out, mu, wg)
	}

	wg.Wait()
}

func getDataSignerCrc32(data string,hashResults []string, pos int, mu *sync.Mutex, wg *sync.WaitGroup) {
		defer wg.Done()

		result := DataSignerCrc32(data)
		mu.Lock()
		hashResults[pos] = result
		mu.Unlock()
}

func calculateSingleHash(data string, out chan interface{}, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	wgSh := &sync.WaitGroup{}
	muSh := &sync.Mutex{}
	hashedResults := make([]string, 2)

	mu.Lock()
	md5Hash := DataSignerMd5(data)
	mu.Unlock()

	wgSh.Add(2)
	go getDataSignerCrc32(data, hashedResults, 0, muSh, wgSh)
	go getDataSignerCrc32(md5Hash, hashedResults, 1, muSh, wgSh)
	wgSh.Wait()

	out <- strings.Join(hashedResults, "~")
}

func calculateMultiHash(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	hashResults := make([]string, ThCount)
	wg = &sync.WaitGroup{}

	for th := 0; th < ThCount; th += 1 {
		strTh := strconv.Itoa(th)
		mu := &sync.Mutex{}
		wg.Add(1)
		go getDataSignerCrc32(strTh+ data, hashResults, th, mu, wg)
	}
	wg.Wait()
	out <- strings.Join(hashResults, "")

}

var MultiHash = func(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		wg.Add(1)
		go calculateMultiHash(data.(string), out, wg)
	}
	wg.Wait()
}

var CombineResults = func(in, out chan interface{}) {
	var results []string
	for data := range in {
		results = append(results, data.(string))
	}
	sort.Strings(results)
	totalResult := strings.Join(results, "_")
	out <- totalResult
}
