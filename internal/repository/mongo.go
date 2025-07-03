package repository

import (
	"context"
	"errors"
	"fmt"
	"jsonjunk/config"
	"jsonjunk/internal/model"
	"jsonjunk/pkg/idgen"
	logger "jsonjunk/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoPasteRepository(dbName string) Repository {
	coll := config.MongoClient.Database(dbName).Collection("paste")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}}, // 1 오름 ,-1 내림
		Options: options.Index().SetUnique(true),
	}

	if _, err := coll.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		logger.Log.Error("failed to create index on 'id' field",
			zap.String("collection", "paste"),
			zap.Error(err),
		)
	}
	return &mongoRepository{coll: coll}
}

func (r *mongoRepository) Insert(ctx context.Context, p model.Paste) error {
	for i := 0; i < 3; i++ {

		_, err := r.coll.InsertOne(ctx, p)
		if err == nil {
			return nil
		}

		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			duplicate := false
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 {
					duplicate = true
					break
				}
			}
			if duplicate {
				p.ID = idgen.GenerateUUID()
				continue
			}
		}
		return fmt.Errorf("failed to insert paste: %w", err)
	}
	return fmt.Errorf("failed to generate unique uuid after 3 attempts")
}

func (r *mongoRepository) SearchPasteByID(ctx context.Context, id string) (*model.Paste, error) {
	var result model.Paste
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find paste: %w", err)
	}
	return &result, nil
}

func (r *mongoRepository) SearchPasteList(ctx context.Context) ([]*model.Paste, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Paste
	for cursor.Next(ctx) {
		var paste model.Paste
		if err := cursor.Decode(&paste); err != nil {
			return nil, err
		}
		results = append(results, &paste)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to found pastes ,%w", err)
	}

	return results, nil
}

func (r *mongoRepository) ModifyPaste(ctx context.Context, id string, fields map[string]interface{}) (paste model.Paste, err error) {
	filter := bson.M{"id": id}
	update := bson.M{"$set": fields}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // 업데이트 이후 문서 반환

	err = r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&paste)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return paste, fmt.Errorf("no paste found with id: %s", id)
		}
		return paste, fmt.Errorf("failed to update paste: %w", err)
	}
	return paste, nil
}

func (r *mongoRepository) DeletePaste(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	result, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete paste , %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no paste found with id: %s", id)
	}
	return nil
}
