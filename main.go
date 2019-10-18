package main

import (
	"log"
	"wenkuProject/functions"
)

func main() {
	// 文库url
	url := "https://wenku.baidu.com/view/4ac57d39fc0a79563c1ec5da50e2524de518d0d8.html"
	name,err := functions.GetDoc(url)
	if err != nil {
		log.Printf("error:%v \n",err)
	} else {
		log.Printf("文件已保存为:%v",name)
	}

}