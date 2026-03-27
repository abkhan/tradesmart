package mongodb

import (
	"context"
	"time"

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

// GetByOrderID finds a record by seller_order_id.
func (r *Repository) GetByOrderID(ctx context.Context, orderID string) (*models.Trade, error) {
	var t models.Trade
	err := r.coll.FindOne(ctx, bson.M{"seller_order_id": orderID}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetByTracking finds a record by tracking number.
func (r *Repository) GetByTracking(ctx context.Context, tracking string) (*models.Trade, error) {
	var t models.Trade
	err := r.coll.FindOne(ctx, bson.M{"tracking": tracking}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetByDate finds all records for a specific date (ignoring time).
func (r *Repository) GetByDate(ctx context.Context, date time.Time) ([]models.Trade, error) {
	// Calculate start and end of the day
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)

	filter := bson.M{
		"order_date": bson.M{"$gte": start, "$lt": end},
	}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trades []models.Trade
	if err := cursor.All(ctx, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

// Search finds records by a general search term (matching seller_order_id, tracking, or order_type).
func (r *Repository) Search(ctx context.Context, query string) ([]models.Trade, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"seller_order_id": bson.M{"$regex": query, "$options": "i"}},
			{"tracking": bson.M{"$regex": query, "$options": "i"}},
			{"order_type": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trades []models.Trade
	if err := cursor.All(ctx, &trades); err != nil {
		return nil, err
	}
	return trades, nil
}

// Summary represents aggregated statistics.
type Summary struct {
	Count   int64   `json:"count"`
	Sales   float64 `json:"sales"`
	Profit  float64 `json:"profit"`
	Returns float64 `json:"returns"`
	Net     float64 `json:"net"`
}

// GetSummary calculates aggregated stats for a date range.
func (r *Repository) GetSummary(ctx context.Context, start, end interface{}) (*Summary, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"order_date": bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":     nil,
			"count":   bson.M{"$sum": 1},
			"sales":   bson.M{"$sum": "$sale_price_1"},
			"profit":  bson.M{"$sum": "$profit"},
			"returns": bson.M{"$sum": "$return_refund"},
			"net":     bson.M{"$sum": "$net"},
		}}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &Summary{}, nil
	}

	res := results[0]
	return &Summary{
		Count:   toInt64(res["count"]),
		Sales:   toDouble(res["sales"]),
		Profit:  toDouble(res["profit"]),
		Returns: toDouble(res["returns"]),
		Net:     toDouble(res["net"]),
	}, nil
}

func toDouble(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	}
	return 0
}

func toInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int32:
		return int64(val)
	case int64:
		return val
	}
	return 0
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
