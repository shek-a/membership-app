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
	UpdateMemberById(ctx *gin.Context)
	DeleteMemberById(ctx *gin.Context)
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
	var newMember models.Member
	if !bindJsonBody(ctx, &newMember) {
		return
	}

	response := m.memberService.CreateMember(ctx, &newMember)
	ctx.JSON(response.StatusCode, response.Body)
}

func (m *MemberHander) GetMemberById(ctx *gin.Context) {
	memberId, valid := extractMemberIdfromUrlPath(ctx)
	if !valid {
		return
	}

	response := m.memberService.GetMemberById(ctx, int(memberId))
	ctx.JSON(response.StatusCode, response.Body)
}

func (m *MemberHander) GetAllMembers(ctx *gin.Context) {
	response := m.memberService.GetAllMembers(ctx)
	ctx.JSON(response.StatusCode, response.Body)
}

func (m *MemberHander) UpdateMemberById(ctx *gin.Context) {
	var updateMember models.UpdateMember
	if !bindJsonBody(ctx, &updateMember) {
		return
	}

	memberId, valid := extractMemberIdfromUrlPath(ctx)
	if !valid {
		return
	}

	response := m.memberService.UpdateMemberById(ctx, &updateMember, int(memberId))
	ctx.JSON(response.StatusCode, response.Body)
}

func (m *MemberHander) DeleteMemberById(ctx *gin.Context) {
	memberId, valid := extractMemberIdfromUrlPath(ctx)
	if !valid {
		return
	}

	response := m.memberService.DeleteMemberById(ctx, int(memberId))
	ctx.JSON(response.StatusCode, response.Body)
}

func bindJsonBody(ctx *gin.Context, obj interface{}) bool {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return false
	}
	return true
}

func extractMemberIdfromUrlPath(ctx *gin.Context) (int64, bool) {
	memberId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid member ID",
		})
		return 0, false
	}
	return memberId, true
}
