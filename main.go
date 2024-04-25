package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {

}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API health is okay")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{ "alive": true}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	log.Info("Starting API server.")
	router := mux.NewRouter()

	router.HandleFunc("/Healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	http.ListenAndServe(":8000", router)
}
