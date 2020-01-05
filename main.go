package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	d := mux.NewRouter()
	f := mux.NewRouter()
	d.HandleFunc("/statusChange", func(w http.ResponseWriter, r *http.Request) {
		changeStatus(w, r)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		t, err := template.ParseFiles("templates/index.html")
		// f, err := os.Open("templates/index.html")
		// stat, err := f.Stat()
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// buf := make([]byte, stat.Size())

		// f.Read(buf)
		// fmt.Printf("%s\n", buf)

		// t, err := template.New("test").Parse(string(buf))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("r: %+v\n w: %+v\n", r, w)
		// for index, _ := range laundryMachines {
		var a RoomTotal
		a.RoomNum = "NineTenApt"
		a.Washer = 1
		a.Dryer = 0
		laundryRoomContent = append(laundryRoomContent, a)
		err = t.Execute(w, laundryRoomContent)
		if err != nil {
			panic(err)
		}
		// }
	})
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir("templates/css/"))))

	//Must make a post request with a JSON in the body containing {"name":"{UUID}"}
	f.HandleFunc("/machineID", func(w http.ResponseWriter, r *http.Request) {
		i := getStatus(w, r)
		fmt.Println(i)
		if i == 0 {
			http.Error(w, "laundry machine not recognized", http.StatusBadRequest)
			return
		} else {
			fmt.Fprintf(w, strconv.Itoa(i))

		}

	})
	http.Handle("/", r)
	http.Handle("/machineID", f)
	http.Handle("/statusChange", d)
	log.Fatal(http.ListenAndServe(":8081", nil))

}
