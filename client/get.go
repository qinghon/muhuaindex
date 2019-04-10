package main

import (
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"time"
)

type dominList struct {
	Domain  []string `json:"domain"`
	Massage string   `json:"massage"`
}

func get_domain_list() []error {
	API_URL := BASE_API

	requests := gorequest.New()
	res, body, err := requests.Get(API_URL).
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		End()
	if err != nil {
		//log.Println(err)
		return err
	}
	if res.StatusCode != 200 {
		log.Println(res)
	}

	d := dominList{}
	json.Unmarshal([]byte(body), &d)
	domain_list = d.Domain
	return err
}
