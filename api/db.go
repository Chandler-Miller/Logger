package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"logger/config"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func addLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("Incorrect method used: %s. Please send the request with a POST method.", r.Method), http.StatusMethodNotAllowed)
	}

	dbName := r.URL.Query().Get("dbname")
	colName := r.URL.Query().Get("collection")
	mes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding log to server %s", err), http.StatusBadRequest)
		return
	}
	log.Println("Received log: ", mes)

	clientOptions := options.Client().ApplyURI(config.DBAddress)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to DB: %s", err), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database(dbName).Collection(colName)
	_, err = collection.InsertOne(context.Background(), bson.M{"message: ": mes})
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to add log to collection: %s", err), http.StatusBadRequest)
	}

	w.WriteHeader(200)
}

func getLogHandler() {
	
}
