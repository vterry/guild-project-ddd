package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	SessionId  string             `json:"sessionId" bson:"sessionId"`
	UserId     string             `json:"userId" bson:"userId"`
	JwtToken   string             `json:"jwtToken" bson:"jwtToken"`
	RenewToken string             `json:"renewToken" bson:"renewToken"`
	CsrfToken  string             `json:"csrfToken" bson:"csrfToken"`
	ExpiresAt  time.Time          `json:"expiresAt" bson:"expiresAt"`
	Revoked    bool               `json:"revoked" bson:"revoked"`
}

type Login struct {
	Id       primitive.ObjectID `json:"id" bson:"_id" validate:"required"`
	LoginID  string             `json:"loginId" bson:"loginId" validate:"required"`
	UserId   string             `json:"userId" bson:"userId" validate:"required"`
	Password string             `json:"password" bson:"password" validate:"required"`
}
