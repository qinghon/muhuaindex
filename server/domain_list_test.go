package main

import "testing"

func Test_get_list(t *testing.T)  {
	initall()
	res:=get_list()
	if res==nil {
		t.Error("this fail ")
	}else {
		t.Log("this ok",res)
	}
	db.Close()
}
func Benchmark_get_list(b *testing.B) {
	b.StopTimer()
	initall()
	b.N=10000
	b.StartTimer()
	for i:=0;i<b.N ; i++ {
		get_list()
	}
	defer db.Close()
}