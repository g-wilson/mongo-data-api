package mongoapi

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type FindOneOptions struct {
	underlying map[string]interface{}
}

func NewFindOneOptions() FindOneOptions {
	return FindOneOptions{underlying: map[string]interface{}{}}
}

func (o FindOneOptions) WithProjection(proj bson.M) FindOneOptions {
	o.underlying["projection"] = proj

	return o
}

func (o FindOneOptions) Apply(params bson.M) {
	for k, v := range o.underlying {
		params[k] = v
	}
}

type FindOneResponse struct {
	Error    error
	Document json.RawMessage
}

func (r FindOneResponse) Decode(dest interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	err := bson.UnmarshalExtJSON(r.Document, false, dest)
	if err != nil {
		return fmt.Errorf("mongoapi: findOne: failed to unmarshal response document: %w", err)
	}

	return nil
}

func (c *Collection) FindOne(ctx context.Context, filter bson.M, options ...FindOneOptions) FindOneResponse {
	body := bson.M{
		"collection": c.name,
		"filter":     filter,
	}

	for _, o := range options {
		o.Apply(body)
	}

	res := struct {
		Document json.RawMessage
	}{}
	err := c.database.do(ctx, ActionFindOne, body, &res)
	if err != nil {
		return FindOneResponse{Error: err}
	}
	if res.Document == nil {
		return FindOneResponse{Error: ErrNoDocuments}
	}

	return FindOneResponse{Document: res.Document}
}
