package network

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wenkuProject/conf"
	"wenkuProject/functions"
)

func getDoc(c *gin.Context) {
	docUrl := c.Query("url")
	docPath,e := functions.GetDoc(docUrl)
	if e != nil {
		c.Error(e)
	}
	if docPath == "" {
		c.Error(e)
		return
	}
	docDownloadUrl := conf.OssConf.DownloadUrl+docPath
	c.SecureJSON(http.StatusOK,docDownloadUrl)
	return
}