package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

const mongoURI string = "mongodb://localhost:27017/"

var mongoClient *mongo.Client

type Greeting struct {
	Message	 string `json:"Message"`
	Name     string `json:"Name"`
	Language string `json:"Language"`
}

type PrevGreetings struct {
	Count     int64      `json:"Count"`
	Greetings []Greeting `json:"Greetings"`
}

type GreetRequest struct {
	Name     string `json:"Name"`
	Language string `json:"Language"`
}

var greetings map[string]string = map[string]string{
	"EN" : "Hello",
	"ES" : "Hola",
	"LA" : "Salve",
	"PL" : "Cześć",
}


func main() {
	fmt.Println("Starting backend...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to MongoDB")


	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/getGreetings", getPrevGreetReq)
	router.HandleFunc("/greet", greetReq)

	log.Fatal(http.ListenAndServe(":1337",
		handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))

}

func getPrevGreetReq(w http.ResponseWriter, r *http.Request) {

	// Return value
	var pg PrevGreetings

	// get collection
	collection := mongoClient.Database("stackdemo").Collection("greetings")


	// set up context for call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()


	// get doc count
	i, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	pg.Count = i

	// get cursor for collection, no filter
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}

	// for each item, add it to the return value
	for cur.Next(ctx) {
		var g Greeting
		err := cur.Decode(&g)
		if err != nil {
			panic(err)
		}
		pg.Greetings = append(pg.Greetings, g)
	}

	json.NewEncoder(w).Encode(pg)

}

func greetReq(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("BAD REQ")
		panic(err)
	}

	var greq GreetRequest

	json.Unmarshal(reqBody, &greq)

	var g Greeting

	g.Name = greq.Name
	g.Language = greq.Language

	g.Message = fmt.Sprintf("%s %s!", greetings[g.Language], g.Name)

	collection := mongoClient.Database("stackdemo").Collection("greetings")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.InsertOne(ctx, g)

	json.NewEncoder(w).Encode(g)
}



