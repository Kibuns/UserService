package DAL

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Kibuns/UserService/Models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// global variable mongodb connection client
var client mongo.Client = newClient()

// ----Create----
func RegisterUser(user Models.User) error {
	result, err := usernameIsUnique(user.UserName)
	if err != nil {
		panic(err)
	}
	if !result {
		return errors.New("username is not unique")
	}

	userCollection := client.Database("UserDB").Collection("users")
	user.Created = time.Now()
	user.UserID = uuid.New().String()
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	fmt.Println("New user added called: " + user.UserName)
	return nil
}

//----Read----

func ReadAllUsers() ([]primitive.M, error) {
	twootCollection := client.Database("UserDB").Collection("users")
	// retrieve all the documents (empty filter)
	cursor, err := twootCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}



func SearchUser(id string) (value primitive.M) {
	twootCollection := client.Database("TwootDB").Collection("twoots")
	// convert the hexadecimal string to an ObjectID type
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	// retrieve the document with the specified _id
	var result bson.M
	err = twootCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}}).Decode(&result)
	if err != nil {
		panic(err)
	}

	// display the retrieved document
	fmt.Println("displaying the result from the search query")
	fmt.Println(result)
	value = result

	return value
}


//----Update----

//----Delete----

// other
func newClient() (value mongo.Client) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://ninoverhaegh:6P77TACMZwsd8pb4@twotterdb.jfx1rk2.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	value = *client

	return
}

func usernameIsUnique(username string) (bool, error) {
	users, err := ReadAllUsers()
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user["username"] == username {
			return false, nil
		}
	}

	return true, nil
}
