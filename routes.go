package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Routes API object
type Routes struct {
	Mux     *mux.Router
	Negroni *negroni.Negroni
	Store   *Store
}

// NewRoutes create API backend for the program
func NewRoutes(store *Store) Routes {
	r := Routes{}
	r.Store = store
	return r
}

// getVersion handler
func getVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(version))
}

func (r *Routes) getAll(w http.ResponseWriter, req *http.Request) {
	docs, err := r.Store.GetAll()
	if err != nil {
		log.Println("Error retrieving documents")
		http.Error(w, "Error retrieving documents", http.StatusInternalServerError)
	}

	docsByte, err := json.Marshal(docs)
	if err != nil {
		log.Println("Error writing documents.")
		http.Error(w, "Error writing documents", http.StatusInternalServerError)
	}
	w.Write(docsByte)
}

func (r *Routes) postDoc(w http.ResponseWriter, req *http.Request) {
	name, ok := req.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		http.Error(w, "Error reading name param", http.StatusBadRequest)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)

	}

	// tags, ok := r.URL.Query()["tags"]
	r.Store.Add(name[0], []string{}, body)
}

func (r *Routes) delDoc(w http.ResponseWriter, req *http.Request) {
	id, ok := req.URL.Query()["id"]

	if !ok || len(id[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading id param."))
	}

	err := r.Store.Delete(id[0])

	if err != nil {
		log.Printf("Error %v deleting id: %s", err, id[0])
		http.Error(w, "Error reading body", http.StatusBadRequest)
	}
}

// Run API function
func (r *Routes) Run(addr string) {

	// Define mux
	r.Mux = mux.NewRouter()

	// API handler
	r.Mux.HandleFunc("/api/version", getVersion).Methods("GET")
	r.Mux.HandleFunc("/api/docs", r.getAll).Methods("GET")
	r.Mux.HandleFunc("/api/doc", r.postDoc).Methods("POST")
	r.Mux.HandleFunc("/api/doc", r.delDoc).Methods("DELETE")

	// Define negroni middleware
	r.Negroni = negroni.New()
	r.Negroni.UseHandler(r.Mux)

	log.Fatal(http.ListenAndServe(addr, r.Negroni))
}
