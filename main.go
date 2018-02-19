package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

type Medicine struct {
	Name      string `json:"Name"`
	Dose      string `json:"Dose"`
	Remaining string `json:"Remaining"`
}

func retriveMedicines(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./medicieDB.db")
	rows, err := db.Query("select name, dose, remaining from medicines")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var medicines []Medicine

	for rows.Next() {
		var name string
		var dose string
		var remaining string
		err = rows.Scan(&name, &dose, &remaining)
		if err != nil {
			log.Fatal(err)
		}
		medicines = append(medicines, Medicine{Name: name, Dose: dose, Remaining: remaining})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(medicines)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/medicines", retriveMedicines)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("server running at localhost:8080")
	handleRequests()
}
