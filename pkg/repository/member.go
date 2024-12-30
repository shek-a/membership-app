package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"members.com/membership/pkg/models"
)

type MemberRepositoryI interface {
	CreateMember(ctx context.Context, member *models.Member) error
	GetMemberById(ctx context.Context, memberId int) (*models.Member, error)
	GetAllMembers(ctx context.Context) ([]models.Member, error)
	UpdateMemberById(ctx context.Context, member *models.UpdateMember, memberId int) error
	DeleteMemberById(ctx context.Context, memberId int) error
}

type MemberRepository struct {
	mongoDb *mongo.Database
}

func NewMembershipRepository(mongo *mongo.Database) MemberRepositoryI {
	return &MemberRepository{
		mongoDb: mongo,
	}
}

func (m *MemberRepository) CreateMember(ctx context.Context, member *models.Member) error {
	_, err := m.mongoDb.Collection("members").InsertOne(ctx, member)
	if err != nil {
		log.Println("error")
	}

	return err
}

func (m *MemberRepository) GetMemberById(ctx context.Context, memberId int) (*models.Member, error) {
	var member models.Member
	filter := bson.D{bson.E{Key: "ID", Value: memberId}}
	err := m.mongoDb.Collection("members").FindOne(ctx, filter).Decode(&member)
	return &member, err
}

func (m *MemberRepository) GetAllMembers(ctx context.Context) ([]models.Member, error) {
	query, err := m.mongoDb.Collection("members").Find(ctx, bson.D{})
	if err != nil {
		return []models.Member{}, err
	}
	defer query.Close(ctx)

	membersList := make([]models.Member, 0)
	for query.Next(ctx) {
		var row models.Member
		err := query.Decode(&row)
		if err != nil {
			log.Println("error decoding member:", err)
		}
		membersList = append(membersList, row)
	}
	return membersList, nil
}

func (m *MemberRepository) UpdateMemberById(ctx context.Context, member *models.UpdateMember, memberId int) error {
	filter := bson.M{"id": memberId}
	update := bson.M{
		"$set": bson.M{
			"FirstName":   member.FirstName,
			"LastName":    member.LastName,
			"Email":       member.Email,
			"DateOfBirth": member.DateOfBirth,
		},
	}

	_, err := m.mongoDb.Collection("members").UpdateOne(ctx, filter, update)
	return err
}

func (m *MemberRepository) DeleteMemberById(ctx context.Context, memberId int) error {
	filter := bson.M{"ID": memberId}
	_, err := m.mongoDb.Collection("members").DeleteOne(ctx, filter)
	if err != nil {
		log.Println("error deleting member:", err)
		return err
	}
	return nil
}
