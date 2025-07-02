package repository

import (
	"context"
	"errors"
	"fmt"
	"jsonjunk/config"
	"jsonjunk/internal/model"
	"jsonjunk/pkg/idgen"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoPasteRepository() Repository {
	coll := config.MongoClient.Database("jsonjunk").Collection("paste")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}}, // 1 오름 ,-1 내림
		Options: options.Index().SetUnique(true),
	}

	if _, err := coll.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		fmt.Printf("failed to created id index : %v\n", err)
	}

	return &mongoRepository{coll: coll}
}

func (r *mongoRepository) Insert(p model.Paste) error {
	for i := 0; i < 3; i++ {
		_, err := r.coll.InsertOne(context.TODO(), p)
		if err == nil {
			return nil
		}

		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 {
					p.ID = idgen.GenerateUUID()
				}
			}
		} else {
			return err
		}
	}
	return fmt.Errorf("failed to create uuid")
}

func (r *mongoRepository) SearchPasteByID(id string) (*model.Paste, error) {
	var result model.Paste
	err := r.coll.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *mongoRepository) TestSearchPastedAll() ([]*model.Paste, error) {
	cursor, err := r.coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []*model.Paste
	for cursor.Next(context.TODO()) {
		var paste model.Paste
		if err := cursor.Decode(&paste); err != nil {
			return nil, err
		}
		results = append(results, &paste)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
