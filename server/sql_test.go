package main

import (
	"encoding/json"
	"testing"
)

func Test_insert_time(t *testing.T)  {
	initdb()
	a:=[]map[string]interface{}{}
	row:=make(map[string]interface{})
	row["domain"]="baidu.com"
	row["time_dns"]=0.062
	row["time_first_package"]=0.102
	row["time_total"]=0.240
	row["IP"]="1.1.1.1"
	row1:=make(map[string]interface{})
	row1["domain"]="www.baidu.com"
	row1["time_dns"]=0.062
	row1["time_first_package"]=0.102
	row1["time_total"]=0.240
	row1["IP"]="2.2.2.2"
	a=append(a,row,row1)
	err:=insert_time(a);
	if err != nil {
		t.Errorf("%s",err)
	}
	db.Close()
}
func Test_query_time(t *testing.T)  {
	initdb()
	_,err:=query_time("qq.com")
	if err!=nil {
		t.Errorf("query fail:%s",err)
	}
	db.Close()
}
func Test_query_city(t *testing.T)  {
	initall()
	res,err:=query_time("qq.com")
	if err!=nil {
		t.Errorf("query fail:%s",err)
	}
	db.Close()
	_,err=query_city(res)
	if err!=nil {
		t.Errorf("query city error:%s",err)
	}
}

func Test_countScore(t *testing.T)  {
	c:=map[string][]Data_For_Query{}
	c["局域网"]=[]Data_For_Query{
		Data_For_Query{0.20, 0.292, 1.502, "0.0.0.0", "200",0,0,0,0,0,0},
		Data_For_Query{0.20, 0.152, 0.902, "0.0.0.0", "200",0,0,0,0,0,0},
		Data_For_Query{0.20, 0.122, 0.852, "0.0.0.0", "200",0,0,00,0,0,0,},
		Data_For_Query{0.190, 0.102, 0.702, "0.0.0.0", "200",0,0,0,0,0,0},
	}
	//t.Log(c)
	res:=countScore(c)
	js,_:=json.Marshal(res)
	t.Log(string(js))
}

func Test_ipip_query(t *testing.T)  {
	res,err:=ipip_query("182.91.170.36")
	if  err!=nil {
		t.Error(err)
	}
	t.Log(res)
}
/*
func Test_geoip_query(t *testing.T)  {
	res,err:=geoip_query("182.91.170.36")
	if err!=nil {
		t.Error(err)
	}
	t.Log(res)
}
*/