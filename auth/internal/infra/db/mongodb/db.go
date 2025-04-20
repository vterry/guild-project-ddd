package mongodb

import (
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBStorage(uri, username, password string) *options.ClientOptions {
	clientOptions := options.Client().ApplyURI(uri).
		SetAuth(options.Credential{
			Username:      username,
			Password:      password,
			AuthSource:    "admin",
			AuthMechanism: "SCRAM-SHA-256",
		})
	return clientOptions
}

func MongoURIBuilder(mongoAdd string) string {
	return strings.Join([]string{"mongodb://", mongoAdd, "/"}, "")
}
