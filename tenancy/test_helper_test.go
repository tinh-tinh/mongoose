package tenancy_test

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseSchema is a common base struct for test models with ID and timestamps
type BaseSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}
