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

func createMedicines(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var NewMed Medicine

	err := decoder.Decode(&NewMed)

	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	log.Println(NewMed)

	db, err := sql.Open("sqlite3", "./medicieDB.db")
	defer db.Close()

	insertCmd := fmt.Sprintf(`insert into medicines(name, dose, remaining) values("%s", %s, %s)`, NewMed.Name, NewMed.Dose, NewMed.Remaining)

	_, err = db.Exec(insertCmd)

	if err != nil {
		log.Fatal(err)
	}

}

func retriveMedicines(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./medicieDB.db")
	defer db.Close()
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
	myRouter.HandleFunc("/medicines", retriveMedicines).Methods("GET")
	myRouter.HandleFunc("/medicines", createMedicines).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("server running at localhost:8080")
	handleRequests()
}
