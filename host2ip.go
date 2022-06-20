package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var inputFile = flag.String("f", "host.txt", "输入的文件")
var numThreads = flag.Int("t", 100, "协程数")
var outputFile = flag.String("o", "ips.txt", "输出的文件")

type runner struct {
	workerchan       chan string
	wgresolveworkers *sync.WaitGroup
	outputchan       chan string
	wgoutputworker   *sync.WaitGroup
}

func new() (*runner, error) {
	r := runner{
		wgoutputworker:   &sync.WaitGroup{},
		wgresolveworkers: &sync.WaitGroup{},
		outputchan:       make(chan string),
		workerchan:       make(chan string),
	}

	return &r, nil
}

func main() {
	r, err := new()
	if err != nil {
		log.Println(err)
	}

	flag.Parse()

	//获取命令行参数
	numThreads := *numThreads
	inputFile := *inputFile
	outputFile := *outputFile

	//时间
	t := time.Now()
	//contents := readFile(*inputFile)
	//r.runner()
	//开启结果处理
	r.startOutputWork(outputFile)
	//获取输入的协程
	go r.inputWork(inputFile)

	//worker
	for i := 0; i < numThreads; i++ {
		r.wgresolveworkers.Add(1)
		go r.worker()
	}
	r.wgresolveworkers.Wait()
	close(r.outputchan)
	r.wgoutputworker.Wait()

	elapsed := time.Since(t)
	fmt.Println("app elapsed:", elapsed)

}

func (r *runner) worker() {
	defer r.wgresolveworkers.Done()
	for domain := range r.workerchan {
		ns, err := net.LookupHost(domain)
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(ns)
		// r.outputchan <- ns
		for _, item := range ns {
			//fmt.Println(domain, item)
			r.outputchan <- item
		}

	}
}

//从文件中读入数据，写入到workchan
func (r *runner) inputWork(path string) {
	defer close(r.workerchan)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r.workerchan <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

//读取outputChan的内容打印并存入到指定文件
func (r *runner) startOutputWork(path string) {
	r.wgoutputworker.Add(1)
	go r.outputWork(path)
}

//读取outputChan的内容打印并存入到指定文件
func (r *runner) outputWork(path string) {
	defer r.wgoutputworker.Done()

	//创建文件句柄
	foutput, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer foutput.Close()
	w := bufio.NewWriter(foutput)
	defer w.Flush()

	ips := map[string]struct{}{}
	for item := range r.outputchan {
		ips[item] = struct{}{}

	}

	//去重
	for ip := range ips {
		//fmt.Println("ip:", ip)
		fmt.Println(ip)
		w.WriteString(ip + "\n")
	}

}
