package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ipipdotnet/ipdb-go"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

type Data_For_Query struct {
	Time_dns           float64 `json:"time_dns"`
	Time_first_package float64 `json:"time_first_package"`
	Time_total         float64 `json:"time_total"`
	IP                 string  `json:"ip"`
	Http_code          string  `json:"http_code"`
	Size_download      int     `json:"size_download"`
	NUM_REDIRECTS      int8    `json:"num_redirects"`
	Speed_download     float64 `json:"speed_download"`
	Time_redirect      float64 `json:"time_redirect"`
	Time_pretransfer   float64 `json:"time_pretransfer"`
	Time_starttransfer float64 `json:"time_starttransfer"`
}
type Signle_data struct {
	Data       map[string]float64 `json:"data"`
	Score      float64            `json:"score"`
	ScoreDns   float64            `json:"score_dns"`
	ScoreFirst float64            `json:"score_first"`
	ScoreTotal float64            `json:"score_total"`
	Status     bool               `json:"status"`
}
type country_data struct {
	CountryScore Signle_data            `json:"country_score"`
	Province     map[string]Signle_data `json:"province"`
}

func initdb() (*sql.DB, error) {
	var err error
	par := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s",
		global_config.DATABASE_USER,
		global_config.DATABASE_PASSWD,
		global_config.DATABASE_HOST,
		global_config.DATABASE_PORT,
		global_config.DATABASE_NAME)

	db, err = sql.Open("mysql", par)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to log mysql: %s ", err)
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(2000)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping mysql: %s", err)
	}
	//db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Hour)
	return db, nil

}

func decode_post(a All, ip string) {
	dict := map[string]interface{}{}
	list_list := []map[string]interface{}{}

	for k, v := range a.Data {
		dict["domain"] = k
		dict["time_dns"] = v.Time.TimeNamelookup
		dict["time_first_package"] = v.Time.TimeStarttransfer - v.Time.TimePretransfer
		dict["time_total"] = v.Time.TimeTotal
		dict["IP"] = ip
		dict["http_code"] = v.Time.HTTPCode
		dict["size_download"] = v.Time.Sizedownload
		dict["NUM_REDIRECTS"] = v.Time.NUMREDIRECTS
		dict["speed_download"] = v.Time.SpeedDownload
		dict["time_redirect"] = v.Time.TimeRedirect
		dict["time_pretransfer"] = v.Time.TimePretransfer
		dict["time_starttransfer"] = v.Time.TimeStarttransfer
		list_list = append(list_list, dict)
		dict = map[string]interface{}{}
	}
	insert_time(list_list)
}

