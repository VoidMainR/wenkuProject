package network

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func StartApi() {
	g := gin.Default()
	g.Use(CORSMiddleware())

	g.GET("/getDoc", getDoc)

	s := &http.Server{
		Addr: ":8991",
		Handler: g,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	}
}
