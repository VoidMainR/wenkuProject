package network

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"wenkuProject/functions"
)

func getDoc(c *gin.Context) {
	docUrl := c.Query("url")
	//log.Printf("uuuurl:%s",docUrl)
	docName,e := functions.GetDoc(docUrl)
	if e != nil {
		c.Error(e)
	}
	if docName == "" {
		c.Error(e)
		return
	}
	re := regexp.MustCompile("/(.*)$")
	name := re.FindStringSubmatch(docName)
	c.FileAttachment(docName,name[1])
	c.Set("data","success")
	return
}