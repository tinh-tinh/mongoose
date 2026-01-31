package mongoose

import (
	"reflect"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/dto/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validation = validator.Validator{}

// Create creates a new document in the collection using the provided input.
// It validates the input data and inserts a new document into the collection.
// Returns the result of the insertion as an *mongo.InsertOneResult and any error encountered.
func (m *Model[M]) Create(input *M) (*mongo.InsertOneResult, error) {
	err := m.beforeInsert(input)
	if err != nil {
		return nil, err
	}

	err = ExecutePreHook(Create, m)
	if err != nil {
		return nil, err
	}
	result, err := m.Collection.InsertOne(m.Ctx, input)
	if err != nil {
		return nil, err
	}

	err = ExecutePostHook(Create, m, result)
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
		err := m.beforeInsert(v)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}

	err := ExecutePreHook(CreateMany, m, input)
	if err != nil {
		return nil, err
	}

	result, err := m.Collection.InsertMany(m.Ctx, data)
	if err != nil {
		return nil, err
	}

	err = ExecutePostHook(CreateMany, m, result)
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
	err := ExecutePreHook(Update, m, filter, data)
	if err != nil {
		return err
	}

	if err := m.sanitizeFilter(filter); err != nil {
		return err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	update, err := m.beforeUpdate(data, false)
	if err != nil {
		return err
	}
	_, err = m.Collection.UpdateOne(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	err = ExecutePostHook(Update, m)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model[M]) UpdateByID(id interface{}, data *M) error {
	query, err := m.getQueryId(id)
	if err != nil {
		return err
	}

	return m.Update(query, data)
}

// UpdateMany updates multiple documents in the collection based on the provided filter and new data.
// It converts the filter to a BSON document using ToDoc.
// The new data is validated and prepared for update using m.validData.
// Finally, it performs the update operation using UpdateMany with the $set operator.
// Returns an error if the update operation fails.
func (m *Model[M]) UpdateMany(filter interface{}, data *M) error {
	err := ExecutePreHook(UpdateMany, m, filter, data)
	if err != nil {
		return err
	}

	if err := m.sanitizeFilter(filter); err != nil {
		return err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	update, err := m.beforeUpdate(data, false)
	if err != nil {
		return err
	}

	_, err = m.Collection.UpdateMany(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	err = ExecutePostHook(UpdateMany, m)
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
	err := ExecutePreHook(Delete, m, filter)
	if err != nil {
		return err
	}

	if err := m.sanitizeFilter(filter); err != nil {
		return err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	_, err = m.Collection.DeleteOne(m.Ctx, query)
	if err != nil {
		return err
	}

	err = ExecutePostHook(Delete, m)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model[M]) DeleteByID(id interface{}) error {
	query, err := m.getQueryId(id)
	if err != nil {
		return err
	}
	return m.Delete(query)
}

// DeleteMany deletes multiple documents in the collection based on the provided filter.
// It converts the filter to a BSON document using ToDoc.
// Finally, it performs the delete operation using DeleteMany.
// Returns an error if the delete operation fails.
func (m *Model[M]) DeleteMany(filter interface{}) error {
	err := ExecutePreHook(DeleteMany, m, filter)
	if err != nil {
		return err
	}

	if err := m.sanitizeFilter(filter); err != nil {
		return err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	_, err = m.Collection.DeleteMany(m.Ctx, query)
	if err != nil {
		return err
	}

	err = ExecutePostHook(DeleteMany, m)
	if err != nil {
		return err
	}

	return nil
}

// beforeInsert validates and prepares the data for insert.
// It iterates over the struct fields of the given data and sets the corresponding field of the model to the value of the field.
// It sets the createdAt and updatedAt to the current time.
func (m *Model[M]) beforeInsert(data *M) error {
	if m.option.Validation {
		err := ExecutePreHook(Validate, m, data)
		if err != nil {
			return err
		}

		err = validation.Validate(data)
		if err != nil {
			return err
		}

		err = ExecutePostHook(Validate, m)
		if err != nil {
			return err
		}
	}

	typeInfo := GetTypeInfo[M]()
	ct := reflect.ValueOf(data).Elem()
	// Use cached field info to find BaseSchema/BaseTimestamp fields (only top-level)
	for _, field := range typeInfo.Fields {
		// Only process top-level fields (IndexPath length == 1)
		if len(field.IndexPath) != 1 {
			continue
		}
		switch field.TypeName {
		case "BaseSchema":
			ct.Field(field.Index).Set(reflect.ValueOf(BaseSchema{
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}))
		case "BaseTimestamp":
			ct.Field(field.Index).Set(reflect.ValueOf(BaseTimestamp{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}))
		}
	}
	return nil
}

// beforeUpdate validates and prepares the data for update/replace.
// It validates the data and constructs a bson.E slice for the update.
// If isReplace is true, it sets createdAt to current time (for replacement).
// It always sets updatedAt to current time.
// It respects "readonly" tags.
func (m *Model[M]) beforeUpdate(data *M, isReplace bool) ([]bson.E, error) {
	if m.option.Validation {
		err := ExecutePreHook(Validate, m, data)
		if err != nil {
			return nil, err
		}

		err = validation.Validate(data)
		if err != nil {
			return nil, err
		}

		err = ExecutePostHook(Validate, m)
		if err != nil {
			return nil, err
		}
	}

	upsert := []bson.E{}
	typeInfo := GetTypeInfo[M]()

	if m.option.Timestamp {
		if isReplace {
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
	// Use cached field info for field iteration (only top-level fields)
	for _, field := range typeInfo.Fields {
		// Only process top-level fields (IndexPath length == 1)
		if len(field.IndexPath) != 1 {
			continue
		}

		// Skip embedded base types
		if field.TypeName == "BaseSchema" || field.TypeName == "BaseTimestamp" {
			continue
		}

		val := ct.Field(field.Index).Interface()

		// Skip readonly fields during update/replace to prevent mass assignment
		if field.MongooseTag == "readonly" {
			continue
		}

		if field.BsonTag != "" && !reflect.ValueOf(val).IsZero() {
			upsert = append(upsert, bson.E{
				Key:   field.BsonTag,
				Value: val,
			})
		}
	}

	return upsert, nil
}
