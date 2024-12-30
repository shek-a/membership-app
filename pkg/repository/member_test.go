package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"members.com/membership/pkg/models"
)

func TestCreateMember(t *testing.T) {
	t.Parallel()

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("members").ClientType(mtest.Mock))

	testCases := []struct {
		name        string
		mongoDbMock func(mt *mtest.T)
		wantErr     bool
	}{
		{
			name: "success creating new member",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: false,
		},
		{
			name: "error creating new member",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   0,
					Code:    11000,
					Message: "duplicate key error",
				}))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mt *mtest.T) {
			tc.mongoDbMock(mt)
			repo := NewMembershipRepository(mt.DB)
			member := &models.Member{
				ID:          1,
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "John.Doe@gmail.com",
				DateOfBirth: "1990-01-01",
			}
			err := repo.CreateMember(context.Background(), member)

			if tc.wantErr {
				assert.Errorf(t, err, "Want error but got: %v", err)
			} else {
				assert.NoErrorf(t, err, "Not expecting error")
			}
		})
	}
}

func TestGetMemberById(t *testing.T) {
	t.Parallel()

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("members").ClientType(mtest.Mock))

	member := &models.Member{
		ID:          1,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	testCases := []struct {
		name        string
		mongoDbMock func(mt *mtest.T)
		wantErr     bool
	}{
		{
			name: "success getting member by id",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "membership.members", mtest.FirstBatch, bson.D{
					{Key: "ID", Value: member.ID},
					{Key: "FirstName", Value: member.FirstName},
					{Key: "LastName", Value: member.LastName},
					{Key: "Email", Value: member.Email},
					{Key: "DateOfBirth", Value: member.DateOfBirth},
				}))
			},
			wantErr: false,
		},
		{
			name: "member by id not found",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(0, "membership.members", mtest.FirstBatch))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mt *mtest.T) {
			tc.mongoDbMock(mt)
			repo := NewMembershipRepository(mt.DB)
			member, err := repo.GetMemberById(context.Background(), member.ID)

			if tc.wantErr {
				assert.Errorf(t, err, "Want error but got: %v", err)
			} else {
				assert.NoErrorf(t, err, "Not expecting error")
				assert.Equal(t, 1, member.ID)
				assert.Equal(t, "John", member.FirstName)
				assert.Equal(t, "Doe", member.LastName)
				assert.Equal(t, "John.Doe@gmail.com", member.Email)
				assert.Equal(t, "1990-01-01", member.DateOfBirth)
			}
		})
	}
}

func TestGetAllMembers(t *testing.T) {
	t.Parallel()

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("members").ClientType(mtest.Mock))

	testCases := []struct {
		name        string
		mongoDbMock func(mt *mtest.T)
		wantErr     bool
	}{
		{
			name: "success getting all members",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(2, "membership.members", mtest.FirstBatch, bson.D{
					{Key: "ID", Value: 1},
					{Key: "FirstName", Value: "John"},
					{Key: "LastName", Value: "Doe"},
					{Key: "Email", Value: "John.Doe@gmail.com"},
					{Key: "DateOfBirth", Value: "1990-01-01"},
				}, bson.D{
					{Key: "ID", Value: 2},
					{Key: "FirstName", Value: "Jane"},
					{Key: "LastName", Value: "Smith"},
					{Key: "Email", Value: "Jane.Smith@gmail.com"},
					{Key: "DateOfBirth", Value: "1985-05-05"},
				}))
			},
			wantErr: false,
		},
		{
			name: "error getting all members",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    11000,
					Message: "fetching members failed",
				}))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mt *mtest.T) {
			tc.mongoDbMock(mt)
			repo := NewMembershipRepository(mt.DB)
			members, err := repo.GetAllMembers(context.Background())

			if tc.wantErr {
				assert.Errorf(t, err, "Want error but got: %v", err)
				assert.Len(t, members, 0)
			} else {
				assert.NoErrorf(t, err, "Not expecting error")
				assert.Len(t, members, 2)

				member1 := members[0]
				assert.Equal(t, 1, member1.ID)
				assert.Equal(t, "John", member1.FirstName)
				assert.Equal(t, "Doe", member1.LastName)
				assert.Equal(t, "John.Doe@gmail.com", member1.Email)
				assert.Equal(t, "1990-01-01", member1.DateOfBirth)

				member2 := members[1]
				assert.Equal(t, 2, member2.ID)
				assert.Equal(t, "Jane", member2.FirstName)
				assert.Equal(t, "Smith", member2.LastName)
				assert.Equal(t, "Jane.Smith@gmail.com", member2.Email)
				assert.Equal(t, "1985-05-05", member2.DateOfBirth)
			}
		})
	}
}

func TestUpdateMemberById(t *testing.T) {
	t.Parallel()

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("members").ClientType(mtest.Mock))

	memberId := 1
	member := models.UpdateMember{
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "John.Doe@gmail.com",
		DateOfBirth: "1990-01-01",
	}

	testCases := []struct {
		name        string
		mongoDbMock func(mt *mtest.T)
		wantErr     bool
	}{
		{
			name: "success updating existing member",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: false,
		},
		{
			name: "update non-existing member",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    11000,
					Message: "member not found",
				}))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mt *mtest.T) {
			tc.mongoDbMock(mt)
			repo := NewMembershipRepository(mt.DB)
			err := repo.UpdateMemberById(context.Background(), &member, memberId)

			if tc.wantErr {
				assert.Errorf(t, err, "Want error but got: %v", err)
			} else {
				assert.NoErrorf(t, err, "Not expecting error")
			}
		})
	}
}

func TestDeleteMemberByID(t *testing.T) {
	t.Parallel()

	mt := mtest.New(t, mtest.NewOptions().DatabaseName("members").ClientType(mtest.Mock))

	testCases := []struct {
		name        string
		mongoDbMock func(mt *mtest.T)
		wantErr     bool
	}{
		{
			name: "success deleting existing member",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateSuccessResponse())
			},
			wantErr: false,
		},
		{
			name: "delete non-existing member",
			mongoDbMock: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
					Index:   0,
					Code:    11000,
					Message: "member not found",
				}))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		mt.Run(tc.name, func(mt *mtest.T) {
			memberId := 123
			tc.mongoDbMock(mt)
			repo := NewMembershipRepository(mt.DB)
			err := repo.DeleteMemberById(context.Background(), memberId)

			if tc.wantErr {
				assert.Errorf(t, err, "Want error but got: %v", err)
			} else {
				assert.NoErrorf(t, err, "Not expecting error")
			}
		})
	}
}
