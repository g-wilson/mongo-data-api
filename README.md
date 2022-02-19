# Mongo Data API Go client

Go client for the [Mongo Data API](https://docs.atlas.mongodb.com/api/data-api/).

Designed to closely resemble the official Go client for Mongo, so that migrating or sharing code between the two clients is as simple as possible.

- Request parameters are defined with bson primitives
- BSON tags used in decode methods when unmarshalling responses

## Warning

This is **extremely work-in-progress** - right now only the `findOne` action is implemented!

## Why might you want this?

If...

- you write Lambda functions in Go (other serverless function services are available...)
- that connect to a Mongo database
- and you get a lot of traffic so lots of Lambas are running
- and you are worried about having too many database connections

## Example:

```go
type User struct {
	ID        string    `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
}

func main() {
	ctx := context.Background()

	mClient := mongoapi.New("https://data.mongodb-api.com/app/your-id/endpoint/data/beta", os.Getenv("MONGODB_API_KEY"))
	db := mClient.Database("cluster-name", "db-name")

	user := User{}
	err = db.Collection("users").
		FindOne(ctx, bson.M{"_id": "1"}).
		Decode(&user)
	if err != nil {
		panic(err)
	}

	fmt.Println(User.ID)
	// ...
}
```
