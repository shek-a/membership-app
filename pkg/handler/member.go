package handler

import (
	"github.com/gin-gonic/gin"
	"members.com/membership/pkg/models"
	"members.com/membership/pkg/service"
)

type MemberHandlerI interface {
	CreateMember(ctx *gin.Context)
}

type MemberHander struct {
	server        *gin.Engine
	memberService service.MemberServiceI
}

func NewMemberHandler(server *gin.Engine, memberService service.MemberServiceI) MemberHandlerI {
	return &MemberHander{
		server:        server,
		memberService: memberService,
	}
}

func (m *MemberHander) CreateMember(ctx *gin.Context) {
	var member models.Member
	if err := ctx.ShouldBindJSON(&member); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	response := m.memberService.CreateMember(ctx, &member)
	ctx.JSON(response.StatusCode, response.Body)
}
