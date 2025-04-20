package mongodb

import (
	"context"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/vterry/ddd-study/auth-server/internal/app/utils"
	"github.com/vterry/ddd-study/auth-server/internal/domain/session"
	"github.com/vterry/ddd-study/auth-server/internal/infra/db/dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionRepository struct {
	collection mongo.Collection
	ctx        context.Context
}

func NewSessionRepository(ctx context.Context, db *mongo.Database) (*SessionRepository, error) {
	repo := SessionRepository{
		collection: *db.Collection("Session"),
		ctx:        ctx,
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "sessionId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := repo.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create unique index on sessionId - %v", err)
	}

	return &repo, nil
}

func (r *SessionRepository) Save(sessionObj session.Session) error {
	sessionDao := dao.SessionToDAO(sessionObj)
	sessionDao.Id = primitive.NewObjectID()

	if err := utils.Validate.Struct(sessionDao); err != nil {
		errors := err.(validator.ValidationErrors)
		return fmt.Errorf("error while parse dao object: %w", errors)
	}

	opts := options.InsertOne().SetBypassDocumentValidation(false)
	_, err := r.collection.InsertOne(r.ctx, sessionDao, opts)
	if err != nil {
		return fmt.Errorf("error while saving object with id %v: %w", sessionObj.SessionID.ID().String(), err)
	}

	return nil
}

func (s *SessionRepository) Update(sess session.Session) (*session.Session, error) {
	sessionDao := dao.SessionToDAO(sess)
	filter := bson.D{{Key: "sessionId", Value: sessionDao.SessionId}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "jwtToken", Value: sessionDao.JwtToken},
			{Key: "renewToken", Value: sessionDao.RenewToken},
			{Key: "csrfToken", Value: sessionDao.CsrfToken},
			{Key: "expiresAt", Value: sessionDao.ExpiresAt},
			{Key: "revoked", Value: sessionDao.Revoked},
		}},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedSession dao.Session
	err := s.collection.FindOneAndUpdate(s.ctx, filter, update, opts).Decode(&updatedSession)
	if err != nil {
		return nil, fmt.Errorf("error while updating session: %w", err)
	}

	domainSession, err := dao.DAOtoSession(updatedSession)
	if err != nil {
		return nil, fmt.Errorf("error while converting session DAO to domain model: %w", err)
	}

	return &domainSession, nil
}

func (s *SessionRepository) FindSessionByID(sessionID session.SessionID) (*session.Session, error) {
	filter := bson.D{{Key: "sessionId", Value: sessionID.ID().String()}}
	var sessionDao dao.Session

	err := s.collection.FindOne(s.ctx, filter).Decode(&sessionDao)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("session not found with id: %v", sessionID.ID().String())
		}
		return nil, fmt.Errorf("error while fetching session: %w", err)
	}

	domainSession, err := dao.DAOtoSession(sessionDao)
	if err != nil {
		return nil, fmt.Errorf("error while converting session DAO to domain model: %w", err)
	}

	return &domainSession, nil
}
