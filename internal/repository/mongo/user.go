package mongo

import (
	"context"
	"fmt"
	"gRPC_Example/internal/core"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

func (repository *UserRepository) GetById(ctx context.Context, id string) (*core.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	userCh := make(chan *core.User)
	errCh := make(chan error)

	var err error

	go repository.retrieveUser(ctxTimeout, id, userCh, errCh)

	if err != nil {
		return nil, err
	}

	select {
	case <-ctxTimeout.Done():
		fmt.Println("Processing timeout in Mongo")
		return nil, ctxTimeout.Err()
	case user := <-userCh:
		fmt.Println("Finished processing in Mongo")
		return user, nil
	case err := <-errCh:
		fmt.Println("error retrive user", err)
		return nil, err
	}

}

func (repository *UserRepository) retrieveUser(ctx context.Context, id string, userCh chan *core.User, errCh chan error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errCh <- err
	}

	user := &core.User{}

	if err := repository.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(user); err != nil {
		errCh <- err
	}

	userCh <- user

}
