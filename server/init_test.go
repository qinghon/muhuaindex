package main

import (
	"encoding/json"
	"testing"
)

func Test_getconf(t *testing.T)  {
	var c Conf
	conf:=c.getconf()
	js,_:=json.Marshal(conf)
	t.Log(string(js))
}
