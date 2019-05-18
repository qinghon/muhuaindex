package main

import (
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Conf struct {
	DATABASE_USER   string  `yaml:"DATABASE_USER"`
	DATABASE_PASSWD string  `yaml:"DATABASE_PASSWD"`
	DATABASE_NAME   string  `yaml:"DATABASE_NAME"`
	DATABASE_PORT   string  `yaml:"DATABASE_PORT"`
	DATABASE_HOST   string  `yaml:"DATABASE_HOST"`
	Control         control `yaml:"control"`
}
type control struct {
	HttpFailThreshold float64 `yaml:"http_fail_threshold"`
	//Geoip_file_path   string  `yaml:"geoip_file_path"`
	Ipip_file_path string `yaml:"ipip_file_path"`
}

func (config *Conf) check() bool {
	return checktype(*config)
}
func checktype(i interface{}) bool {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		switch t.Field(i).Type.String() {
		case "string":
			if v.Field(i).String() == "" {
				return false
			}
		case "int":
			if v.Field(i).Int() == 0 {
				return false
			}
		case "float64":
			if v.Field(i).Float() == 0 {
				return false
			}
		case "main.control":
			//fmt.Println(v.Field(i))
			return checktype(v.Field(i).Interface())
		default:
			return true
		}
	}
	return true
}
func (config *Conf) getconf() (Conf, error) {
	//读取config.yaml文件
	yamlfile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Println(err)
		config = nil
		return Conf{}, err
	}
	err = yaml.UnmarshalStrict(yamlfile, config)
	if err != nil {
		log.Println(err)
		return Conf{}, err
	}
	return *config, nil
}

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
func (config *Conf) getconf_env() (Conf, error) {

	//读取环境变量
	config.DATABASE_USER = os.Getenv("DATABASE_USER")
	config.DATABASE_PASSWD = os.Getenv("DATABASE_PASSWD")
	config.DATABASE_NAME = os.Getenv("DATABASE_NAME")
	config.DATABASE_PORT = os.Getenv("DATABASE_PORT")
	config.DATABASE_HOST = os.Getenv("DATABASE_HOST")
	if config.DATABASE_NAME == "" || config.DATABASE_USER == "" ||
		config.DATABASE_PASSWD == "" || config.DATABASE_PORT == "" ||
		config.DATABASE_HOST == "" {
		//log.Println("DATABSE config not set ,exit")

		return *config, errors.New("DATABSE env config not set ")
	}
	config.Control.HttpFailThreshold = 0.74
	http_fail_threshold := os.Getenv("http_fail_threshold")
	if http_fail_threshold == "" {
		log.Println("http_fail_threshold not set,defult set 0.74")
	} else if tmp, err := strconv.ParseFloat(http_fail_threshold,64); err != nil {
		log.Println("Conversion http_fail_threshold err,set defult 0.74", err)
	} else {
		config.Control.HttpFailThreshold = tmp
	}
	config.Control.Ipip_file_path = os.Getenv("ipip_file_path")
	if config.Control.Ipip_file_path == "" {
		log.Println("env \"ipip_file_path\" not foud ,reading path ./ipipfree.ipdb")
		if PathExist("./ipipfree.ipdb") {
			config.Control.Ipip_file_path = "./ipipfree.ipdb"
		}else{
			log.Println("./ipipfree.ipdb file not found ")
			return Conf{}, errors.New("ipip_file_path not foud not set ,exit")
		}
	}
	return *config, nil
}


func initdb(config Conf)  {
	//初始化数据库
	if PathExist("init.lock") {
		log.Println("init.lock exists ")
		return
	}else {
		f,err:=os.Create("init.lock")
		if err!=nil {
			log.Println(err)
		}
		defer f.Close()
	}
	var err error
	par := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&timeout=30s",
		config.DATABASE_USER,
		config.DATABASE_PASSWD,
		config.DATABASE_HOST,
		config.DATABASE_PORT,
		config.DATABASE_NAME)

	db, err = sql.Open("mysql", par)
	if err != nil {
		log.Println("Failed to connect to log mysql: ", err)
	}
	file,err:=ioutil.ReadFile("init.sql")
	if err!=nil {
		log.Println(err)
	}
	resquest:=strings.Split(string(file),";\n")
	for _,req:=range resquest {
		_,err:=db.Exec(req)
		if err!=nil {
			log.Println("Error run:",err)
		}
	}
	log.Println("Create table over!")
}
/* Global Variables */
var global_config Conf

func initall() {
	var err error
	if global_config, err = global_config.getconf_env(); err != nil {
		log.Println(err)
		log.Println("reading config.yaml file ...")
		err = nil
		global_config, err = global_config.getconf()
		//fmt.Println(global_config)
		if err != nil {
			log.Fatalln(err)
		}
	}
	initdb(global_config)
	db, err = connectdb()

	if err != nil {
		log.Fatal(err)
	}

}
