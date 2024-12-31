package service

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"members.com/membership/pkg/models"
	"members.com/membership/pkg/repository"
	"members.com/membership/pkg/utils"
)

type MemberServiceI interface {
	CreateMember(ctx context.Context, member *models.Member) models.Response
	GetMemberById(ctx context.Context, memberId int) models.Response
	GetAllMembers(ctx context.Context) models.Response
	UpdateMemberById(ctx context.Context, member *models.UpdateMember, memberId int) models.Response
	DeleteMemberById(ctx context.Context, memberId int) models.Response
}

type MemberService struct {
	memberRepository repository.MemberRepositoryI
}

func NewMemberService(memberRepository repository.MemberRepositoryI) MemberServiceI {
	return &MemberService{
		memberRepository: memberRepository,
	}
}

func (m *MemberService) CreateMember(ctx context.Context, member *models.Member) models.Response {
	member.ID = utils.GenerateRandomNumber()

	if !utils.IsValidEmail(member.Email) {
		return createErrorResponse(http.StatusBadRequest, "Invalid email")
	}

	if !utils.IsValidDate(member.DateOfBirth) {
		return createErrorResponse(http.StatusBadRequest, "Invalid date of birth")
	}

	err := m.memberRepository.CreateMember(ctx, member)
	if err != nil {
		return createErrorResponse(http.StatusInternalServerError, "Error creating member")
	}
	return models.Response{
		StatusCode: http.StatusCreated,
		Body:       member,
	}
}

func (m *MemberService) GetMemberById(ctx context.Context, memberId int) models.Response {
	member, err := m.memberRepository.GetMemberById(ctx, memberId)
	if err != nil {
		return handleMemberFetchError(err, memberId)
	}
	return models.Response{
		StatusCode: http.StatusOK,
		Body:       member,
	}
}

func (m *MemberService) GetAllMembers(ctx context.Context) models.Response {
	members, err := m.memberRepository.GetAllMembers(ctx)
	if err != nil {
		return createErrorResponse(http.StatusInternalServerError, "Error fetching members")
	}
	return models.Response{
		StatusCode: http.StatusOK,
		Body:       members,
	}
}

func (m *MemberService) UpdateMemberById(ctx context.Context, member *models.UpdateMember, memberId int) models.Response {
	if member.Email != "" && !utils.IsValidEmail(member.Email) {
		return createErrorResponse(http.StatusBadRequest, "Invalid email")
	}

	if member.DateOfBirth != "" && !utils.IsValidDate(member.DateOfBirth) {
		return createErrorResponse(http.StatusBadRequest, "Invalid date of birth")
	}

	fetchedMember, err := m.memberRepository.GetMemberById(ctx, memberId)
	if err != nil {
		return handleMemberFetchError(err, memberId)
	}

	fetchedMember = mergeFields(fetchedMember, member)

	err = m.memberRepository.UpdateMemberById(ctx, member, memberId)
	if err != nil {
		return createErrorResponse(http.StatusInternalServerError, "Error updating member")
	}
	return models.Response{
		StatusCode: http.StatusOK,
		Body:       fetchedMember,
	}
}

func (m *MemberService) DeleteMemberById(ctx context.Context, memberId int) models.Response {
	_, err := m.memberRepository.GetMemberById(ctx, memberId)
	if err != nil {
		return handleMemberFetchError(err, memberId)
	}

	err = m.memberRepository.DeleteMemberById(ctx, memberId)
	if err != nil {
		return createErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Could not delete Member %d", memberId))
	}
	return createSuccessResponse(http.StatusOK, fmt.Sprintf("Member %d deleted", memberId))
}

func mergeFields(member *models.Member, updateMember *models.UpdateMember) *models.Member {
	if updateMember.FirstName != "" {
		member.FirstName = updateMember.FirstName
	}

	if updateMember.LastName != "" {
		member.LastName = updateMember.LastName
	}

	if updateMember.Email != "" {
		member.Email = updateMember.Email
	}

	if updateMember.DateOfBirth != "" {
		member.DateOfBirth = updateMember.DateOfBirth
	}
	return member
}

func handleMemberFetchError(err error, memberId int) models.Response {
	if err == mongo.ErrNoDocuments {
		return createErrorResponse(http.StatusNotFound, fmt.Sprintf("Member %d not found", memberId))
	}
	return createErrorResponse(http.StatusInternalServerError, "Error fetching member")
}

func createSuccessResponse(statusCode int, successMessage string) models.Response {
	return models.Response{
		StatusCode: statusCode,
		Body: models.SuccessMessage{
			Message: successMessage,
		},
	}
}

func createErrorResponse(statusCode int, errorMessage string) models.Response {
	return models.Response{
		StatusCode: statusCode,
		Body: models.ErrorMessage{
			Error: errorMessage,
		},
	}
}
