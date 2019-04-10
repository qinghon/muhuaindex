package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	initall()
	engine := gin.Default()
	//engine.GET("/", WebRoot)
	engine.POST("/get_list", Handle)  //客户端上传API
	engine.GET("/get_list", get_list_json) //获取网站列表
	engine.POST("/add_list" ,add_list_json)
	engine.GET("/get_score",get_score)
	engine.Run("0.0.0.0:6666")
	defer db.Close()
}
/*
func WebRoot(context *gin.Context) {
	context.String(http.StatusOK, "<h1>hello world</h1>")
}*/
func Handle(c *gin.Context) {
	var a All
	err := c.ShouldBind(&a)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
	go decode_post(a, c.ClientIP())
}

type A_time struct {
	HTTPCode          string   `json:"http_code"`       //http 状态码
	TimeNamelookup    float64 `json:"time_namelookup"` //DNS解析时间
	TimeConnect       float64 `json:"time_connect"`    //TCP建立时间
	TimeRedirect      float64 `json:"time_redirect"`
	TimePretransfer   float64 `json:"time_pretransfer"`
	TimeStarttransfer float64 `json:"time_starttransfer"`
	TimeTotal         float64 `json:"time_total"`     //请求总时间
	SpeedDownload     float64 `json:"speed_download"` //下载速度
	NUMREDIRECTS      int8    `json:"NUM_REDIRECTS"`  //重定向数
	Sizedownload      int     `json:"size_download"`
}

type B_web struct {
	Url     string `json:"url"`
	Time    A_time `json:"time"`
	Massage string `json:"massage"`
}
type All struct {
	Id   string           `json:"id"`
	Data map[string]B_web `json:"data"`
}
