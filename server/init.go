package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
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
	Geoip_file_path   string  `yaml:"geoip_file_path"`
	Ipip_file_path    string  `yaml:"ipip_file_path"`
}

func (config *Conf) getconf() (*Conf) {
	yamlfile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.UnmarshalStrict(yamlfile, config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

/* Global Variables */
var global_config Conf

func initall() {
	global_config = *global_config.getconf()
	//fmt.Println(global_config)
	var err error
	db, err = initdb()
	if err != nil {
		log.Fatal(err)
	}

}
