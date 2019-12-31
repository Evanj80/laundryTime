package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Machine struct {
	RoomNum     string
	MachineType string
	Status      int
}

type RoomTotal struct {
	RoomNum string
	Washer  int
	Dryer   int
}

var laundryMachines []Machine
var laundryRoomContent []RoomTotal

func changeStatus(w http.ResponseWriter, r *http.Request) {

	var l Machine

	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(l)
	fmt.Print(l.Status)
	counter := 0
	for index, _ := range laundryMachines {
		if laundryMachines[index].MachineType == l.MachineType && laundryMachines[index].RoomNum == l.RoomNum {
			laundryMachines[index].Status = l.Status
			counter++
		}
	}
	if counter == 0 {
		laundryMachines = append(laundryMachines, l)
	}
	counter = 0
	if l.MachineType == "Washer" {
		for index, _ := range laundryRoomContent {
			if laundryRoomContent[index].RoomNum == l.RoomNum {
				laundryRoomContent[index].Washer = laundryRoomContent[index].Washer + l.Status
				counter++
			}
		}
		if counter == 0 {
			var x RoomTotal
			x.RoomNum = l.RoomNum
			x.Washer = 1
			x.Dryer = 0
			laundryRoomContent = append(laundryRoomContent, x)
		}

	}
	if l.MachineType == "Dryer" {
		for index, _ := range laundryRoomContent {
			if laundryRoomContent[index].RoomNum == l.RoomNum {
				laundryRoomContent[index].Dryer = laundryRoomContent[index].Dryer + l.Status
				counter++
			}
		}
		if counter == 0 {
			var x RoomTotal
			x.RoomNum = l.RoomNum
			x.Washer = 0
			x.Dryer = 1
			laundryRoomContent = append(laundryRoomContent, x)
		}

	}
	fmt.Println(laundryMachines)
	fmt.Println(laundryRoomContent)
}

func getStatus(w http.ResponseWriter, r *http.Request) int {
	var l Machine
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return -1
	}
	found := 0
	for index, v := range laundryMachines {
		if v.MachineType == l.MachineType {
			found = index
		}
	}
	var i = laundryMachines[found]
	fmt.Println(i)
	print(i.Status)
	return i.Status
}

func main() {
	//Must make a post request with a JSON in the body Update status to either 1 or -1
	//Example Curl
	//curl -d '{"roomnum" : "collegenine", "MachineType":"Dryer", "status" : -1}' -H "Content-Type: application/json" -X POST http://localhost:8081/statusChange

	http.HandleFunc("/statusChange", func(w http.ResponseWriter, r *http.Request) {
		changeStatus(w, r)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// for index, _ := range laundryMachines {
		err = t.Execute(w, laundryRoomContent)
		if err != nil {
			panic(err)
		}
		// }
	})

	//Must make a post request with a JSON in the body containing {"name":"{UUID}"}
	http.HandleFunc("/machineID", func(w http.ResponseWriter, r *http.Request) {
		i := getStatus(w, r)
		fmt.Println(i)
		if i == 0 {
			http.Error(w, "laundry machine not recognized", http.StatusBadRequest)
			return
		} else {
			fmt.Fprintf(w, strconv.Itoa(i))

		}

	})
	log.Fatal(http.ListenAndServe(":8081", nil))

}
