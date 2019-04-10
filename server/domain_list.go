package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"strings"
)

type List_domain struct {
	Domain []string `json:"domain"`
	Massage string `json:"massage"`
}

func (l *List_domain) check() *List_domain {
	d:=List_domain{}
	for _,v:=range l.Domain {
		if strings.Contains(v, "http") {
			u, _ := url.Parse(v)
			d.Domain =append(d.Domain,u.Host)
		} else {
			d.Domain=append(d.Domain,v)
		}
	}
	return &d
}


func get_list() []string {
	rows, err := db.Query("SELECT domain from Domain_list ")
	if err != nil {
		return nil
	}
	defer rows.Close()

	d := []string{}
	for rows.Next() {
		var domain string
		rows.Scan( &domain)
		d = append(d, domain)
	}
	return d
}
func get_list_json(c *gin.Context)  {

	d :=List_domain{}
	d.Domain=get_list()
	if d.Domain==nil {
		c.String(503,"server error")
	}else {
		d.Massage="ok"
		c.JSON(200,d)
		//log.Print(d.Domain)
	}
}


func add_list(data []string) error {

	data=deduplication(data,get_list())
	vals := []interface{}{}
	sqlStr:="INSERT into Domain_list (domain) VALUES "
	for _,v:=range data {
		sqlStr+="(?),"
		vals=append(vals,v)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	stmt,err:=db.Prepare(sqlStr)
	if err!=nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(vals...)
	if err!=nil {
		return err
	}
	return nil
}

func add_list_json(c *gin.Context)  {
	var d  List_domain
	err:=c.ShouldBind(&d)
	if err!=nil {
		log.Print(err)
	}
	if d.Domain!=nil {
		d=*d.check()
		err=add_list(d.Domain)
		if err!=nil {
			c.JSON(503,gin.H{"massage":"add fail","domain":[]string{}})
		}else {
			d.Domain=get_list()
			c.JSON(200,d)
		}
	}else {
		c.JSON(503,gin.H{"massage":false,"domain":[]string{}})
	}


}
func index(a string,l *[]string) int {
	for k,v:=range *l {
		if a==v {
			return k
		}
	}
	return -1
}

func deduplication(a ,b []string) []string {
	d:=[]string{}
	for _,v:=range a {
		if index(v,&b)==-1 {
			d=append(d,v)
		}
	}
	return d

}