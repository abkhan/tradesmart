package mongodb

import (
	"context"

	"mongotest/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository handles database operations for Trade records.
type Repository struct {
	coll *mongo.Collection
}

// NewRepository creates a new Trade repository.
func NewRepository(db *mongo.Database, collectionName string) *Repository {
	return &Repository{
		coll: db.Collection(collectionName),
	}
}

// CreateUniqueTrackingIndex ensures the tracking field has a unique index.
func (r *Repository) CreateUniqueTrackingIndex(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tracking", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

// BulkUpsertByTracking performs a bulk upsert of trade records based on the tracking field.
func (r *Repository) BulkUpsertByTracking(ctx context.Context, trades []models.Trade) (*mongo.BulkWriteResult, error) {
	if len(trades) == 0 {
		return &mongo.BulkWriteResult{}, nil
	}

	var models []mongo.WriteModel
	for _, t := range trades {
		filter := bson.D{{Key: "tracking", Value: t.Tracking}}
		update := bson.D{{Key: "$set", Value: t}}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false)
	return r.coll.BulkWrite(ctx, models, opts)
}

// InsertMany inserts multiple trade records.
func (r *Repository) InsertMany(ctx context.Context, trades []models.Trade) (*mongo.InsertManyResult, error) {
	if len(trades) == 0 {
		return &mongo.InsertManyResult{}, nil
	}

	var docs []interface{}
	for _, t := range trades {
		docs = append(docs, t)
	}

	return r.coll.InsertMany(ctx, docs)
}
