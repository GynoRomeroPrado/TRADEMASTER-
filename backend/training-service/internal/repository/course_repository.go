package repository

import (
	"context"
	"errors"

	"github.com/trademaster/backend/training-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CourseRepository interface {
	Create(course *models.Course) error
	FindByID(id string) (*models.Course, error)
	FindAll(filters map[string]interface{}, limit, offset int) ([]*models.Course, error)
	FindByTrade(trade string) ([]*models.Course, error)
	Update(course *models.Course) error
	Delete(id string) error
}

type courseRepository struct {
	collection *mongo.Collection
}

func NewCourseRepository(db *mongo.Database) CourseRepository {
	return &courseRepository{
		collection: db.Collection("courses"),
	}
}

func (r *courseRepository) Create(course *models.Course) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, course)
	return err
}

func (r *courseRepository) FindByID(id string) (*models.Course, error) {
	ctx := context.Background()
	var course models.Course
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&course)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindAll(filters map[string]interface{}, limit, offset int) ([]*models.Course, error) {
	ctx := context.Background()

	// Build query
	query := bson.M{"is_published": true}
	for k, v := range filters {
		query[k] = v
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var courses []*models.Course
	if err := cursor.All(ctx, &courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) FindByTrade(trade string) ([]*models.Course, error) {
	return r.FindAll(map[string]interface{}{"trade": trade}, 100, 0)
}

func (r *courseRepository) Update(course *models.Course) error {
	ctx := context.Background()
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": course.ID}, course)
	return err
}

func (r *courseRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
