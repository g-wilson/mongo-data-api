package mongoapi

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

type FindOptions struct {
	underlying map[string]interface{}
}

func NewFindOptions() FindOptions {
	return FindOptions{underlying: map[string]interface{}{}}
}

func (o FindOptions) WithProjection(proj bson.M) FindOptions {
	o.underlying["projection"] = proj

	return o
}

func (o FindOptions) WithSort(sort bson.D) FindOptions {
	o.underlying["sort"] = sort

	return o
}

func (o FindOptions) WithLimit(limit uint64) FindOptions {
	o.underlying["limit"] = limit

	return o
}

func (o FindOptions) WithSkip(skip uint64) FindOptions {
	o.underlying["skip"] = skip

	return o
}

func (o FindOptions) Apply(params bson.M) {
	for k, v := range o.underlying {
		params[k] = v
	}
}

type FindResponse struct {
	Error     error
	Documents []json.RawMessage
}

func (r FindResponse) All(dest interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	resultsVal := reflect.ValueOf(dest)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("mongoapi: results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}

	sliceVal := resultsVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("mongoapi: results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
	}

	elementType := sliceVal.Type().Elem()

	for _, d := range r.Documents {
		newElem := reflect.New(elementType)

		err := bson.UnmarshalExtJSON(d, false, newElem.Interface())
		if err != nil {
			return fmt.Errorf("mongoapi: find: failed to unmarshal response document: %w", err)
		}

		sliceVal = reflect.Append(sliceVal, newElem.Elem())
	}

	resultsVal.Elem().Set(sliceVal.Slice(0, sliceVal.Cap()))
	return nil
}

func (c *Collection) Find(ctx context.Context, filter bson.M, options ...FindOptions) FindResponse {
	body := bson.M{
		"collection": c.name,
		"filter":     filter,
	}

	for _, o := range options {
		o.Apply(body)
	}

	res := struct {
		Documents []json.RawMessage
	}{}
	err := c.database.do(ctx, ActionFind, body, &res)
	if err != nil {
		return FindResponse{Error: err}
	}
	if len(res.Documents) < 1 {
		return FindResponse{Error: ErrNoDocuments}
	}

	return FindResponse{Documents: res.Documents}
}
