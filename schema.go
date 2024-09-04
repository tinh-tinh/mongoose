package mongoose

import (
	"time"

	"dario.cat/mergo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema map[string]interface{}

type BaseSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

func NewSchema(schema interface{}) Schema {
	base := BaseSchema{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	baseMap := make(Schema, 3)

	if err := mergo.Map(&baseMap, base); err != nil {
		panic(err)
	}

	modelMap := make(Schema)
	if err := mergo.Map(&modelMap, schema); err != nil {
		panic(err)
	}

	if err := mergo.Map(&baseMap, modelMap); err != nil {
		panic(err)
	}

	return baseMap
}
