package network

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func StartApi() {
	g := gin.Default()

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

