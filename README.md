# Mongo Data API Go client

Go client for the [Mongo Data API](https://docs.atlas.mongodb.com/api/data-api/).

Designed to closely resemble the official Go client for Mongo, so that migrating or sharing code between the two clients is as simple as possible.

- Request parameters are defined with bson primitives
- BSON tags used in decode methods when unmarshalling responses

**Extremely work-in-progress**

Right now only `findOne` action is implemented!

### Example:

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
