package mongoose

import "go.mongodb.org/mongo-driver/bson/primitive"

func IsValidateObjectID(str string) bool {
	_, err := primitive.ObjectIDFromHex(str)
	return err == nil
}

func ToObjectID(str string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		panic(err)
	}
	return id
}
