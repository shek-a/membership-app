package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"members.com/membership/pkg/models"
	"members.com/membership/pkg/service"
)

type MemberHandlerI interface {
	CreateMember(ctx *gin.Context)
	GetMemberById(ctx *gin.Context)
	GetAllMembers(ctx *gin.Context)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response := m.memberService.CreateMember(ctx, &member)
	ctx.JSON(response.StatusCode, response.Body)
}

func (m *MemberHander) GetMemberById(ctx *gin.Context) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid member ID",
		})
		return
	}

	response := m.memberService.GetMemberById(ctx, int(memberId))
	ctx.JSON(response.StatusCode, response.Body)
}

func (m *MemberHander) GetAllMembers(ctx *gin.Context) {
	response := m.memberService.GetAllMembers(ctx)
	ctx.JSON(response.StatusCode, response.Body)
}
