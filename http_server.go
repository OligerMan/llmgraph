package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewHTTPServer(addr string) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/instance_capability", setInstanceCapabilities).Methods("POST")
	r.HandleFunc("/delete_instance", deleteInstance).Methods("POST")
	r.HandleFunc("/execution_graph", setExecutionGraph).Methods("POST")
	r.HandleFunc("/delete_execution_graph", deleteExecutionGraph).Methods("POST")

	r.HandleFunc("/execute", execute).Methods("POST")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
