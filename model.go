package mongoose

import (
	"fmt"
	"time"

	"dario.cat/mergo"
)

type Model struct {
	// ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func Merge(model interface{}) map[string]interface{} {
	base := Model{
		// ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	baseMap := make(map[string]interface{})
	if err := mergo.Map(&baseMap, base); err != nil {
		panic(err)
	}
	fmt.Printf("Base map is: %v\n", baseMap)

	modelMap := make(map[string]interface{})
	if err := mergo.Map(&modelMap, model); err != nil {
		panic(err)
	}
	fmt.Printf("Model map is: %v\n", modelMap)

	if err := mergo.Map(&baseMap, modelMap); err != nil {
		panic(err)
	}

	fmt.Printf("Merged map is: %v\n", baseMap)

	return modelMap
}
