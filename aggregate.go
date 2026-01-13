package mongoose

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *Model[M]) Aggregate(pipeline mongo.Pipeline) ([]bson.M, error) {
	cursor, err := m.Collection.Aggregate(m.Ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M

	if err = cursor.All(m.Ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *Model[M]) FindOne(filter interface{}, opts ...QueryOptions) (*M, error) {
	err := ExecutePreHook(FindOne, m, filter)
	if err != nil {
		return nil, err
	}

	if err := m.sanitizeFilter(filter); err != nil {
		return nil, err
	}

	pipeline := []bson.M{}
	// Filter by search

	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}
	aggSearch := bson.M{"$match": query}
	pipeline = append(pipeline, aggSearch)

	var opt QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Ref != nil {
		for _, ref := range opt.Ref {
			refPath := m.getRefPath(ref)
			aggLookup := bson.M{"$lookup": bson.M{
				"from":         refPath.From,
				"localField":   refPath.ForeignKey,
				"foreignField": "_id",
				"as":           refPath.As,
			}}
			aggUnwind := bson.M{"$unwind": bson.M{
				"path": fmt.Sprintf("$%s", refPath.As),
			}}
			pipeline = append(pipeline, aggLookup, aggUnwind)
		}
	}

	if opt.Projection != nil {
		pipeline = append(pipeline, bson.M{"$project": opt.Projection})
	}

	if opt.Sort != nil {
		pipeline = append(pipeline, bson.M{"$sort": opt.Sort})
	}

	pipeline = append(pipeline, bson.M{"$limit": 1})

	cursor, err := m.Collection.Aggregate(m.Ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var data []*M
	for cursor.Next(m.Ctx) {
		var t M
		err := cursor.Decode(&t)
		if err != nil {
			return nil, err
		}
		data = append(data, &t)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(m.Ctx)

	if len(data) == 0 {
		return nil, nil
	}

	err = ExecutePostHook(FindOne, m, data[0])
	if err != nil {
		return nil, err
	}
	return data[0], nil
}

func (m *Model[M]) Find(filter interface{}, opts ...QueriesOptions) ([]*M, error) {
	err := ExecutePreHook(Find, m, filter)
	if err != nil {
		return nil, err
	}

	if err := m.sanitizeFilter(filter); err != nil {
		return nil, err
	}

	pipeline := []bson.M{}
	// Filter by search

	query, err := ToDoc(filter)
	if err != nil {
		return nil, err
	}
	aggSearch := bson.M{"$match": query}
	pipeline = append(pipeline, aggSearch)

	var opt QueriesOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Ref != nil {
		for _, ref := range opt.Ref {
			refPath := m.getRefPath(ref)
			aggLookup := bson.M{"$lookup": bson.M{
				"from":         refPath.From,
				"localField":   refPath.ForeignKey,
				"foreignField": "_id",
				"as":           refPath.As,
			}}
			aggUnwind := bson.M{"$unwind": bson.M{
				"path": fmt.Sprintf("$%s", refPath.As),
			}}
			pipeline = append(pipeline, aggLookup, aggUnwind)
		}
	}

	if opt.Projection != nil {
		pipeline = append(pipeline, bson.M{"$project": opt.Projection})
	}

	if opt.Sort != nil {
		pipeline = append(pipeline, bson.M{"$sort": opt.Sort})
	}

	if opt.Skip != 0 {
		pipeline = append(pipeline, bson.M{"$skip": opt.Skip})
	}

	if opt.Limit != 0 {
		pipeline = append(pipeline, bson.M{"$limit": opt.Limit})
	}

	cursor, err := m.Collection.Aggregate(m.Ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var data []*M
	for cursor.Next(m.Ctx) {
		var t M
		err := cursor.Decode(&t)
		if err != nil {
			return nil, err
		}
		data = append(data, &t)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(m.Ctx)

	err = ExecutePostHook(Find, m, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type RefPath struct {
	From       string
	ForeignKey string
	As         string
}

func (m *Model[M]) getRefPath(ref string) *RefPath {
	var model M
	ct := reflect.ValueOf(&model).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		tag := field.Tag.Get("ref")
		if tag != "" {
			refStr := strings.Split(tag, "->")
			foreignKey := refStr[0]
			as := field.Tag.Get("bson")
			foreignCol := refStr[1]
			if foreignKey == ref {
				return &RefPath{
					From:       foreignCol,
					ForeignKey: foreignKey,
					As:         as,
				}
			}
		}
	}
	return nil
}
