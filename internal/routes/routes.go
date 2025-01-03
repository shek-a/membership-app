package routes

import (
	"github.com/gin-gonic/gin"
	"members.com/membership/pkg/handler"
)

func RegisterRoutes(server *gin.Engine, handler handler.MemberHandlerI) {
	server.POST("/member", handler.CreateMember)
	server.GET("/member/:id", handler.GetMemberById)
	server.GET("/members", handler.GetAllMembers)
	server.PUT("/member/:id", handler.UpdateMemberById)
	server.DELETE("/member/:id", handler.DeleteMemberById)
}
