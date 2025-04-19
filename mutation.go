package mongoose

import (
	"reflect"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create creates a new document in the collection using the provided input.
// It validates the input data and inserts a new document into the collection.
// Returns the result of the insertion as an *mongo.InsertOneResult and any error encountered.
func (m *Model[M]) Create(input *M) (*mongo.InsertOneResult, error) {
	if m.option.Validation {
		err := validator.Scanner(input)
		if err != nil {
			return nil, err
		}
	}

	m.serializeData(input, "insert")
	result, err := m.Collection.InsertOne(m.Ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateMany creates multiple new documents in the collection using the provided input.
// It validates each document data and inserts new documents into the collection.
// Returns the result of the insertion as an *mongo.InsertManyResult and any error encountered.
func (m *Model[M]) CreateMany(input []*M) (*mongo.InsertManyResult, error) {
	data := make([]interface{}, 0)

	for _, v := range input {
		if m.option.Validation {
			err := validator.Scanner(v)
			if err != nil {
				return nil, err
			}
		}
		ct := reflect.ValueOf(v).Elem()
		for i := range ct.NumField() {
			field := ct.Type().Field(i)
			name := field.Type.Name()
			if name == "BaseSchema" {
				ct.FieldByName(name).Set(reflect.ValueOf(BaseSchema{
					ID:        primitive.NewObjectID(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}))
			} else if name == "BaseTimestamp" {
				ct.FieldByName(name).Set(reflect.ValueOf(BaseTimestamp{
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}))
			}
		}
		data = append(data, v)
	}

	result, err := m.Collection.InsertMany(m.Ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates a single document in the collection based on the provided filter and new data.
// It converts the filter to a BSON document using ToDoc.
// The new data is validated and prepared for update using m.validData.
// Finally, it performs the update operation using UpdateOne with the $set operator.
// Returns an error if the update operation fails.
func (m *Model[M]) Update(filter interface{}, data *M) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	if m.option.Validation {
		err := validator.Scanner(data)
		if err != nil {
			return err
		}
	}
	update := m.serializeData(data, "update")
	_, err = m.Collection.UpdateOne(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	return nil
}

// UpdateMany updates multiple documents in the collection based on the provided filter and new data.
// It converts the filter to a BSON document using ToDoc.
// The new data is validated and prepared for update using m.validData.
// Finally, it performs the update operation using UpdateMany with the $set operator.
// Returns an error if the update operation fails.
func (m *Model[M]) UpdateMany(filter interface{}, data *M) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	if m.option.Validation {
		err := validator.Scanner(data)
		if err != nil {
			return err
		}
	}
	update := m.serializeData(data, "update")

	_, err = m.Collection.UpdateMany(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a single document in the collection based on the provided filter.
// It converts the filter to a BSON document using ToDoc.
// Finally, it performs the delete operation using DeleteOne.
// Returns an error if the delete operation fails.
func (m *Model[M]) Delete(filter interface{}) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	_, err = m.Collection.DeleteOne(m.Ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMany deletes multiple documents in the collection based on the provided filter.
// It converts the filter to a BSON document using ToDoc.
// Finally, it performs the delete operation using DeleteMany.
// Returns an error if the delete operation fails.
func (m *Model[M]) DeleteMany(filter interface{}) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	_, err = m.Collection.DeleteMany(m.Ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// serializeData validates and prepares the data for update/insert.
// It iterates over the struct fields of the given data and sets the corresponding field of the model to the value of the field.
// If the field is not found in the model, it is not set.
// If the field is tagged with bson, the tag value is used as the key.
// If mutation is "insert", it sets the createdAt and updatedAt to the current time.
// If mutation is "replace", it sets the createdAt to the current time.
// If mutation is "update", it sets the updatedAt to the current time.
// The given data must be a struct.
// It returns the validated data as a bson.E slice.
func (m *Model[M]) serializeData(data *M, mutation string) []bson.E {
	upsert := []bson.E{}

	if mutation == "insert" {
		ct := reflect.ValueOf(data).Elem()
		for i := range ct.NumField() {
			field := ct.Type().Field(i)
			name := field.Type.Name()
			if name == "BaseSchema" {
				ct.FieldByName(name).Set(reflect.ValueOf(BaseSchema{
					ID:        primitive.NewObjectID(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}))
			}
		}
	} else {
		if m.option.Timestamp {
			if mutation == "replace" {
				upsert = append(upsert, bson.E{
					Key:   "createdAt",
					Value: time.Now(),
				})
			}
			upsert = append(upsert, bson.E{
				Key:   "updatedAt",
				Value: time.Now(),
			})
		}
		ct := reflect.ValueOf(data).Elem()
		for i := range ct.NumField() {
			field := ct.Type().Field(i)

			nameTag := field.Tag.Get("bson")
			name := field.Type.Name()
			val := ct.Field(i).Interface()
			if nameTag != "" && name != "BaseSchema" && !reflect.ValueOf(val).IsZero() {
				upsert = append(upsert, bson.E{
					Key:   nameTag,
					Value: val,
				})
			}
		}
	}

	return upsert
}
