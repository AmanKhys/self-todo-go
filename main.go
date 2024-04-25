package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "root"
	dbName   = "todolist"
)

var db, _ = gorm.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

type TodoItemModel struct {
	Id          int    `gorm:"primary key"`
	Description string `gorm:"not null"`
	Completed   bool   `gorm:"not null"`
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new todo Item. Saving to database")
	todo := &TodoItemModel{Description: description, Completed: false}
	//	log.Println(todo)
	db.Create(&todo)
	//	log.Println(todo)
	result := db.Last(&todo)
	//	log.Println(todo)
	//	log.Println(result.Value)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)

}

func GetItemByID(id int) bool {
	todo := &TodoItemModel{}
	result := db.First(&todo, id)
	if result.Error != nil {
		log.Warn("TodoItem not found in the database")
		return false
	}
	return true
}
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	completed, _ := strconv.ParseBool(r.FormValue("completed"))
	err := GetItemByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated":false, "error":"Record not found"}`)
	} else {
		log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
		todo := &TodoItemModel{Id: id}
		db.First(&todo)
		todo.Completed = completed
		db.Save(&todo)

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated":true}`)
	}
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
	defer db.Close()
	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})
	log.Info("Starting API server.")
	router := mux.NewRouter()

	router.HandleFunc("/Healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateItem).Methods("PUT")
	http.ListenAndServe(":8000", router)
}
