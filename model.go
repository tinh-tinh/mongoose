package mongoose

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type BaseTimestamp struct {
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type Model[M any] struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

func NewModel[M any](connect *Connect, name string) *Model[M] {
	return &Model[M]{
		Ctx:        connect.Ctx,
		Collection: connect.Client.Database(connect.DB).Collection(name),
	}
}

func (m Model[M]) Create(input *M) (*mongo.InsertOneResult, error) {
	ct := reflect.ValueOf(input).Elem()
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
	result, err := m.Collection.InsertOne(m.Ctx, input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m Model[M]) CreateMany(input []*M) (*mongo.InsertManyResult, error) {
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

func (m Model[M]) Find(filter interface{}, opt ...*options.FindOptions) ([]*M, error) {
	var data []*M

	query, err := ToDoc(filter)
	if err != nil {
		return data, err
	}
	cur, err := m.Collection.Find(m.Ctx, query, opt...)
	if err != nil {
		return data, err
	}

	for cur.Next(m.Ctx) {
		var t M
		err := cur.Decode(&t)
		if err != nil {
			return data, err
		}
		data = append(data, &t)
	}

	if err := cur.Err(); err != nil {
		return data, err
	}

	cur.Close(m.Ctx)

	return data, nil
}

func (m Model[M]) FindOne(filter interface{}, opt ...*options.FindOneOptions) (*M, error) {
	var data M

	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}

	err = m.Collection.FindOne(m.Ctx, query, opt...).Decode(&data)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &data, nil
}

func (m Model[M]) FindByID(id string, opt ...*options.FindOneOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}

	return m.FindOne(query, opt...)
}

func (m Model[M]) Update(filter interface{}, data interface{}) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	update := []bson.E{{
		Key:   "updatedAt",
		Value: time.Now(),
	}}
	ct := reflect.ValueOf(data).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)

		nameTag := field.Tag.Get("bson")
		name := field.Type.Name()
		val := ct.Field(i).Interface()
		if nameTag != "" && name != "BaseSchema" && !reflect.ValueOf(val).IsZero() {
			update = append(update, bson.E{
				Key:   nameTag,
				Value: val,
			})
		}
	}

	_, err = m.Collection.UpdateOne(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	return nil
}

func (m Model[M]) UpdateMany(filter interface{}, data interface{}) error {
	query, err := ToDoc(filter)
	if err != nil {
		return err
	}

	update := []bson.E{{
		Key:   "updatedAt",
		Value: time.Now(),
	}}
	ct := reflect.ValueOf(data).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)

		nameTag := field.Tag.Get("bson")
		name := field.Type.Name()
		val := ct.Field(i).Interface()
		if nameTag != "" && name != "BaseSchema" && !reflect.ValueOf(val).IsZero() {
			update = append(update, bson.E{
				Key:   nameTag,
				Value: val,
			})
		}
	}

	_, err = m.Collection.UpdateMany(m.Ctx, query, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}

	return nil
}

func (m Model[M]) DeleteOne(filter interface{}) error {
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

func (m Model[M]) DeleteMany(filter interface{}) error {
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

func ToDoc(v interface{}) (doc *bson.D, err error) {
	if v == nil {
		return &bson.D{}, nil
	}
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
