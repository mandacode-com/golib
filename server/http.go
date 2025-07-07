package server

import "github.com/gin-gonic/gin"

type HTTPRegisterer interface {
	Register(group *gin.RouterGroup) error
}
