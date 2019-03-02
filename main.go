package main

import (
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"database/sql"
 	 "fmt"

)
var db *sql.DB
const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "admin"
  dbname   = "postgres"
)

func main(){
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
   		 "password=%s dbname=%s sslmode=disable",
   	 host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
	  panic(err)
	}
	defer db.Close()
	fmt.Println("Connected to database")
	r := mux.NewRouter()
	r.PathPrefix("/jobs/").Handler(http.StripPrefix("/jobs/", http.FileServer(http.Dir("./jobs"))))
	fmt.Println("Started server")
	http.ListenAndServe(":8080",r)
}
