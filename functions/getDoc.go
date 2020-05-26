package functions

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"wenkuProject/conf"
)

func convertString(src string, srcCode string) string {
	srcDecoder := mahonia.NewDecoder(srcCode)
	convertedSrc := srcDecoder.ConvertString(src)
	return convertedSrc
}

// 获取百度文库文档脚本
func GetDoc(url string) (string,error) {
	htc := &http.Client{}
	//htc.Timeout = time.Second * 5
	req,_ := http.NewRequest("GET",url,nil)

	// 请求头添加
	req.Header.Set("User-Agent", "uuid")
	//req.Header.Set("Accept", "*/*")
	//req.Header.Set("Cache-Control", "no-cache")
	//req.Header.Set("Postman-Token", "b8b71aa9-7f2e-49b3-9516-5d6ce6c9eb39")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//req.Header.Set("Cookie", "BAIDUID=7F9269C5CCCE1F789B46CF429F9CD8EF:FG=1")
	//req.Header.Set("Connection", "keep-alive")

	rsp,err := htc.Do(req)

	if rsp != nil {
		defer rsp.Body.Close()
	}


	if err != nil {
		log.Printf("RequestErr:%v",err)
		return "",err
	}
	htmlBytes,err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("getHtmlErr:%v",err)
	}
	htmlStr := string(htmlBytes)
	
	if htmlStr == "" {
		return "", nil
	}

	buffer, titleStr := handleBody(htmlStr)

	if buffer.Len() > 0 {
		docPath,err := upLoadOSS(bytes.NewReader(buffer.Bytes()),titleStr)
		return docPath, err
	} else {
		return "",err
	}
}

func handleBody(htmlStr string) (*bytes.Buffer, string)  {
	var buffer bytes.Buffer

	// 获取title
	titleRE := regexp.MustCompile("<title>(.*?)</title>")
	titleStr := titleRE.FindStringSubmatch(htmlStr)[1]
	titleStr = strings.Replace(titleStr,"&nbsp;","",-1)
	titleStr = strings.Replace(titleStr," ","",-1)
	titleStr = convertString(titleStr,"gbk")


	re := regexp.MustCompile("(https.*?0.json.*?)x22}")
	urlArr := re.FindAllStringSubmatch(htmlStr,-1)
	if urlArr == nil {
		return &buffer, ""
	}
	newUrlArr := urlArr[:len(urlArr)/2]

	y := ""

	for _,textUrl := range newUrlArr {
		trueUrl := strings.ReplaceAll(textUrl[1],"\\","")
		textRsp,err := http.Get(trueUrl)
		if err != nil {
			log.Printf("getUrlErr:%v",err)
		}

		textBytes, err := ioutil.ReadAll(textRsp.Body)
		textStr := string(textBytes)
		tre := regexp.MustCompile("\"c\":\"(.*?)\".*?\"y\":(.*?),")
		texts:= tre.FindAllStringSubmatch(textStr,-1)
		for _,text := range texts {
			n := ""
			zhString,err:= strconv.Unquote(`"` + text[1] + `"`)

			if err != nil {
				log.Printf("ee:%v,%v",err,text[1])
				continue
			}

			if text[2] != y {
				y = text[2]
				n = "\n"
			} else {
				n = ""
			}

			_,err = buffer.WriteString(n + zhString)
			if err != nil {
				log.Printf("writeBufferErr:%v",err)
			}
		}
		textRsp.Body.Close()
	}

	return &buffer, titleStr
}


// 上传文件到oss
func upLoadOSS(strReader io.Reader, docName string) (string,error) {
	ossClient,err := oss.New(conf.OssConf.EndPoint, conf.OssConf.AccessKeyId, conf.OssConf.AccessKeySecret)
	if err != nil {
		log.Printf("ossErr:%v",err)
		return "",err
	}

	bucket,err := ossClient.Bucket(conf.OssConf.Bucket)
	if err != nil {
		log.Printf("bucket Err:%v",err)
		return "",err
	}

	docPath := fmt.Sprintf("tools/outputs/%v.doc", docName)

	err = bucket.PutObject(docPath,strReader)
	if err != nil {
		log.Printf("save doc err:%v",err)
	}
	return docPath,nil
}