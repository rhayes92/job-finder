package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

func loadDB() {

	r, err := db.Exec(`create table jobs (
	Job_ID varchar,
	Agency varchar,
	Posting_Type varchar,
	Business_Title varchar,
	Civil_Service_Title varchar,
	Title_Code varchar,
	Job_Level varchar,
	Job_Category varchar,
	Full_Part_Time_indicator varchar,
	Salary_Range_Begin float,
	Salary_Range_End float,
	Salary_Frequency varchar,
	Work_Location varchar,
	Division_Unit varchar,
	Job_Description	varchar,
	Minimum_Qual_Requirements varchar,
	Preffered_Skills varchar,
	Posting timestamp,
	PRIMARY KEY(Job_ID)
	);`)
	fmt.Println(r)
	if err != nil {
		fmt.Println(err)
	}

	csvFile, _ := os.Open("NYC_Jobs.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		baseStatement := "insert into public.jobs "
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		cols := "(Job_ID "
		vals := "('" + line[0] + "'"
		if line[1] != "" {
			cols = cols + ", Agency"
			vals = vals + ", '" + line[1] + "'"
		}
		if line[2] != "" {
			cols = cols + ", Posting_Type"
			vals = vals + ", '" + line[2] + "'"
		}
		if line[3] != "" {
			cols = cols + ", Business_Title"
			vals = vals + ", '" + line[3] + "'"
		}
		if line[4] != "" {
			cols = cols + ", Civil_Service_Title"
			vals = vals + ", '" + line[4] + "'"
		}
		if line[5] != "" {
			cols = cols + ", Title_Code"
			vals = vals + ", '" + line[5] + "'"
		}
		if line[6] != "" {
			cols = cols + ", Job_Level"
			vals = vals + ", '" + line[6] + "'"
		}
		if line[7] != "" {
			cols = cols + ", Job_Category"
			vals = vals + ", '" + line[7] + "'"
		}
		if line[8] != "" {
			cols = cols + ", Full_Part_Time_indicator"
			vals = vals + ", '" + line[8] + "'"
		}
		if line[9] != "" {
			cols = cols + ", Salary_Range_Begin"
			vals = vals + ", " + line[9]
		}
		if line[10] != "" {
			cols = cols + ", Salary_Range_End"
			vals = vals + ", " + line[10]
		}
		if line[11] != "" {
			cols = cols + ", Salary_Frequency"
			vals = vals + ", '" + line[11] + "'"
		}
		if line[12] != "" {
			cols = cols + ", Work_Location"
			vals = vals + ", '" + line[12] + "'"
		}
		if line[13] != "" {
			cols = cols + ", Division_Unit"
			vals = vals + ", '" + line[13] + "'"
		}
		if line[14] != "" {
			cols = cols + ", Job_Description"
			vals = vals + ", '" + line[14] + "'"
		}
		if line[15] != "" {
			cols = cols + ", Minimum_Qual_Requirements"
			vals = vals + ", '" + line[15] + "'"
		}
		if line[16] != "" {
			cols = cols + ", Preffered_Skills"
			vals = vals + ", '" + line[16] + "'"
		}
		if line[17] != "" {
			cols = cols + ", Posting"
			vals = vals + ", '" + line[17] + "'"
		}
		cols = cols + ")"
		vals = vals + ")"
		statement := baseStatement + cols + " VALUES " + vals
		db.Exec(statement)
	}

}
func main() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	loadDB()

	fmt.Println("Connected to database")
	r := mux.NewRouter()
	r.PathPrefix("/jobs/").Handler(http.StripPrefix("/jobs/", http.FileServer(http.Dir("./jobs"))))
	fmt.Println("Started server")
	http.ListenAndServe(":8080", r)
}
