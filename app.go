package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"hybrid-kvstore/badgerkv"
	"log"
	"net/http"
	"time"
)

type XRouter struct {
	Badger *badgerkv.BadgerStorage
}

func main() {
	xRouter := XRouter{
		Badger: badgerkv.Init(),
	}

	r := mux.NewRouter()
	r.HandleFunc("/", xRouter.handler)
	r.HandleFunc("/item", xRouter.newItem).Methods("POST")
	r.HandleFunc("/item/volatile", xRouter.getVolatileList).Methods("GET")
	r.HandleFunc("/item/persistence", xRouter.getPersistenceList).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func (x *XRouter) handler(w http.ResponseWriter, r *http.Request) {

	return
}

func (x *XRouter) newItem(w http.ResponseWriter, r *http.Request) {
	formData := badgerkv.MigratedCsvFile{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Printf("error parse form data : %+v", err)
		return
	}
	index := time.Now().Format("m_02_01_2006_15_04_05")
	x.Badger.AddNewItem(index, formData)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
	return
}

func (x *XRouter) getVolatileList(w http.ResponseWriter, r *http.Request) {
	data := x.Badger.GetVolatileItemList()
	w.Header().Set("Content-Type", "application/json")
	byteResult, err := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byteResult)
	if err != nil {
		log.Printf("[writeErrorAuthResponse] error write response buffer  : %#v", err)
	}
	return
}

func (x *XRouter) getPersistenceList(w http.ResponseWriter, r *http.Request) {
	data := x.Badger.GetPersistenceItemList()
	w.Header().Set("Content-Type", "application/json")
	byteResult, err := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byteResult)
	if err != nil {
		log.Printf("[writeErrorAuthResponse] error write response buffer  : %#v", err)
	}
	return
}
