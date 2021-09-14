package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type Greeting struct {
	Message	 string `json:"Message"`
	Name     string `json:"Name"`
	Language string `json:"Language"`
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

	handleRequests()
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", worldReq)
	router.HandleFunc("/greet", greetReq)

	//log.Fatal(http.ListenAndServe(":1337", router))

	log.Fatal(http.ListenAndServe(":1337",
		handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))

}

func worldReq(w http.ResponseWriter, r *http.Request) {
	g := Greeting{"Hello world!", "world", "EN"}

	json.NewEncoder(w).Encode(g)

	fmt.Println("Hello World REQ recieved")
}

func greetReq(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("BAD REQ")
		return
	}

	var greq GreetRequest

	json.Unmarshal(reqBody, &greq)

	var g Greeting

	g.Name = greq.Name
	g.Language = greq.Language

	g.Message = fmt.Sprintf("%s %s!", greetings[g.Language], g.Name)

	json.NewEncoder(w).Encode(g)

}

