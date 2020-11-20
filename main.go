package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/chenjiandongx/collections"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	timeout      int = 3
	ip           string
	help         bool
	thread_count int
	filepath     string
	outpath      string = "result.txt"
	percision    int    = 4
)

func get_reqeusts_byurl(url string) {
	fmt.Println(url)
}
func check_tcp(wg *sync.WaitGroup, q *collections.Queue, result_list *[]string) {
	// var address_item interface{}
	// var wait *sync.WaitGroup
	// wait=wg
	for !q.IsEmpty() {
		address_item, _ := q.Get()
		address := address_item.(string)
		conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
		var result string
		if err == nil {
			conn.Close()
			result = strings.Join([]string{address, "open"}, ",")
		} else {
			result = strings.Join([]string{address, "close"}, ",")
		}
		fmt.Println(result)
		*result_list = append(*result_list, result)
		defer wg.Done()
	}

	// fmt.Println(address)

}

func ReadFile(filepath string) []string {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var target_list []string
	for scanner.Scan() {
		line_text := scanner.Text()
		// fmt.Println(line_text)
		target_list = append(target_list, line_text)

	}
	return target_list
}

func WriteFile(outpath string, input_data *[]string) {
	file, _ := os.Create(outpath)
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, data := range *input_data {
		_, err := writer.WriteString(data + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
	writer.Flush()
}

func ConvertHMSTime(nanosecond int) string {
	nanosecondString := strconv.Itoa(nanosecond)
	integerString := ""
	decimalString := ""
	second := 0
	indexBound := len(nanosecondString) - 9
	if len(nanosecondString) > 9 {
		integerString += nanosecondString[:indexBound]
		second, _ = strconv.Atoi(integerString)
		integerTime := ""
		// fmt.Println(second)
		days := second / 86400
		if days > 0 {
			integerTime += strconv.Itoa(days) + "d"
			second = second % 86400
		}
		hours := second / 3600
		if hours > 0 {
			integerTime += strconv.Itoa(hours) + "h"
			second = second % 3600
		}
		minutes := second / 60
		if minutes > 0 {
			integerTime += strconv.Itoa(minutes) + "m"
			second = second % 60
		}
		integerTime += strconv.Itoa(second)
		if integerTime == "" {
			integerTime = "0"
		}
		integerString = integerTime
	} else {
		integerString = "0"
	}
	decimalString += nanosecondString[indexBound:percision]
	hsmTime := integerString + "." + decimalString + "s"
	return hsmTime
}

func main() {

	flag.StringVar(&filepath, "f", "target.txt", "target file ip:port")
	flag.BoolVar(&help, "h", false, "show  help")
	flag.IntVar(&thread_count, "t", 10, "thread count")
	flag.IntVar(&timeout, "timeout", 3, "timeout")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	_, err := os.Stat(filepath)
	if err != nil || os.IsNotExist(err) {
		fmt.Printf("%s not exist!\n", filepath)
		os.Exit(0)
	}
	start_time := time.Now().UnixNano()
	target_list := ReadFile(filepath)
	if len(target_list) == 0 {
		fmt.Printf("%s is empty!", filepath)
		os.Exit(0)
	}
	var result_list []string
	var wg sync.WaitGroup
	q := collections.NewQueue()
	for _, address := range target_list {
		wg.Add(1)
		q.Put(address)
	}
	target_size := len(target_list)
	if thread_count > target_size {
		thread_count = target_size
	}
	fmt.Println("target file: " + filepath)
	fmt.Println("current thread count:" + strconv.Itoa(thread_count))
	fmt.Printf("target size:%s\n\n", strconv.Itoa(target_size))
	for i := 0; i < thread_count; i++ {
		go check_tcp(&wg, q, &result_list)
	}
	wg.Wait()
	fmt.Println("---------------------------")
	if len(result_list) > 0 {
		outpath = time.Now().Format("result_2006-01-02_150405") + ".csv"
		fmt.Printf("save to file %s ……\n", outpath)
		WriteFile(outpath, &result_list)
	}
	end_time := time.Now().UnixNano()
	fmt.Printf("Done!elapsed time: %s\n", ConvertHMSTime(int(end_time-start_time)))
}
