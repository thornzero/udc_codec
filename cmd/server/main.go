package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/thornzero/udc_codec/pkg/db"
)

var store *db.Store

func main() {
	s, err := db.OpenDB("tags.db")
	if err != nil {
		panic(err)
	}
	if err := s.Migrate(); err != nil {
		panic(err)
	}
	store = s

	router := mux.NewRouter()
	router.HandleFunc("/tags/{tag}", GetTag).Methods("GET")
	router.HandleFunc("/tags", InsertTag).Methods("POST")

	fmt.Println("API server running on :8080")
	http.ListenAndServe(":8080", router)
}

func GetTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tag := vars["tag"]
	record, err := store.LookupTag(tag)
	if err != nil {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(record)
}

func InsertTag(w http.ResponseWriter, r *http.Request) {
	var t db.TagRecord
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := store.InsertTag(&t); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(t)
}
