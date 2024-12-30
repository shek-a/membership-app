package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"members.com/membership/pkg/models"
)

type MockMemberRepository struct {
	mock.Mock
}

func TestCreateMember(t *testing.T) {
	t.Parallel()

	member := &models.Member{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	testCases := []struct {
		name               string
		createMember       *models.Member
		memberRepoMock     func(ctx context.Context, mockRepo *MockMemberRepository)
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name:         "Success creating new member",
			createMember: member,
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("CreateMember", ctx, member).Return(nil)
			},
			expectedStatusCode: 201,
			expectedBody:       member,
		},
		{
			name:         "Error creating new member",
			createMember: member,
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("CreateMember", ctx, member).Return(errors.New("repository error"))
			},
			expectedStatusCode: 500,
			expectedBody:       models.ErrorMessage{Error: "Error creating member"},
		},
		{
			name: "Invalid email",
			createMember: &models.Member{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "John.Doegmail.com",
				DateOfBirth: "1990-01-01",
			},
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
			},
			expectedStatusCode: 400,
			expectedBody:       models.ErrorMessage{Error: "Invalid email"},
		},
		{
			name: "Invalid date of birth",
			createMember: &models.Member{
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "John.Doe@gmail.com",
				DateOfBirth: "1st April 1990",
			},
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
			},
			expectedStatusCode: 400,
			expectedBody:       models.ErrorMessage{Error: "Invalid date of birth"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockMemberRepository)
			tc.memberRepoMock(ctx, mockRepo)

			memberService := NewMemberService(mockRepo)
			response := memberService.CreateMember(ctx, tc.createMember)

			assert.Equal(t, tc.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tc.expectedBody, response.Body)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetMemberById(t *testing.T) {
	t.Parallel()

	memberId := 1
	member := &models.Member{
		ID:          memberId,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	testCases := []struct {
		name               string
		memberRepoMock     func(ctx context.Context, mockRepo *MockMemberRepository)
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "Success getting member by id",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(member, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       member,
		},
		{
			name: "Member is not found",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(nil, mongo.ErrNoDocuments)
			},
			expectedStatusCode: 404,
			expectedBody:       models.ErrorMessage{Error: "Member 1 not found"},
		},
		{
			name: "Error getting member by id",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(nil, errors.New("repository error"))
			},
			expectedStatusCode: 500,
			expectedBody:       models.ErrorMessage{Error: "Error fetching member"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockMemberRepository)
			tc.memberRepoMock(ctx, mockRepo)

			memberService := NewMemberService(mockRepo)
			response := memberService.GetMemberById(ctx, memberId)

			assert.Equal(t, tc.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tc.expectedBody, response.Body)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllMembers(t *testing.T) {
	t.Parallel()

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

	testCases := []struct {
		name               string
		memberRepoMock     func(ctx context.Context, mockRepo *MockMemberRepository)
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "Success getting all members",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetAllMembers", ctx).Return(members, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       members,
		},
		{
			name: "Error getting all members",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetAllMembers", ctx).Return(nil, errors.New("repository error"))
			},
			expectedStatusCode: 500,
			expectedBody:       models.ErrorMessage{Error: "Error fetching members"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockMemberRepository)
			tc.memberRepoMock(ctx, mockRepo)

			memberService := NewMemberService(mockRepo)
			response := memberService.GetAllMembers(ctx)

			assert.Equal(t, tc.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tc.expectedBody, response.Body)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateMemberById(t *testing.T) {
	t.Parallel()

	memberId := 1
	updateMember := &models.UpdateMember{
		FirstName: "Jonathan",
		Email:     "Jonathan.Doe@gmail.com",
	}
	member := &models.Member{
		ID:          memberId,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	testCases := []struct {
		name               string
		updateMember       *models.UpdateMember
		memberRepoMock     func(ctx context.Context, mockRepo *MockMemberRepository)
		expectedStatusCode int
		expectedBody       any
		wantErr            bool
	}{
		{
			name:         "Success updating member by id",
			updateMember: updateMember,
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(member, nil)
				mockRepo.On("UpdateMemberById", ctx, updateMember, memberId).Return(nil)
			},
			expectedStatusCode: 200,
			wantErr:            false,
		},
		{
			name:         "Member is not found",
			updateMember: updateMember,
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(nil, mongo.ErrNoDocuments)
			},
			expectedStatusCode: 404,
			expectedBody:       models.ErrorMessage{Error: "Member 1 not found"},
			wantErr:            true,
		},
		{
			name:         "Error getting member by id",
			updateMember: updateMember,
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(nil, errors.New("repository error"))
			},
			expectedStatusCode: 500,
			expectedBody:       models.ErrorMessage{Error: "Error fetching member"},
			wantErr:            true,
		},
		{
			name:         "Error updating member by id",
			updateMember: updateMember,
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(member, nil)
				mockRepo.On("UpdateMemberById", ctx, updateMember, memberId).Return(errors.New("repository error"))
			},
			expectedStatusCode: 500,
			expectedBody:       models.ErrorMessage{Error: "Error updating member"},
			wantErr:            true,
		},
		{
			name: "Invalid updated email",
			updateMember: &models.UpdateMember{
				Email: "Jonathan.Doegmail.com",
			},
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
			},
			expectedStatusCode: 400,
			expectedBody:       models.ErrorMessage{Error: "Invalid email"},
			wantErr:            true,
		},
		{
			name: "Invalid updated date of birth",
			updateMember: &models.UpdateMember{
				DateOfBirth: "1990-01-01T00:00:00Z",
			},
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
			},
			expectedStatusCode: 400,
			expectedBody:       models.ErrorMessage{Error: "Invalid date of birth"},
			wantErr:            true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockMemberRepository)
			tc.memberRepoMock(ctx, mockRepo)

			memberService := NewMemberService(mockRepo)
			response := memberService.UpdateMemberById(ctx, tc.updateMember, memberId)

			assert.Equal(t, tc.expectedStatusCode, response.StatusCode)
			if !tc.wantErr {
				updatedMember := response.Body.(*models.Member)
				assert.Equal(t, 1, updatedMember.ID)
				assert.Equal(t, "Jonathan", updatedMember.FirstName)
				assert.Equal(t, "Doe", updatedMember.LastName)
				assert.Equal(t, "Jonathan.Doe@gmail.com", updatedMember.Email)
				assert.Equal(t, "1990-01-01", updatedMember.DateOfBirth)
			} else {
				assert.Equal(t, tc.expectedBody, response.Body)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteMemberById(t *testing.T) {
	t.Parallel()

	memberId := 1
	member := &models.Member{
		ID:          memberId,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	testCases := []struct {
		name               string
		memberRepoMock     func(ctx context.Context, mockRepo *MockMemberRepository)
		expectedStatusCode int
		expectedBody       any
	}{
		{
			name: "Success deleting existing member",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(member, nil)
				mockRepo.On("DeleteMemberById", ctx, memberId).Return(nil)
			},
			expectedStatusCode: 204,
			expectedBody:       nil,
		},
		{
			name: "Member is not found",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(nil, mongo.ErrNoDocuments)
			},
			expectedStatusCode: 404,
			expectedBody:       models.ErrorMessage{Error: "Member 1 not found"},
		},
		{
			name: "Error deleting existing member",
			memberRepoMock: func(ctx context.Context, mockRepo *MockMemberRepository) {
				mockRepo.On("GetMemberById", ctx, memberId).Return(member, nil)
				mockRepo.On("DeleteMemberById", ctx, memberId).Return(errors.New("repository error"))
			},
			expectedStatusCode: 500,
			expectedBody:       models.ErrorMessage{Error: "Could not delete Member 1"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := new(MockMemberRepository)
			tc.memberRepoMock(ctx, mockRepo)

			memberService := NewMemberService(mockRepo)
			response := memberService.DeleteMemberById(ctx, memberId)

			assert.Equal(t, tc.expectedStatusCode, response.StatusCode)
			assert.Equal(t, tc.expectedBody, response.Body)
			mockRepo.AssertExpectations(t)
		})
	}
}

func (m *MockMemberRepository) CreateMember(ctx context.Context, member *models.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockMemberRepository) GetMemberById(ctx context.Context, memberId int) (*models.Member, error) {
	args := m.Called(ctx, memberId)
	member, ok := args.Get(0).(*models.Member)
	if !ok {
		return nil, args.Error(1)
	}
	return member, args.Error(1)
}

func (m *MockMemberRepository) GetAllMembers(ctx context.Context) ([]models.Member, error) {
	args := m.Called(ctx)
	members, ok := args.Get(0).([]models.Member)
	if !ok {
		return nil, args.Error(1)
	}
	return members, args.Error(1)
}

func (m *MockMemberRepository) UpdateMemberById(ctx context.Context, member *models.UpdateMember, memberId int) error {
	args := m.Called(ctx, member, memberId)
	return args.Error(0)
}

func (m *MockMemberRepository) DeleteMemberById(ctx context.Context, memberId int) error {
	args := m.Called(ctx, memberId)
	return args.Error(0)
}
