package main

import (
	"encoding/json"
	"testing"
)

func Test_getconf(t *testing.T)  {
	var c Conf
	conf,err:=c.getconf()

	if err!=nil {
		t.Error(err)
	}
	if ! checktype(conf) {
		t.Error("check error")
	}
	js,_:=json.Marshal(conf)
	t.Log(string(js))
}
func Test_check(t *testing.T)  {
	config:=Conf{"1","2","3","4","5",control{0,""}}
	if config.check()==false{
		t.Error("check error")
	}

}
func Test_initdb(t *testing.T)  {
	var c Conf
	global_config,err:=c.getconf()
	if err !=nil{
		t.Error(err)
	}
	initdb(global_config)
}