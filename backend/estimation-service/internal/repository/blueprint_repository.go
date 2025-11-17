package repository

import (
	"context"
	"errors"

	"github.com/trademaster/backend/estimation-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlueprintRepository interface {
	Create(blueprint *models.Blueprint) error
	FindByID(id string) (*models.Blueprint, error)
	FindByEstimateID(estimateID string) ([]*models.Blueprint, error)
	Update(blueprint *models.Blueprint) error
	Delete(id string) error
}

type blueprintRepository struct {
	collection *mongo.Collection
}

func NewBlueprintRepository(db *mongo.Database) BlueprintRepository {
	return &blueprintRepository{
		collection: db.Collection("blueprints"),
	}
}

func (r *blueprintRepository) Create(blueprint *models.Blueprint) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, blueprint)
	return err
}

func (r *blueprintRepository) FindByID(id string) (*models.Blueprint, error) {
	ctx := context.Background()
	var blueprint models.Blueprint
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&blueprint)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("blueprint not found")
		}
		return nil, err
	}
	return &blueprint, nil
}

func (r *blueprintRepository) FindByEstimateID(estimateID string) ([]*models.Blueprint, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{"estimate_id": estimateID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blueprints []*models.Blueprint
	if err := cursor.All(ctx, &blueprints); err != nil {
		return nil, err
	}
	return blueprints, nil
}

func (r *blueprintRepository) Update(blueprint *models.Blueprint) error {
	ctx := context.Background()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": blueprint.ID}, blueprint)
	return err
}

func (r *blueprintRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
