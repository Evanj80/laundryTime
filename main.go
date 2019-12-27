package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

type RoomID struct {
	machines []Machine
}

type Machine struct {
	Name   string
	Status int
}

var laundryMachines []Machine

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
		if laundryMachines[index].Name == l.Name {
			laundryMachines[index].Status = l.Status
			counter++
		}
	}
	if counter == 0 {
		laundryMachines = append(laundryMachines, l)
	}
	fmt.Println(laundryMachines)

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
		if v.Name == l.Name {
			found = index
		}
	}
	var i = laundryMachines[found]
	fmt.Println(i)
	print(i.Status)
	return i.Status
}

func main() {
	//Must make a post request with a JSON in the body containing {"name":"{UUID}","status": {newstatus}} Update status to either 1 or 2
	http.HandleFunc("/statusChange", func(w http.ResponseWriter, r *http.Request) {
		changeStatus(w, r)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fp := path.Join("templates", "index.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for index, _ := range laundryMachines {
			tmpl.Execute(w, laundryMachines[index])
		}
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
