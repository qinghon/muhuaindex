package main

import (
	"github.com/parnurzeal/gorequest"
	"log"
)

/*
type real_time struct {
	Time_DNS float64 			`json:"time_dns"` //dns 解析时间
	Time_first_package float64 	`json:"time_tcp_connect"` //首包时间
	TimeTotal float64 			`json:"time_total"` //总时间
}
type number_and_time struct {
	Domain string `json:"domain"`
	Number float64 `json:"number"`
	Time real_time `json:"time"`
}

func newnumber_and_time(url string) *number_and_time {
	var a number_and_time
	a.Domain=url
	return &a
}
type up_js struct {
	
	data []number_and_time
}
func count_time(t A_time) real_time {
	var a real_time
	a.Time_DNS=t.TimeNamelookup
	a.Time_first_package=t.TimePretransfer
	a.TimeTotal=t.TimeTotal

	return a
}
*/
func count_number(a []All) All {
	var uploadData All
	uploadData.Id = a[0].Id
	l := len(a)
	log.Printf("Len data list:%v", l)
	b := B_web{}
	c := make(map[string]B_web)
	for v, _ := range a[0].Data {
		b.Url = v

		for _, d := range a {
			b.Time.TimeTotal += d.Data[v].Time.TimeTotal
			b.Time.TimePretransfer += d.Data[v].Time.TimePretransfer
			b.Time.TimeNamelookup += d.Data[v].Time.TimeNamelookup
			b.Time.SpeedDownload += d.Data[v].Time.SpeedDownload
			b.Time.TimeConnect += d.Data[v].Time.TimeConnect
			b.Time.TimeRedirect += d.Data[v].Time.TimeRedirect
			b.Time.TimeStarttransfer += d.Data[v].Time.TimeStarttransfer
			b.Time.HTTPCode = d.Data[v].Time.HTTPCode
			b.Time.NUMREDIRECTS += d.Data[v].Time.NUMREDIRECTS
			b.Time.Sizedownload += d.Data[v].Time.Sizedownload
		}
		b.Time.NUMREDIRECTS /= int8(l)
		b.Time.TimeStarttransfer /= float64(l)
		b.Time.TimeRedirect /= float64(l)
		b.Time.TimeConnect /= float64(l)
		b.Time.SpeedDownload /= float64(l)
		b.Time.TimeTotal /= float64(l)
		b.Time.TimePretransfer /= float64(l)
		b.Time.TimeNamelookup /= float64(l)
		b.Time.Sizedownload /= l
		c[v] = b
		b = B_web{}
		uploadData.Data = c
	}

	return uploadData
}

func upload(id string, a All) {
	url := BASE_API
	req := gorequest.New()
	_, body, err := req.Post(url).Send(a).End()
	if err != nil {
		log.Println(err)
	}
	log.Println(body)
}
