package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection URI
var (
	USERNAME = os.Getenv("PERCONASERVERMONGODB_MONGODB_USER_ADMIN_USER")
	PASSWORD = os.Getenv("PERCONASERVERMONGODB_MONGODB_USER_ADMIN_PASSWORD")
	HOST     = os.Getenv("PERCONASERVERMONGODB_HOST")
	uri      = fmt.Sprintf("mongodb://%s:%s@%s:27017/?maxPoolSize=20&w=majority", USERNAME, PASSWORD, HOST)
)

func main() {
	http.HandleFunc("/", Ping)
	http.HandleFunc("/connect", ConnectToMongoDB)
	http.ListenAndServe("0.0.0.0:8080", nil)

}

func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
	// fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func ConnectToMongoDB(w http.ResponseWriter, r *http.Request) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Fprintf(w, "failed to connect to the server: %s", err.Error())
		return
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			fmt.Fprintf(w, "failed to disconnect: %s", err.Error())
			return
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Fprintf(w, "unable to connect to the server: %s", err.Error())
		return
	}
	fmt.Fprint(w, "Successfully pinged the server.")
}
