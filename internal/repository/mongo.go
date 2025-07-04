package repository

import (
	"context"
	"errors"
	"fmt"
	"jsonjunk/config"
	"jsonjunk/internal/model"
	"jsonjunk/pkg/idgen"
	logger "jsonjunk/pkg/logging"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type mongoRepository struct {
	coll *mongo.Collection
}

func NewMongoPasteRepository(ctx context.Context, dbName string) Repository {
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

func (r *mongoRepository) InsertPaste(ctx context.Context, p model.Paste) error {
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
		return fmt.Errorf("%w: %v", model.ErrInsertFailed, err)
	}
	return fmt.Errorf("%w: failed to generate unique ID after 3 attempts", model.ErrDuplicatePasteID)
}

func (r *mongoRepository) SearchPasteByID(ctx context.Context, id string) (*model.Paste, error) {
	var result model.Paste
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", model.ErrDatabase, err)
	}
	return &result, nil
}

func (r *mongoRepository) SearchPasteList(ctx context.Context) ([]*model.Paste, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrDatabase, err)
	}
	defer cursor.Close(ctx)

	var results []*model.Paste
	for cursor.Next(ctx) {
		var paste model.Paste
		if err := cursor.Decode(&paste); err != nil {
			return nil, fmt.Errorf("%w: %v", model.ErrDatabase, err)
		}
		results = append(results, &paste)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrDatabase, err)
	}

	return results, nil
}

func (r *mongoRepository) UpdatePasteByID(ctx context.Context, id string, fields map[string]interface{}) (paste model.Paste, err error) {
	filter := bson.M{"id": id}
	update := bson.M{"$set": fields}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // 업데이트 이후 문서 반환

	err = r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&paste)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return paste, model.ErrPasteNotFound
		}
		return paste, fmt.Errorf("%w: %v", model.ErrDatabase, err)
	}
	return paste, nil
}

func (r *mongoRepository) DeletePasteByID(ctx context.Context, id string) error {
	filter := bson.M{
		"id":         id,
		"is_deleted": bson.M{"$ne": true}, // 이미 삭제된 건 제외
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	result, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%w: %v", model.ErrDatabase, err)
	}

	if result.MatchedCount == 0 {
		return model.ErrPasteNotFound
	}
	return nil
}
