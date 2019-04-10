package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var BASE_API string
var domain_list []string

var UPLOAD_TIME int
var TEST_TIME int

type A_time struct {
	HTTPCode          string  `json:"http_code"`       //http 状态码
	TimeNamelookup    float64 `json:"time_namelookup"` //DNS解析时间
	TimeConnect       float64 `json:"time_connect"`    //TCP建立时间
	TimeRedirect      float64 `json:"time_redirect"`
	TimePretransfer   float64 `json:"time_pretransfer"`
	TimeStarttransfer float64 `json:"time_starttransfer"`
	TimeTotal         float64 `json:"time_total"`     //请求总时间
	SpeedDownload     float64 `json:"speed_download"` //下载速度
	NUMREDIRECTS      int8    `json:"NUM_REDIRECTS"`  //重定向数
	Sizedownload      int     `json:"size_download"`
}

type B_web struct {
	Url     string `json:"url"`
	Time    A_time `json:"time"`
	Massage string `json:"massage"`
}
type All struct {
	Id   string           `json:"id"`
	Data map[string]B_web `json:"data"`
}

func runcmd(cmd string) ([]byte, int) {

	a := exec.Command("/bin/sh", "-c", cmd)
	c, err := a.Output()
	if err != nil {
		fmt.Println(err)
	}
	ret := a.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	return c, ret
}

func runcurl(url string) (B_web) {
	var curl = "curl -o /dev/null -sL --connect-timeout 10 -w '{" +
		"\"http_code\":\"%{http_code}\"," +
		"\"time_namelookup\":%{time_namelookup}," +
		"\"time_connect\":%{time_connect}," +
		"\"time_redirect\":%{time_redirect}," +
		"\"time_pretransfer\":%{time_pretransfer}," +
		"\"time_starttransfer\":%{time_starttransfer}," +
		"\"time_total\":%{time_total}," +
		"\"speed_download\":%{speed_download}," +
		"\"NUM_REDIRECTS\":%{NUM_REDIRECTS}," +
		"\"size_download\":%{size_download} }'  " +
		"'" + url + "'"
	res, ret := runcmd(curl)
	var f A_time
	var b B_web
	err := json.Unmarshal(res, &f)
	if err != nil {
		fmt.Println(err, url)
	}
	f.TimeNamelookup = rundig(url)
	b.Time = f
	b.Url = url
	b.Massage = EXIT_CODE[ret]
	return b
}
func rundig(url_link string) float64 {
	var domain string
	if strings.Contains(url_link, "http") {
		u, _ := url.Parse(url_link)
		domain = u.Host
	} else {
		domain = url_link
	}
	cmd := "dig " + domain + " +noall +stat|grep msec|egrep -o '[0-9]+'"
	dnstime, _ := runcmd(cmd)
	flt32, _ := strconv.Atoi(string(dnstime)[:len(dnstime)-1])
	f := float64(flt32) / 1000
	return f
}
func worker_sun(id int, jobs <-chan string, result chan<- B_web) {
	for j := range jobs {
		b := runcurl(j)
		result <- b
		//log.Println(id)
	}
}
func worker() All {
	jobs := make(chan string, 100)
	results := make(chan B_web, 100)

	for i := 0; i < runtime.NumCPU(); i++ {
		go worker_sun(i, jobs, results)
	}
	for _, v := range domain_list {
		jobs <- v
	}
	close(jobs)

	var A All
	d := make(map[string]B_web)
	for i := 0; i < len(domain_list); i++ {
		c := <-results
		url1 := c.Url
		d[url1] = c
		A.Data = d
	}
	return A
}

func main() {
	BASE_API = os.Getenv("API_URL")
	UPLOAD_TIME_s := os.Getenv("UPLOAD_TIME")
	TEST_TIME_s := os.Getenv("TEST_TIME")
	if UPLOAD_TIME_s == "" || TEST_TIME_s == "" {
		TEST_TIME = 600
		UPLOAD_TIME = 3600
		log.Println("TEST_TIME or UPLOAD_TIME not found,set TEST_TIME=", TEST_TIME, " and set UPLOAD_TIME =", UPLOAD_TIME)
	} else {
		TEST_TIME, _ = strconv.Atoi(TEST_TIME_s)
		UPLOAD_TIME, _ = strconv.Atoi(UPLOAD_TIME_s)
	}
	if BASE_API == "" {
		log.Println("API_URL not found")
		os.Exit(1)
	}
	newtk()
}
