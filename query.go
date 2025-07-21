package mongoose

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// FindByID returns a single document that matches the id. The id is the
// "_id" field in the document. The function takes a variable number of
// QueryOptions, which can be used to limit, sort, skip, or project the
// result. If no QueryOptions are provided, the first document in the collection
// is returned.
//
// FindByID returns an error if there is a problem with the query or the document
// cannot be decoded. If no document matches the id, the function returns
// nil, nil.
func (m *Model[M]) FindByID(id interface{}, opt ...QueryOptions) (*M, error) {
	query, err := m.getQueryId(id)
	if err != nil {
		return nil, err
	}

	return m.FindOne(query, opt...)
}

// Count returns the number of documents that match the filter. The filter can be any
// type that can be marshaled to a bson.D. It converts the filter to a BSON document
// using ToDoc. The function then counts the number of documents that match the query
// and returns the count as an int64. It returns an error if there is a problem with
// the query or the counting operation fails.
func (m *Model[M]) Count(filter interface{}) (int64, error) {
	err := ExecutePreHook(Count, m, filter)
	if err != nil {
		return 0, err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return 0, err
	}

	count, err := m.Collection.CountDocuments(m.Ctx, query)
	if err != nil {
		return 0, err
	}

	err = ExecutePostHook(Count, m)
	if err != nil {
		return 0, err
	}
	return count, nil
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
	err := ExecutePreHook(FindOneAndUpdate, m, filter, data)
	if err != nil {
		return nil, err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}
	upsert, err := m.serializeData(data, "update")
	if err != nil {
		return nil, err
	}

	var model M
	err = m.Collection.FindOneAndUpdate(m.Ctx, query, bson.D{{Key: "$set", Value: upsert}}, opt...).Decode(&model)
	if err != nil {
		return nil, err
	}

	err = ExecutePostHook(FindOneAndUpdate, m, model)
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
func (m *Model[M]) FindByIDAndUpdate(id any, data *M, opt ...*options.FindOneAndUpdateOptions) (*M, error) {
	query, err := m.getQueryId(id)
	if err != nil {
		return nil, err
	}

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
	err := ExecutePreHook(FindOneAndDelete, m, filter)
	if err != nil {
		return nil, err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}

	var model M
	err = m.Collection.FindOneAndDelete(m.Ctx, query, opt...).Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	err = ExecutePostHook(FindOneAndDelete, m, model)
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
func (m *Model[M]) FindByIDAndDelete(id any, opt ...*options.FindOneAndDeleteOptions) (*M, error) {
	query, err := m.getQueryId(id)
	if err != nil {
		return nil, err
	}

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
	err := ExecutePreHook(FindOneAndReplace, m, filter, data)
	if err != nil {
		return nil, err
	}

	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}

	update, err := m.serializeData(data, "replace")
	if err != nil {
		return nil, err
	}

	var model M
	err = m.Collection.FindOneAndReplace(m.Ctx, query, update, opt...).Decode(&model)
	if err != nil {
		return nil, err
	}

	err = ExecutePostHook(FindOneAndReplace, m, model)
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
func (m *Model[M]) FindByIDAndReplace(id any, data *M, opt ...*options.FindOneAndReplaceOptions) (*M, error) {
	query, err := m.getQueryId(id)
	if err != nil {
		return nil, err
	}

	return m.FindOneAndReplace(query, data, opt...)
}
