package main

import "testing"

func Test_get_domain_list(t *testing.T)  {
	err:=get_domain_list()
	if err!=nil {
		t.Errorf("get_domain_error")
	}
}
