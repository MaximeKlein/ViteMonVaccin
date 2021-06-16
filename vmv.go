package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

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

var locations []Location
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
	json.NewEncoder(writer).Encode(locations)
}

func getLocation(id int) (*Location, error) {
	for _, loc := range locations {
		if loc.Id == id {
			return &loc, nil
		}
	}
	return nil, fmt.Errorf("couldnt find locations %d in database", id)
}

func getAvailabilities(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: get_availabilities")
	ids := request.URL.Query()["id"]

	if len(ids) > 1 {
		fmt.Fprintf(writer, "Too many Ids")
		return
	} else if len(ids) == 0 {
		fmt.Fprintf(writer, "Location ID is mandatory")
		return
	}

	id, error := strconv.Atoi(ids[0])

	if error != nil {
		log.Fatal(error)
		return
	}

	location, error := getLocation(id)

	if error != nil {
		fmt.Println(error)
		json.NewEncoder(writer).Encode(Error{Body: error.Error()})
		return
	}

	json.NewEncoder(writer).Encode(location.Timeslots)
}

func bookAppointment(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: bookAppointment, Method: ", request.Method)
	if request.Method != "POST" {
		json.NewEncoder(writer).Encode(Error{Body: "bookAppointment is only open to POST Method"})
		return
	}

}

func handleRequests() {
	http.HandleFunc("/book_apointment", bookAppointment)
	http.HandleFunc("/get_availabilities", getAvailabilities)
	http.HandleFunc("/get_locations", getLocations)
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	locations = []Location{
		Location{Title: "Centre de vaccination", Adress: "1 Rue du docteur Maboul", Id: 0, Timeslots: []Timeslot{Timeslot{Date: "10/06/2021", Time: "10:30"}, Timeslot{Date: "10/06/2021", Time: "10:40"}}},
		Location{Title: "Pharmacie du coin", Adress: "1 place du patient anglais", Id: 1, Timeslots: []Timeslot{Timeslot{Date: "10/06/2021", Time: "10:30"}, Timeslot{Date: "10/06/2021", Time: "10:40"}}},
	}

	handleRequests()
}
