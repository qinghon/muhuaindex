package main

import (
	"log"
	"time"
)

func newtk()  {
	tickerm :=time.NewTicker(time.Second*time.Duration(TEST_TIME))
	tickerh :=time.NewTicker(time.Second*time.Duration(UPLOAD_TIME))
	var  AList []All

	err:=get_domain_list()
	if err!=nil {
		ticker3s:=time.NewTicker(time.Second*3)
		log.Println(err)

		for err!=nil  {
			select {
			case <-ticker3s.C:
				err=get_domain_list()
				log.Println(err)
			}
		}
	}
	for   {
		select {
		case <-tickerm.C:
			log.Print("Test domain begin")
			a:=worker()
			AList=append(AList,a)
		case <-tickerh.C:
			A:=count_number(AList)
			upload("1",A)
			AList=nil

		}
	}
}