func insert_time(data []map[string]interface{}) error {
	sqlStr := "INSERT into domain_time( domain ,time_dns,time_first_package,time_total,IP," +
		"http_code,size_download,NUM_REDIRECTS,speed_download,time_redirect,time_pretransfer," +
		"time_starttransfer) VALUES "
	vals := []interface{}{}
	for _, row := range data {
		sqlStr += "(?,?,?,?,?,?,?,?,?,?,?,?),"
		vals = append(vals, row["domain"], row["time_dns"], row["time_first_package"], row["time_total"],
			row["IP"], row["http_code"], row["size_download"], row["NUM_REDIRECTS"], row["speed_download"],
			row["time_redirect"], row["time_pretransfer"], row["time_starttransfer"])
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(vals...)
	if err != nil {
		return err
	}
	return nil
}

func query_time(domain string) ([]Data_For_Query, error) {
	unixtime := time.Now().Unix() - 360
	sqlStr := "SELECT time_dns,time_first_package,time_total,IP," +
		"http_code,size_download,NUM_REDIRECTS,speed_download,time_redirect,time_pretransfer," +
		"time_starttransfer FROM domain_time WHERE unix_timestamp(timestamp)>? AND " + `domain=?`
	stmt, _ := db.Prepare(sqlStr)
	rows, err := stmt.Query(unixtime, domain)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	data := Data_For_Query{}
	data_list := []Data_For_Query{}
	for rows.Next() {
		err := rows.Scan(&data.Time_dns, &data.Time_first_package, &data.Time_total, &data.IP, &data.Http_code,
			&data.Size_download, &data.NUM_REDIRECTS, &data.Speed_download, &data.Time_redirect,
			&data.Time_pretransfer, &data.Time_starttransfer)
		if err != nil {
			log.Printf("rows.scan err:%s", err)
			return nil, err
		}
		data_list = append(data_list, data)
	}
	defer rows.Close()
	return data_list, nil
}

func ipip_query(s string) (string, error) {
	db, err := ipdb.NewCity(global_config.Control.Ipip_file_path)
	if err != nil {
		return "", err
	}
	res, err := db.FindInfo(s, "CN")
	return res.RegionName, nil
}
func geoip_query(s string) (uint, error) {
	geodb, err := geoip2.Open(global_config.Control.Geoip_file_path)
	if err != nil {
		log.Println("geoip open error:", err)
		return 0, err
	}
	defer geodb.Close()
	city, _ := geodb.City(net.ParseIP(s))
	return city.City.GeoNameID, nil
}

func query_city(d []Data_For_Query) (map[string][]Data_For_Query, error) {
	tmp := map[string][]Data_For_Query{}
	for _, v := range d {
		cityid, _ := ipip_query(v.IP)
		tmp[cityid] = append(tmp[cityid], v)
	}
	return tmp, nil
}
func countScore(c map[string][]Data_For_Query) map[string]Signle_data {
	cityMax := Data_For_Query{}
	cityAve := Data_For_Query{}
	cityMin := Data_For_Query{}
	cityScoreTmp := Signle_data{}
	cityScore := map[string]Signle_data{}

	var count_http_code float64
	for k, v := range c {
		/*
		cityMaxVO:=reflect.ValueOf(&cityMax).Elem()
		cityMinVO:=reflect.ValueOf(&cityMin).Elem()
		for _,d:=range v {
			vo:=reflect.ValueOf(d)
			to:=reflect.TypeOf(d)
			for i:=0;i<vo.NumField();i++ {
				if vo.Field(i).Type()==reflect.Float32{
					a:=cityMaxVO.FieldByName(to.Field(i).Name)
					if vo.Field(i).Interface()>a {
						cityMaxVO.FieldByName(to.Field(i).Name)=vo.Field(i)
					}
				}
			}
		}*/
		//烂代码开始
		l := len(v)
		cityMax = Data_For_Query{}
		cityAve = Data_For_Query{}
		count_http_code = 0
		if l == 0 {
			continue
		} else {
			cityMin = v[0]
		}
		for _, d := range v {
			if d.Http_code == "200" {

				if cityMax.Time_first_package < d.Time_first_package {
					cityMax.Time_first_package = d.Time_first_package
				}
				if cityMin.Time_first_package > d.Time_first_package {
					cityMin.Time_first_package = d.Time_first_package
				}
				if cityMax.Time_dns < d.Time_dns {
					cityMax.Time_dns = d.Time_dns
				}
				if cityMin.Time_dns > d.Time_dns {
					cityMin.Time_dns = d.Time_dns
				}
				if cityMax.Time_total < d.Time_total {
					cityMax.Time_total = d.Time_total
				}
				if cityMin.Time_total > d.Time_total {
					cityMin.Time_total = d.Time_total
				}
				//平均和
				cityAve.Time_dns += d.Time_dns / float64(l)
				cityAve.Time_total += d.Time_total / float64(l)
				cityAve.Time_first_package += d.Time_first_package / float64(l)
			} else {
				//计数不成功的请求数
				count_http_code += 1
			}
		}
		//超过比例即为不能访问
		if count_http_code/float64(len(v)) > global_config.Control.HttpFailThreshold {
			cityScoreTmp.Status = false
		} else {
			//fmt.Println(cityMin, cityMax)
			scoreDns := (cityMax.Time_dns - cityAve.Time_dns) * 100 / (cityMax.Time_dns - cityMin.Time_dns)
			scoreFirst := (cityMax.Time_first_package - cityAve.Time_first_package) * 100 / (cityMax.Time_first_package - cityMin.Time_first_package)
			scoreTotal := (cityMax.Time_total - cityAve.Time_total) * 100 / (cityMax.Time_total - cityMin.Time_total)
			score := (scoreDns + scoreFirst + scoreTotal) * (1.0 / 3)
			cityScoreTmp.Status = true
			cityScoreTmp.Score = score
			cityScoreTmp.ScoreDns = scoreDns
			cityScoreTmp.ScoreTotal = scoreTotal
			cityScoreTmp.ScoreFirst = scoreFirst
			cityScoreTmp.Data = map[string]float64{
				"timeDns":   cityAve.Time_dns,
				"timeFirst": cityAve.Time_first_package,
				"timeTotal": cityAve.Time_total,
			}
			cityScore[k] = cityScoreTmp
		}
	}
	return cityScore
}

func count_country_score(p map[string]Signle_data) Signle_data {
	s := Signle_data{}
	l := float64(len(p) - 1)
	for k, v := range p {
		if k == "局域网" || v.Status == false {
			continue
		}
		s.ScoreDns += v.ScoreDns / l
		s.ScoreFirst += v.ScoreFirst / l
		s.ScoreTotal += v.ScoreTotal / l

	}
	s.Score = (s.ScoreTotal + s.ScoreFirst + s.ScoreDns) * (1.0 / 3)
	return s
}

func get_score(c *gin.Context) {
	if c.Query("domain") == "" {
		c.JSON(403, gin.H{"massage": false})
	}
	data, err := query_time(c.Query("domain"))
	if err != nil {
		c.JSON(403, gin.H{"massage": "query error"})
		log.Print(err)
	}
	if len(data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"massage": "domain data not found"})
	}
	tmp, err := query_city(data)
	if err != nil {
		c.JSON(503, gin.H{"massage": "query city error"})
	}
	score := countScore(tmp)
	country_score:=country_data{}
	country_score.Province=score
	country_score.CountryScore=count_country_score(score)
	c.JSON(200, country_score)

}
