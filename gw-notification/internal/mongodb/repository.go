package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"gw-notification/internal/models"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client, dbName, collName string) *Repository {
	return &Repository{
		collection: client.Database(dbName).Collection(collName),
	}
}

// SaveTransaction сохраняет транзакцию в MongoDB с проверкой идемпотентности
func (r *Repository) SaveTransaction(tx models.TransactionHistory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, есть ли уже такой transaction_id
	filter := bson.M{"transaction_id": tx.TransactionID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		// Уже сохранено → идемпотентно
		return nil
	}

	tx.CreatedAt = time.Now()
	_, err = r.collection.InsertOne(ctx, tx)
	return err
}
