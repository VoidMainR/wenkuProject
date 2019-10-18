package functions

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetDoc(url string) (string,error) {
	// 这里过滤一下空格
	htc := http.Client{}
	htc.Timeout = time.Second * 5
	rsp,err := htc.Get(url)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		log.Printf("访问网页出错,%v",err)
		return "",err
	}
	htmlBytes,err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("获取html出错%v",err)
	}
	htmlStr := string(htmlBytes)
	re := regexp.MustCompile("(https.*?0.json.*?)x22}")
	urlArr := re.FindAllStringSubmatch(htmlStr,-1)
	newUrlArr := urlArr[:len(urlArr)/2]

	fileName := ""
	var docFile *os.File = nil

	y := ""

	for i,textUrl := range newUrlArr {
		trueUrl := strings.ReplaceAll(textUrl[1],"\\","")
		textRsp,err := htc.Get(trueUrl)
		if err != nil {
			log.Printf("获取文章出错,%v",err)
		}

		textBytes, err := ioutil.ReadAll(textRsp.Body)
		textStr := string(textBytes)
		tre := regexp.MustCompile("\"c\":\"(.*?)\".*?\"y\":(.*?),")
		texts:= tre.FindAllStringSubmatch(textStr,-1)
		for k,text := range texts {
			n := ""
			zhString,err:= strconv.Unquote(`"` + text[1] + `"`)
			if text[2] != y {
				y = text[2]
				n = "\n"
			} else {
				n = ""
			}
			if err != nil {
				log.Printf("ee:%v,%v",err,text[1])
				continue
			}

			// 这里赋值一个文件名
			if i == 0 && k == 0 {
				fileName = zhString + ".doc"
				docFile,err = os.Create(fileName)
				if err != nil {
					log.Printf("创建文件失败:%v",err)
					return "",err
				}
			}
			if docFile == nil {
				break
			}

			_,err = docFile.WriteString(n + zhString)
			if err != nil {
				log.Printf("写入文件失败:%v",err)
			}
		}
		textRsp.Body.Close()
	}

	if docFile != nil {
		docFile.Close()
	}

	return fileName,err
}