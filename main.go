package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
		"fmt"
		"strings"
	"github.com/bradfitz/slice"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)


var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

type jobInfo struct {
	JobID            string
	BusTitle string
	SalaryRangeBegin float64
	SalaryRangeEnd   float64
	DivisionUnit     string
	JobCategory      string
}
type jobScoreInfo struct {
	Score float64
	JobID string
	BusTitle string
	SalaryRangeBegin float64
	SalaryRangeEnd   float64
	DivisionUnit     float64
	JobCategory      float64
}
type category struct {
	JobCategory string  `json:"category"`
	Rank        float64 `json:"rank"`
	Score       float64
}
type Weights struct {
	JobCategories float64  `json:"category"`
	Divisions float64  `json:"divisions"`
	End float64  `json:"end"`
	Begin float64  `json:"begin"`
}
type evalStruct struct {
	JobCategories []category `json:"categories"`
	Divisions     []category `json:"divisions"`
	Weight Weights `json:"weight"`
}
type catsStruct struct {
	JobCategories []string `json:"categories"`
	Divisions     []string `json:"divisions"`
}

var jobScore map[string]jobScoreInfo
var jobs []jobInfo

var jobCategories []string
var divisions []string
func getCat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request")
	var cat catsStruct
	cat.JobCategories = jobCategories
	cat.Divisions = divisions
	resp, _ := json.Marshal(cat)
	fmt.Fprint(w,string(resp))
}
func ScoreEval(w http.ResponseWriter, r *http.Request){
	body, _ := ioutil.ReadAll(r.Body)
		fmt.Println("tmep",string(body))
	var eval evalStruct
	 json.Unmarshal(body, &eval)
//	fmt.Println(eval)
	fmt.Println(eval.JobCategories)
	linearCat("division", eval.Divisions)
	linearCat("jobcat", eval.JobCategories)
	results := calcScore(eval.Weight)
	resp,_ := json.Marshal(results)
	fmt.Println("scoring")
	fmt.Fprint(w,string(resp))
}
func calcScore(w Weights) map[string]jobScoreInfo{
	var weight [][]float64
	weight = make([][]float64,4)
	for i :=0; i < 4; i++{
		weight[i] =make([]float64,4)
	}
	for  i := range weight{
		var total float64
		if w.JobCategories <= float64(i) +1.0{
			weight[i][0]= 1.0
		}
		if w.Divisions <= float64(i) +1.0{
			weight[i][1]= 1.0
		}
		if w.End <= float64(i) +1.0{
			weight[i][2]= 1.0
		}
		if w.Begin <= float64(i) +1.0{
			weight[i][3]= 1.0
		}
		for j := range weight[i]{
			total = total +  weight[i][j]
		}
		if total != 0{
		 weight[i][0]= weight[i][0]/total
		 weight[i][1]= weight[i][1]/total
		 weight[i][2]= weight[i][2]/total
		 weight[i][3]= weight[i][3]/total
	 }
		 fmt.Println("weight",weight)
	}
		w.JobCategories =  (weight[0][0] +  weight[1][0] +  weight[2][0] +  weight[3][0])/4
		w.Divisions =  (weight[0][1] +  weight[1][1] +  weight[2][1]+  weight[3][1])/4
		w.End =  (weight[0][2] +  weight[1][2] +  weight[2][2]+  weight[3][2])/4
		w.Begin =  (weight[0][3] +  weight[1][3] +  weight[2][3]+  weight[3][3])/4
		fmt.Println(w)
		tempArray := make([]jobScoreInfo,0)
		for key := range jobScore {
			job := jobScore[key]
			job.Score = job.SalaryRangeBegin * w.Begin + job.SalaryRangeEnd *  w.End  + job.DivisionUnit *	w.Divisions + job.JobCategory *	w.JobCategories
			job.JobID = key
			jobScore[key] = job
			tempArray = append(tempArray,job)
		}
		slice.Sort(tempArray[:], func(i, j int) bool {
			return tempArray[i].Score > tempArray[j].Score
		})
		 tempScore := make(map[string]jobScoreInfo)
		for j := range tempArray{
			tempScore[tempArray[j].JobID]	= tempArray[j]
			if j > 9 {
				break
			}
		}
		return tempScore
}
func isUnique( val string, arry []string) bool{
	for _,x := range arry{
		if val == x{
			return false
		}
	}
	return true
}
func loadDB() {
	jobCategories = make([]string,0)
	divisions = make([]string,0)
	jobScore = make(map[string]jobScoreInfo)
	jobs = make([]jobInfo, 0)
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
		//	panic("fuck")
	}

	csvFile, _ := os.Open("NYC_Jobs.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		var job jobInfo
		baseStatement := "insert into public.jobs "
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		cols := "(Job_ID "
		vals := "('" + line[0] + "'"
		job.JobID = line[0]
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
			job.BusTitle  =  line[3]
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
			//strings.SplitAfter("a,b,c", ",")
			line[7]= strings.Trim(line[7],`"`)
			line[7]= strings.Replace(line[7],"&","and",-1)
			job.JobCategory = line[7]
			if isUnique(line[7],jobCategories) {
				jobCategories = append(jobCategories,line[7])
			}
		}
		if line[8] != "" {
			cols = cols + ", Full_Part_Time_indicator"
			vals = vals + ", '" + line[8] + "'"
		}
		if line[9] != "" {
			cols = cols + ", Salary_Range_Begin"
			vals = vals + ", " + line[9]
			s, err := strconv.ParseFloat(line[9], 32)
			if err != nil {
				continue
			}
			if line[11] == "Hourly" {
				s = s * 52 * 40
			} else if line[11] == "Daily" {
				s = s * 52 * 5
			}
			job.SalaryRangeBegin = s
		}
		if line[10] != "" {
			cols = cols + ", Salary_Range_End"
			vals = vals + ", " + line[10]
			s, err := strconv.ParseFloat(line[10], 32)
			if err != nil {
				continue
			}
			if line[11] == "Hourly" {
				s = s * 52 * 40
			} else if line[11] == "Daily" {
				s = s * 52 * 5
			}
			job.SalaryRangeEnd = s
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
			line[13]= strings.Trim(line[13],`"`)
			line[13]= strings.Replace(line[13],"&","and",-1)
			job.DivisionUnit = line[13]
			if isUnique(line[13],divisions) {
				divisions = append(divisions,line[13])
			}
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
		r, err = db.Exec(statement)
		if err != nil {
			//	fmt.Println("error", err)
			continue
		}
		jobs = append(jobs, job)
	}
	fmt.Println(len(jobCategories),len(divisions))
}
func linearCat(typeOfVal string, cat []category) {
	for i := range cat {
		if cat[i].Rank == 0 {
			cat[i].Rank = 1
		}
		cat[i].Score = (cat[i].Rank - 1) / (7 - 1)
	}
	for _, job := range jobs {
		var score float64
		var tempScore jobScoreInfo
		var found bool
		if typeOfVal == "division" {
			for _, catVal := range cat {
				if job.DivisionUnit == catVal.JobCategory {
					found = true
					score = catVal.Score
				}
			}
			if !found{
				score = 0
			}
			tempScore = jobScore[job.JobID]
			tempScore.DivisionUnit = score
			jobScore[job.JobID] = tempScore
		} else {
			for _, catVal := range cat {
				if job.JobCategory == catVal.JobCategory {
					found = true
					score = catVal.Score
				}
			}
			if !found{
				score = 0
			}
			tempScore = jobScore[job.JobID]
			tempScore.JobCategory = score
			jobScore[job.JobID] = tempScore
		}

	}
}
func linearEnd() {
	slice.Sort(jobs[:], func(i, j int) bool {
		return jobs[i].SalaryRangeEnd > jobs[j].SalaryRangeEnd
	})
	//fmt.Println(jobs[0].SalaryRangeEnd, "--", jobs[len(jobs)-1].SalaryRangeEnd)
	for _, val := range jobs {
		var jobScoreVal jobScoreInfo
		jobScoreVal.BusTitle = val.BusTitle
		jobScoreVal.SalaryRangeEnd = (val.SalaryRangeEnd - jobs[len(jobs)-1].SalaryRangeEnd) / (jobs[0].SalaryRangeEnd - jobs[len(jobs)-1].SalaryRangeEnd)
		jobScore[val.JobID] = jobScoreVal
	}
}
func linearBegin() {
	slice.Sort(jobs[:], func(i, j int) bool {
		return jobs[i].SalaryRangeBegin > jobs[j].SalaryRangeBegin
	})
	//fmt.Println(jobs[0].SalaryRangeBegin, "--", jobs[len(jobs)-1].SalaryRangeBegin)
	for _, val := range jobs {
		var jobScoreVal jobScoreInfo
		jobScoreVal =	jobScore[val.JobID]
		jobScoreVal.SalaryRangeBegin = (val.SalaryRangeBegin - jobs[len(jobs)-1].SalaryRangeBegin) / (jobs[0].SalaryRangeBegin - jobs[len(jobs)-1].SalaryRangeBegin)
		jobScore[val.JobID] = jobScoreVal
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
	linearEnd()
	linearBegin()
	fmt.Println("Connected to database")
	r := mux.NewRouter()
	r.HandleFunc("/ScoreEval", ScoreEval).Methods("POST")
	r.HandleFunc("/cat", getCat).Methods("GET")
	r.PathPrefix("/jobs/").Handler(http.StripPrefix("/jobs/", http.FileServer(http.Dir("./jobs"))))
	fmt.Println("Started server")
	http.ListenAndServe(":8080", r)
}
