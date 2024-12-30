package routes

import (
	"github.com/gin-gonic/gin"
	"members.com/membership/pkg/handler"
)

func RegisterRoutes(server *gin.Engine, handler handler.MemberHandlerI) {
	server.POST("/member", handler.CreateMember)
}
