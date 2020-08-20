package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/knowledge/helper"
	"github.com/knowledge/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getKnowledges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// we created  array
	var knowledges []models.Knowledge

	//Connection mongoDB with helper class
	collection := helper.ConnectDB()

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var knowledge models.Knowledge
		// & character returns the memory address of the following variable.
		err := cur.Decode(&knowledge) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		knowledges = append(knowledges, knowledge)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(knowledges) // encode similar to serialize process.
}

func getKnowledge(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var knowledge models.Knowledge
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := helper.ConnectDB()

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&knowledge)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(knowledge)
}

func createKnowledge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var knowledge models.Knowledge

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&knowledge)

	// connect db
	collection := helper.ConnectDB()

	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), knowledge)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func updateKnowledge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var knowledge models.Knowledge

	collection := helper.ConnectDB()

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&knowledge)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"title", knowledge.Title},
			{"detail", knowledge.Detail},
			{"viewcount", knowledge.Viewcount},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&knowledge)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	knowledge.ID = id

	json.NewEncoder(w).Encode(knowledge)
}

func deleteKnowledge(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	collection := helper.ConnectDB()

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

// var client *mongo.Client

func main() {
	//Init Router
	r := mux.NewRouter()

	//r.HandleFunc("/api/knowledge", getKnowledges).Methods("GET")
	r.HandleFunc("/api/knowledge/{id}", getKnowledge).Methods("GET")
	//r.HandleFunc("/api/knowledge", createKnowledge).Methods("POST")
	//r.HandleFunc("/api/knowledge/{id}", updateKnowledge).Methods("PUT")
	//r.HandleFunc("/api/knowledge/{id}", deleteKnowledge).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3404", r))

}
