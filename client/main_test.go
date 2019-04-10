package main

import (
	"log"
	"testing"
)

func Test_runcurl(t *testing.T) {
	b := runcurl("sogou.com")
	t.Log(b)
}

func Test_rundig(t *testing.T) {
	ula := "http://baidu.com"

	a := rundig(ula)
	if a == 0 {
		t.Errorf("query fail")
	}
	t.Log(a)
}
func Test_worker(t *testing.T) {
	log.Println(worker())
}
