package student

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentRepo interface {
	GetAll() ([]Student, error)
	Create(s Student) error
}

type studentRepo struct {
	db *mongo.Collection
}

func NewStudentRepo(db *mongo.Database) StudentRepo {
	return &studentRepo{
		db: db.Collection("student"),
	}
}

func (s *studentRepo) Create(data Student) error {
	ctx := context.Background()
	_, err := s.db.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *studentRepo) GetAll() ([]Student, error) {
	ctx := context.Background()
	filter := bson.D{}
	opt := options.Find()
	opt.SetLimit(10)
	cursor, err := s.db.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	var results []Student
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
