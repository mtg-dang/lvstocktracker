package main

import (
	"encoding/json"
	"example.com/lvapi"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnItemFamily(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Endpoint Hit: Item Family for SKU: " + vars["sku"])
	json.NewEncoder(w).Encode(lvapi.GetLVAlternativeStyleProductIndentifierAndAvailabilityForSKU(vars["sku"]))
}

func returnItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println("Endpoint Hit: Item Family for SKU: " + vars["sku"])
	json.NewEncoder(w).Encode(lvapi.GetLVProductAvailabilityBySKU(vars["sku"]))
}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/api/itemfamily/{sku}", returnItemFamily)
	r.HandleFunc("/api/item/{sku}", returnItem)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	handleRequests()
}
