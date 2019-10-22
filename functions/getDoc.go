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
	"time"
	"wenkuProject/conf"
)

func convertString(src string, srcCode string) string {
	srcDecoder := mahonia.NewDecoder(srcCode)
	convertedSrc := srcDecoder.ConvertString(src)
	return convertedSrc
}

func GetDoc(url string) (string,error) {
	htc := http.Client{}
	htc.Timeout = time.Second * 5
	rsp,err := htc.Get(url)

	if rsp != nil {
		defer rsp.Body.Close()
	}
	var buffer bytes.Buffer

	if err != nil {
		log.Printf("RequestErr:%v",err)
		return "",err
	}
	htmlBytes,err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("getHtmlErr:%v",err)
	}
	htmlStr := string(htmlBytes)

	// 获取title
	titleRE := regexp.MustCompile("<title>(.*?)</title>")
	titleStr := titleRE.FindStringSubmatch(htmlStr)[1]
	titleStr = strings.Replace(titleStr,"&nbsp;","",-1)
	titleStr = strings.Replace(titleStr," ","",-1)
	titleStr = convertString(titleStr,"gbk")


	re := regexp.MustCompile("(https.*?0.json.*?)x22}")
	urlArr := re.FindAllStringSubmatch(htmlStr,-1)
	newUrlArr := urlArr[:len(urlArr)/2]

	y := ""

	for _,textUrl := range newUrlArr {
		trueUrl := strings.ReplaceAll(textUrl[1],"\\","")
		textRsp,err := htc.Get(trueUrl)
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

	if buffer.Len() > 0 {
	  docPath,err := upLoadOSS(bytes.NewReader(buffer.Bytes()),titleStr)
	  return docPath, err
	} else {
		return "",err
	}

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