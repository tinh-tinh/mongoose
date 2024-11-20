package mongoose

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryOptions struct {
	Sort       bson.D
	Projection bson.D
	Ref        []string
}

type QueriesOptions struct {
	Sort       bson.D
	Skip       int64
	Limit      int64
	Projection bson.D
	Ref        []string
}

// ParseFindOneOptions converts a QueryOptions to a *options.FindOneOptions.
//
// The QueryOptions can have a Projection or Sort set. If the Projection is
// non-nil, it is passed to the FindOneOptions.SetProjection method. If the Sort
// is non-nil, it is passed to the FindOneOptions.SetSort method.
// func ParseFindOneOptions(opt QueryOptions) *options.FindOneOptions {
// 	opts := options.FindOne()
// 	if opt.Projection != nil {
// 		opts.SetProjection(opt.Projection)
// 	}
// 	if opt.Sort != nil {
// 		opts.SetSort(opt.Sort)
// 	}

// 	return opts
// }

// ParseFindOptions converts a QueriesOptions to a *options.FindOptions.
//
// The QueriesOptions can have a Projection, Sort, Skip, or Limit set. If the
// Projection is non-nil, it is passed to the FindOptions.SetProjection method.
// If the Sort is non-nil, it is passed to the FindOptions.SetSort method. If
// the Skip is non-zero, it is passed to the FindOptions.SetSkip method. If the
// Limit is non-zero, it is passed to the FindOptions.SetLimit method.
// func ParseFindOptions(opt QueriesOptions) *options.FindOptions {
// 	opts := options.Find()
// 	if opt.Sort != nil {
// 		opts.SetSort(opt.Sort)
// 	}
// 	if opt.Projection != nil {
// 		opts.SetProjection(opt.Projection)
// 	}
// 	if opt.Skip != 0 {
// 		opts.SetSkip(opt.Skip)
// 	}
// 	if opt.Limit != 0 {
// 		opts.SetLimit(opt.Limit)
// 	}

// 	return opts
// }

// Find returns a slice of documents that match the filter. The filter can be any
// type that can be marshaled to a bson.D. The function takes a variable number
// of QueriesOptions, which can be used to limit, sort, skip, or project the
// results. If no QueriesOptions are provided, all documents are returned.
//
// Find returns an error if there is a problem with the query or the documents
// cannot be decoded.
// func (m *Model[M]) Find(filter interface{}, opt ...QueriesOptions) ([]*M, error) {
// 	var data []*M

// 	query, err := ToDoc(filter)
// 	if err != nil {
// 		return data, err
// 	}

// 	opts := []*options.FindOptions{}
// 	for i := 0; i < len(opt); i++ {
// 		opts = append(opts, ParseFindOptions(opt[i]))
// 	}

// 	cur, err := m.Collection.Find(m.Ctx, query, opts...)
// 	if err != nil {
// 		return data, err
// 	}

// 	for cur.Next(m.Ctx) {
// 		var t M
// 		err := cur.Decode(&t)
// 		if err != nil {
// 			return data, err
// 		}
// 		data = append(data, &t)
// 	}

// 	if err := cur.Err(); err != nil {
// 		return data, err
// 	}

// 	cur.Close(m.Ctx)

// 	return data, nil
// }

// FindOne returns a single document that matches the filter. The filter can be
// any type that can be marshaled to a bson.D. The function takes a variable
// number of QueryOptions, which can be used to limit, sort, skip, or project the
// result. If no QueryOptions are provided, the first document in the collection
// is returned.
//
// FindOne returns an error if there is a problem with the query or the document
// cannot be decoded. If no document matches the filter, the function returns
// nil, nil.
// func (m *Model[M]) FindOne(filter interface{}, opt ...QueryOptions) (*M, error) {
// 	var data M

// 	query, err := ToDoc(filter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	opts := []*options.FindOneOptions{}
// 	for i := 0; i < len(opt); i++ {
// 		opts = append(opts, ParseFindOneOptions(opt[i]))
// 	}

// 	err = m.Collection.FindOne(m.Ctx, query, opts...).Decode(&data)
// 	if err != nil && err != mongo.ErrNoDocuments {
// 		return nil, err
// 	}
// 	if err == mongo.ErrNoDocuments {
// 		return nil, nil
// 	}

// 	return &data, nil
// }

