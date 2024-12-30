package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"members.com/membership/pkg/models"
)

type MockMemberService struct {
	mock.Mock
}

func TestCreateMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctx := context.Background()

	member := &models.Member{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: models.Date{Time: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)},
	}

	mockService := new(MockMemberService)
	mockService.On("CreateMember", ctx, mock.Anything).Return(createResponse(http.StatusOK, member))

	memberHandler := NewMemberHandler(router, mockService)
	router.POST("/member", memberHandler.CreateMember)

	t.Run("valid request", func(t *testing.T) {
		body := `{"firstName": "John", "lastName": "Doe", "dateOfBirth": "1990-01-01"}`
		req, _ := http.NewRequest(http.MethodPost, "/member", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "member created")
	})

	t.Run("invalid request", func(t *testing.T) {
		body := `{"name": "John Doe"}`
		req, _ := http.NewRequest(http.MethodPost, "/member", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request")
	})
}

func createResponse(statusCode int, body any) models.Response {
	return models.Response{
		StatusCode: statusCode,
		Body:       body,
	}
}

func (m *MockMemberService) CreateMember(ctx context.Context, member *models.Member) models.Response {
	args := m.Called(ctx, member)
	return args.Get(0).(models.Response)
}

// DeleteMemberById implements service.MemberServiceI.
func (m *MockMemberService) DeleteMemberById(ctx context.Context, memberId int) models.Response {
	panic("unimplemented")
}

// GetAllMembers implements service.MemberServiceI.
func (m *MockMemberService) GetAllMembers(ctx context.Context) models.Response {
	panic("unimplemented")
}

// GetMemberById implements service.MemberServiceI.
func (m *MockMemberService) GetMemberById(ctx context.Context, memberId int) models.Response {
	panic("unimplemented")
}

// UpdateMemberById implements service.MemberServiceI.
func (m *MockMemberService) UpdateMemberById(ctx context.Context, member *models.Member, memberId int) models.Response {
	panic("unimplemented")
}
