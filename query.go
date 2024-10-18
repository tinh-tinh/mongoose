package mongoose

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryOptions struct {
	Sort       bson.D
	Projection bson.D
}

type QueriesOptions struct {
	Sort       bson.D
	Skip       int64
	Limit      int64
	Projection bson.D
}

func ParseFindOneOptions(opt QueryOptions) *options.FindOneOptions {
	opts := options.FindOne()
	if opt.Projection != nil {
		opts.SetProjection(opt.Projection)
	}
	if opt.Sort != nil {
		opts.SetSort(opt.Sort)
	}

	return opts
}

func ParseFindOptions(opt QueriesOptions) *options.FindOptions {
	opts := options.Find()
	if opt.Sort != nil {
		opts.SetSort(opt.Sort)
	}
	if opt.Projection != nil {
		opts.SetProjection(opt.Projection)
	}
	if opt.Skip != 0 {
		opts.SetSkip(opt.Skip)
	}
	if opt.Limit != 0 {
		opts.SetLimit(opt.Limit)
	}

	return opts
}

func (m *Model[M]) Find(filter interface{}, opt ...QueriesOptions) ([]*M, error) {
	var data []*M

	query, err := ToDoc(filter)
	if err != nil {
		return data, err
	}

	opts := []*options.FindOptions{}
	for i := 0; i < len(opt); i++ {
		opts = append(opts, ParseFindOptions(opt[i]))
	}

	cur, err := m.Collection.Find(m.Ctx, query, opts...)
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

func (m *Model[M]) FindOne(filter interface{}, opt ...QueryOptions) (*M, error) {
	var data M

	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}

	opts := []*options.FindOneOptions{}
	for i := 0; i < len(opt); i++ {
		opts = append(opts, ParseFindOneOptions(opt[i]))
	}

	err = m.Collection.FindOne(m.Ctx, query, opts...).Decode(&data)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &data, nil
}

func (m *Model[M]) FindByID(id string, opt ...QueryOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}

	return m.FindOne(query, opt...)
}

func (m *Model[M]) Count(filter interface{}) (int64, error) {
	query, err := ToDoc(filter)
	if err != nil {
		return 0, err
	}

	return m.Collection.CountDocuments(m.Ctx, query)
}

func (m *Model[M]) FindOneAndUpdate(filter interface{}, data *M, opt ...*options.FindOneAndUpdateOptions) (*M, error) {
	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}

	upsert := m.validData(data, "update")

	var model M
	err = m.Collection.FindOneAndUpdate(m.Ctx, query, bson.D{{Key: "$set", Value: upsert}}, opt...).Decode(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (m *Model[M]) FindByIDAndUpdate(id string, data *M, opt ...*options.FindOneAndUpdateOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}
	return m.FindOneAndUpdate(query, data, opt...)
}

func (m *Model[M]) FindOneAndDelete(filter interface{}, opt ...*options.FindOneAndDeleteOptions) (*M, error) {
	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}

	var model M
	err = m.Collection.FindOneAndDelete(m.Ctx, query, opt...).Decode(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (m *Model[M]) FindByIDAndDelete(id string, opt ...*options.FindOneAndDeleteOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}
	return m.FindOneAndDelete(query, opt...)
}

func (m *Model[M]) FindOneAndReplace(filter interface{}, data *M, opt ...*options.FindOneAndReplaceOptions) (*M, error) {
	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}

	update := m.validData(data, "replace")
	var model M
	err = m.Collection.FindOneAndReplace(m.Ctx, query, update, opt...).Decode(&model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (m *Model[M]) FindByIDAndReplace(id string, data *M, opt ...*options.FindOneAndReplaceOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}
	return m.FindOneAndReplace(query, data, opt...)
}
