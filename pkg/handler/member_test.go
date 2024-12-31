package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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

	member := &models.Member{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	mockService := new(MockMemberService)
	mockService.On("CreateMember", mock.Anything, mock.Anything).Return(createResponse(http.StatusOK, member))

	memberHandler := NewMemberHandler(router, mockService)
	router.POST("/member", memberHandler.CreateMember)

	testCases := []struct {
		name                 string
		requestBody          string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "Success creating new member",
			requestBody:          `{"firstName": "John", "lastName": "Doe", "email": "John.Doe@gmail.com", "dateOfBirth": "1990-01-01"}`,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1,\"firstName\":\"John\",\"lastName\":\"Doe\",\"email\":\"John.Doe@gmail.com\",\"dateOfBirth\":\"1990-01-01\"}",
		},
		{
			name:                 "Invalid request",
			requestBody:          `{"firstName": "John", "lastName": "Doe", "dateOfBirth": "1990-01-01"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"error\":\"Invalid request\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/member", bytes.NewBufferString(tc.requestBody))
			request.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedResponseBody)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetMemberById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	member := &models.Member{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	mockService := new(MockMemberService)

	memberHandler := NewMemberHandler(router, mockService)
	router.GET("/member/:id", memberHandler.GetMemberById)

	testCases := []struct {
		name                 string
		memberId             string
		mockMemberService    func(mockService *MockMemberService)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Success getting member by id",
			memberId: "1",
			mockMemberService: func(mockService *MockMemberService) {
				mockService.On("GetMemberById", mock.Anything, 1).Return(createResponse(http.StatusOK, member))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1,\"firstName\":\"John\",\"lastName\":\"Doe\",\"email\":\"John.Doe@gmail.com\",\"dateOfBirth\":\"1990-01-01\"}",
		},
		{
			name:     "Invalid member ID",
			memberId: "1x",
			mockMemberService: func(mockService *MockMemberService) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"error\":\"Invalid member ID\"}",
		},
		{
			name:     "Member not found",
			memberId: "1",
			mockMemberService: func(mockService *MockMemberService) {
				mockService.On("GetMemberById", mock.Anything, 1).Return(createResponse(http.StatusNotFound, models.ErrorMessage{Error: "Member 1 not found"}))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"error\":\"Member 1 not found\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockMemberService(mockService)
			request, _ := http.NewRequest(http.MethodGet, "/member/"+tc.memberId, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedResponseBody)
			mockService.AssertExpectations(t)
			mockService.ExpectedCalls = nil
		})
	}
}

func TestGetAllMembers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	members := []models.Member{
		{
			ID:          1,
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "John.Doe@gmail.com",
			DateOfBirth: "1990-01-01",
		},
		{
			ID:          2,
			FirstName:   "Jane",
			LastName:    "Smith",
			Email:       "Jane.Smith@gmail.com",
			DateOfBirth: "1985-05-05",
		},
	}

	mockService := new(MockMemberService)
	mockService.On("GetAllMembers", mock.Anything).Return(createResponse(http.StatusOK, members))

	memberHandler := NewMemberHandler(router, mockService)
	router.GET("/members", memberHandler.GetAllMembers)

	request, _ := http.NewRequest(http.MethodGet, "/members", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "[{\"id\":1,\"firstName\":\"John\",\"lastName\":\"Doe\",\"email\":\"John.Doe@gmail.com\",\"dateOfBirth\":\"1990-01-01\"},{\"id\":2,\"firstName\":\"Jane\",\"lastName\":\"Smith\",\"email\":\"Jane.Smith@gmail.com\",\"dateOfBirth\":\"1985-05-05\"}]")
	mockService.AssertExpectations(t)
}

func TestUpdateMemberById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	member := &models.Member{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	mockService := new(MockMemberService)

	memberHandler := NewMemberHandler(router, mockService)
	router.PUT("/member/:id", memberHandler.UpdateMemberById)

	testCases := []struct {
		name                 string
		memberId             string
		requestBody          string
		mockMemberService    func(mockService *MockMemberService)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Success updating member by id",
			memberId:    "1",
			requestBody: `{"email": "John.Doe@gmail.com"}`,
			mockMemberService: func(mockService *MockMemberService) {
				mockService.On("UpdateMemberById", mock.Anything, mock.Anything, 1).Return(createResponse(http.StatusOK, member))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"id\":1,\"firstName\":\"John\",\"lastName\":\"Doe\",\"email\":\"John.Doe@gmail.com\",\"dateOfBirth\":\"1990-01-01\"}",
		},
		{
			name:        "Invalid request where member ID is invalid",
			memberId:    "1x",
			requestBody: `{"email": "John.Doe@gmail.com"}`,
			mockMemberService: func(mockService *MockMemberService) {

			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"error\":\"Invalid member ID\"}",
		},
		{
			name:        "Invalid request where email is invalid",
			memberId:    "1",
			requestBody: `{"email": 123}`,
			mockMemberService: func(mockService *MockMemberService) {

			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"error\":\"Invalid request\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockMemberService(mockService)
			request, _ := http.NewRequest(http.MethodPut, "/member/"+tc.memberId, bytes.NewBufferString(tc.requestBody))
			request.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedResponseBody)
			mockService.AssertExpectations(t)
			mockService.ExpectedCalls = nil
		})
	}
}

func TestDeleteMemberById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockService := new(MockMemberService)

	memberHandler := NewMemberHandler(router, mockService)
	router.DELETE("/member/:id", memberHandler.DeleteMemberById)

	testCases := []struct {
		name                 string
		memberId             string
		mockMemberService    func(mockService *MockMemberService)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "Success deleting existing member by id",
			memberId: "1",
			mockMemberService: func(mockService *MockMemberService) {
				mockService.On("DeleteMemberById", mock.Anything, 1).Return(createResponse(http.StatusOK, models.SuccessMessage{Message: "Member 1 deleted"}))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"message\":\"Member 1 deleted\"}",
		},
		{
			name:     "Invalid member ID",
			memberId: "1x",
			mockMemberService: func(mockService *MockMemberService) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "{\"error\":\"Invalid member ID\"}",
		},
		{
			name:     "Member not found",
			memberId: "1",
			mockMemberService: func(mockService *MockMemberService) {
				mockService.On("DeleteMemberById", mock.Anything, 1).Return(createResponse(http.StatusNotFound, models.ErrorMessage{Error: "Member 1 not found"}))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"error\":\"Member 1 not found\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockMemberService(mockService)
			request, _ := http.NewRequest(http.MethodDelete, "/member/"+tc.memberId, nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedResponseBody)
			mockService.AssertExpectations(t)
			mockService.ExpectedCalls = nil
		})
	}
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

func (m *MockMemberService) GetMemberById(ctx context.Context, memberId int) models.Response {
	args := m.Called(ctx, memberId)
	return args.Get(0).(models.Response)
}

func (m *MockMemberService) GetAllMembers(ctx context.Context) models.Response {
	args := m.Called(ctx)
	return args.Get(0).(models.Response)
}

func (m *MockMemberService) UpdateMemberById(ctx context.Context, member *models.UpdateMember, memberId int) models.Response {
	args := m.Called(ctx, member, memberId)
	return args.Get(0).(models.Response)
}

func (m *MockMemberService) DeleteMemberById(ctx context.Context, memberId int) models.Response {
	args := m.Called(ctx, memberId)
	return args.Get(0).(models.Response)
}
