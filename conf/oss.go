package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ossConf struct {
	EndPoint string 	`json:"endPoint"`
	AccessKeyId string	`json:"accessKeyId"`
	AccessKeySecret string	`json:"accessKeySecret"`
	Bucket string			`json:"bucket"`
	DownloadUrl string		`json:"downloadUrl"`
}

var OssConf ossConf

func init() {
	bytes,err := ioutil.ReadFile("conf/oss.json")
	if err != nil {
		log.Printf("读取配置出错,%v",err)
		return
	}

	err = json.Unmarshal(bytes,&OssConf)
	if err != nil {
		log.Printf("jsonMarshal出错：%v",err)
	}
}