// FindByID returns a single document that matches the id. The id is the
// "_id" field in the document. The function takes a variable number of
// QueryOptions, which can be used to limit, sort, skip, or project the
// result. If no QueryOptions are provided, the first document in the collection
// is returned.
//
// FindByID returns an error if there is a problem with the query or the document
// cannot be decoded. If no document matches the id, the function returns
// nil, nil.
func (m *Model[M]) FindByID(id string, opt ...QueryOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}

	return m.FindOne(query, opt...)
}

// Count returns the number of documents that match the filter. The filter can be any
// type that can be marshaled to a bson.D. It converts the filter to a BSON document
// using ToDoc. The function then counts the number of documents that match the query
// and returns the count as an int64. It returns an error if there is a problem with
// the query or the counting operation fails.
func (m *Model[M]) Count(filter interface{}) (int64, error) {
	query, err := ToDoc(filter)
	if err != nil {
		return 0, err
	}

	return m.Collection.CountDocuments(m.Ctx, query)
}

// FindOneAndUpdate returns a single document that matches the filter and updates it with the
// new data. The filter can be any type that can be marshaled to a bson.D. The new data is
// validated and prepared for update using m.validData. The function takes a variable number
// of FindOneAndUpdateOptions, which can be used to control the update operation. If no
// FindOneAndUpdateOptions are provided, the first document in the collection that matches the
// filter is returned and updated.
//
// FindOneAndUpdate returns an error if there is a problem with the query or the document
// cannot be decoded. If no document matches the filter, the function returns nil, nil.
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

// FindByIDAndUpdate updates a single document that matches the given id with the new data.
// The id is the "_id" field in the document. The new data is validated and prepared for update
// using m.validData. The function takes a variable number of FindOneAndUpdateOptions,
// which can be used to control the update operation. If no FindOneAndUpdateOptions are provided,
// the first document in the collection that matches the id is returned and updated.
//
// FindByIDAndUpdate returns an error if there is a problem with the query or the document
// cannot be decoded. If no document matches the id, the function returns nil, nil.
func (m *Model[M]) FindByIDAndUpdate(id string, data *M, opt ...*options.FindOneAndUpdateOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}
	return m.FindOneAndUpdate(query, data, opt...)
}

// FindOneAndDelete deletes a single document that matches the filter and returns the deleted document.
// The filter can be any type that can be marshaled to a bson.D. The function takes a variable number
// of FindOneAndDeleteOptions, which can be used to control the delete operation. If no
// FindOneAndDeleteOptions are provided, the first document in the collection that matches the
// filter is deleted.
//
// FindOneAndDelete returns an error if there is a problem with the query or the document cannot
// be decoded. If no document matches the filter, the function returns nil, nil.
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

// FindByIDAndDelete deletes a single document that matches the id and returns the deleted document.
// The function takes a variable number of FindOneAndDeleteOptions, which can be used to control the
// delete operation. If no FindOneAndDeleteOptions are provided, the first document in the collection
// that matches the id is deleted.
//
// FindByIDAndDelete returns an error if there is a problem with the query or the document cannot
// be decoded. If no document matches the id, the function returns nil, nil.
func (m *Model[M]) FindByIDAndDelete(id string, opt ...*options.FindOneAndDeleteOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}
	return m.FindOneAndDelete(query, opt...)
}

// FindOneAndReplace replaces a single document that matches the filter with the new data.
// The filter can be any type that can be marshaled to a bson.D. The new data is
// validated and prepared for update using m.validData. The function takes a variable
// number of FindOneAndReplaceOptions, which can be used to control the replace operation.
// If no FindOneAndReplaceOptions are provided, the first document in the collection that
// matches the filter is replaced.
//
// FindOneAndReplace returns an error if there is a problem with the query or the document
// cannot be decoded. If no document matches the filter, the function returns nil, nil.
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

// FindByIDAndReplace replaces a single document that matches the given id with the new data.
// The id corresponds to the "_id" field in the document. The new data is validated and prepared
// for replacement using m.validData. The function takes a variable number of
// FindOneAndReplaceOptions, which can be used to control the replace operation. If no
// FindOneAndReplaceOptions are provided, the first document in the collection that matches the
// id is replaced.
//
// FindByIDAndReplace returns an error if there is a problem with the query or the document cannot
// be decoded. If no document matches the id, the function returns nil, nil.
func (m *Model[M]) FindByIDAndReplace(id string, data *M, opt ...*options.FindOneAndReplaceOptions) (*M, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	query := bson.M{"_id": objId}
	return m.FindOneAndReplace(query, data, opt...)
}
