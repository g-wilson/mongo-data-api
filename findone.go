package mongoapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrNoDocuments = errors.New("mongoapi: no documents found in result")
)

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

func (c *Collection) FindOne(ctx context.Context, filter bson.M) FindOneResponse {
	body := bson.M{
		"collection": c.name,
		"filter":     filter,
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

func (c *Collection) FindOneWithProjection(ctx context.Context, filter, projection bson.M) FindOneResponse {
	body := bson.M{
		"collection": c.name,
		"filter":     filter,
		"projection": projection,
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
