package main

import (
	"context"
	"os"
	"time"

	mongoapi "github.com/g-wilson/mongo-data-api"

	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	ID        string    `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	mClient := mongoapi.New(os.Getenv("MONGODB_API_URL"), os.Getenv("MONGODB_API_KEY"))
	db := mClient.Database(os.Getenv("MONGODB_CLUSTER_NAME"), os.Getenv("MONGODB_DB_NAME"))

	ctx := context.Background()

	user := User{}
	err = db.Collection("users").
		FindOne(ctx, bson.M{"_id": "1"}).
		Decode(&user)
	if err != nil {
		panic(err)
	}

	pretty.Println(user)

	allUsers := []User{}
	err = db.Collection("users").
		Find(ctx, bson.M{}, mongoapi.NewFindOptions().
			WithLimit(2).
			WithSort(bson.D{{"created_at", -1}}),
		).
		All(&allUsers)
	if err != nil {
		panic(err)
	}

	pretty.Println(allUsers)
}
