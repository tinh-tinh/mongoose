package mongoose

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Model[M]) Create(input *M) (*mongo.InsertOneResult, error) {
	m.validData(input, "insert")
	result, err := m.Collection.InsertOne(m.Ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Model[M]) CreateMany(input []*M) (*mongo.InsertManyResult, error) {
	data := make([]interface{}, 0)

	for _, v := range input {
		ct := reflect.ValueOf(v).Elem()
		for i := 0; i < ct.NumField(); i++ {
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
		data = append(data, v)
	}

	result, err := m.Collection.InsertMany(m.Ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Model[M]) Update(filter interface{}, data *M) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	update := m.validData(data, "update")
	_, err = m.Collection.UpdateOne(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	return nil
}

func (m *Model[M]) UpdateMany(filter interface{}, data *M) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	update := m.validData(data, "update")

	_, err = m.Collection.UpdateMany(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	return nil
}

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

func (m *Model[M]) validData(data *M, mutation string) []bson.E {
	upsert := []bson.E{}
	if mutation == "insert" {
		ct := reflect.ValueOf(data).Elem()
		for i := 0; i < ct.NumField(); i++ {
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
		ct := reflect.ValueOf(data).Elem()
		for i := 0; i < ct.NumField(); i++ {
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
