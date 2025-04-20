package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"github.com/vterry/ddd-study/auth-server/internal/app/utils"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
	"github.com/vterry/ddd-study/auth-server/internal/domain/login"
	"github.com/vterry/ddd-study/auth-server/internal/infra/db/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewLoginRepository(ctx context.Context, db *mongo.Database) (*LoginRepository, error) {
	repo := &LoginRepository{
		collection: db.Collection("Login"),
		ctx:        ctx,
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := repo.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create unique index on userId - %v", err)
	}

	return repo, nil
}

func (r *LoginRepository) Save(loginObj login.Login) error {
	loginDao := dao.LoginToDAO(loginObj)
	loginDao.Id = primitive.NewObjectID()

	if err := utils.Validate.Struct(loginDao); err != nil {
		errors := err.(validator.ValidationErrors)
		return fmt.Errorf("error while parse dao object: %w", errors)
	}

	opts := options.InsertOne().SetBypassDocumentValidation(false)
	_, err := r.collection.InsertOne(r.ctx, loginDao, opts)
	if err != nil {
		return fmt.Errorf("error while saving object with id %v: %w", loginObj.UserId().ID(), err)
	}

	return nil
}

func (r *LoginRepository) FindLoginByUserID(userId valueobjects.UserID) (*login.Login, error) {
	loginDao := dao.Login{}
	filter := bson.D{{Key: "userId", Value: userId.ID().String()}}

	opts := options.FindOne().SetMaxTime(2 * time.Second)
	err := r.collection.FindOne(r.ctx, filter, opts).Decode(&loginDao)
	if err != nil {
		return nil, fmt.Errorf("error while fetching userId - %v: %w", userId.ID().String(), err)
	}

	resultLogin, err := dao.DAOtoLogin(loginDao)
	if err != nil {
		return nil, fmt.Errorf("error while parsing DAO to Model: %w", err)
	}

	return &resultLogin, nil
}

func (r *LoginRepository) UpdatePassword(userId valueobjects.UserID, password string) error {
	filter := bson.D{{Key: "userId", Value: userId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: password}}}}

	opts := options.Update().SetUpsert(false)
	result, err := r.collection.UpdateOne(r.ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("error while updating password for userId - %v: %w", userId.ID().String(), err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no register found with id: %v", userId.ID().String())
	}

	return nil
}
