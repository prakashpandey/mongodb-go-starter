package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string
	Age   int
	Phone string
	City  string
}

func (u *User) String() string {
	return fmt.Sprintf("Name[%s], Age[%d], Phone[%s], City[%s]", u.Name, u.Age, u.Phone, u.City)
}

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	dbName := "pcpdb"
	connectionString := "mongodb://localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	Database = client.Database(dbName)
	createUniqueIndexOnPCPCollection()
	u := &User{
		Name:  "pcp",
		Phone: "+919403678957",
		Age:   25,
		City:  "Pune",
	}
	// insert(u)
	update(u)
	// if uFromDB, err := FindUserByUsername(u.Name); err != nil {
	// 	log.Fatalf("error finding user[%s]", err.Error())
	// } else {
	// 	fmt.Printf("user retrieved is: %s", uFromDB.String())
	// }
}
