package api

import "github.com/gin-gonic/gin"

func Register(e *gin.Engine) {
	e.POST("/", getUserRecord)
}
