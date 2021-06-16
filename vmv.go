package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SlotBooker struct {
	Name string
	Id   string
}

type Timeslot struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Location struct {
	Title     string     `json:"title"`
	Adress    string     `json:"adress"`
	Id        int        `json:"id"`
	Timeslots []Timeslot `json:"timeslots"`
}

type Error struct {
	Body string `json:"error"`
}

var db mongoDb

func homePage(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: homePage")
	fmt.Fprintf(writer, "Bienvenue sur l'API de ViteMonVaccin!\n")
	fmt.Fprintf(writer, "Voici les differents endpoints:\n")
	fmt.Fprintf(writer, "http://localhost:10000/book_apointment\n")
	fmt.Fprintf(writer, "http://localhost:10000/get_availabilities?id=0\n")
	fmt.Fprintf(writer, "http://localhost:10000/get_locations\n")
}

func getLocations(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: get_locations")
	json.NewEncoder(writer).Encode(db.getAllVaccinationCenters())
}

func getAvailabilities(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: get_availabilities")
	ids := request.URL.Query()["id"]

	if len(ids) > 1 {
		json.NewEncoder(writer).Encode(Error{Body: "Too many Location IDs"})
		return
	} else if len(ids) == 0 {
		json.NewEncoder(writer).Encode(Error{Body: "Location ID is mandatory"})
		return
	}

	timeslots := db.getTimeslots(ids[0], true)

	json.NewEncoder(writer).Encode(timeslots)
}

func bookAppointment(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: bookAppointment, Method: ", request.Method)
	if request.Method != "POST" {
		json.NewEncoder(writer).Encode(Error{Body: "bookAppointment is only open to POST Method"})
		return
	}
	var slotBooker SlotBooker

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(request.Body).Decode(&slotBooker)
	if err != nil {
		json.NewEncoder(writer).Encode(err)
		return
	}

	err = db.reserveTimeslot(slotBooker.Id, slotBooker.Name)
	if err != nil {
		json.NewEncoder(writer).Encode(err.Error())
		return
	}
}

func handleRequests() {
	http.HandleFunc("/book_apointment", bookAppointment)
	http.HandleFunc("/get_timeslots", getAvailabilities)
	http.HandleFunc("/get_locations", getLocations)
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {

	client, context := initDb()
	db = mongoDb{client, context}

	handleRequests()
}
