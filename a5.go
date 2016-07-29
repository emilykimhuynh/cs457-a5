/******************************************************
*
*       Emily Huynh
*       CS 457: Assignment 5
*
*       Amusement Parks Database
*       Tools: go & sqlite3
*
******************************************************/

package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

var validPath = regexp.MustCompile("^/(query|display)/([a-zA-Z0-9]+)$")
var db = openDatabase()

//Miniworld collects all relations in a global struct
var Miniworld = struct {
	ParkCollection             map[int]AmusementPark
	AttractionCollection       map[string]Attraction
	LocationCollection         map[string]Location
	VisitorCollection          map[int]Visitor
	EmployeeCollection         map[int]Employee
	GoesToCollection           map[string]GoesTo
	InteractsWithCollection    map[string]InteractsWith
	ParkNameAndAddressView     map[string]ViewParkNameAndAddress
	VisitorNameAndParkNameView map[string]ViewVisitorNameAndParkName
	QueryStmt                  string
}{
	ParkCollection:             make(map[int]AmusementPark),
	AttractionCollection:       make(map[string]Attraction),
	LocationCollection:         make(map[string]Location),
	VisitorCollection:          make(map[int]Visitor),
	EmployeeCollection:         make(map[int]Employee),
	GoesToCollection:           make(map[string]GoesTo),        //key is "ParkID" + "VisitorTicketID"
	InteractsWithCollection:    make(map[string]InteractsWith), //"VisitorTicketID" + "AttractionName"
	ParkNameAndAddressView:     make(map[string]ViewParkNameAndAddress),
	VisitorNameAndParkNameView: make(map[string]ViewVisitorNameAndParkName),
}

/**********     types for miniworld relations     **********/

//AmusementPark relation representation
type AmusementPark struct {
	Name string
	ID   int //PK
}

//Attraction relation representation
type Attraction struct {
	Name   string //PK
	MgrID  int
	ParkID int
}

//Location relation representation
type Location struct {
	Address string //PK
	ParkID  int
}

//Visitor relation representation
type Visitor struct {
	FastPass bool
	TicketID int //PK
	Name     string
}

//Employee relation representation
type Employee struct {
	Fname                  string
	Minit                  string
	Lname                  string
	ID                     int //PK
	BirthDate              string
	Salary                 string
	HireDate               string
	Address                string
	SupervisorID           sql.NullInt64
	AttractionName         string
	AttractionDateAssigned string
}

//InteractsWith relation representation
type InteractsWith struct {
	VisitorTicketID int    //PK
	AttractionName  string //PK
}

//GoesTo relation representation
type GoesTo struct {
	ParkID          int //PK
	VisitorTicketID int //PK
}

//ViewParkNameAndAddress view representation
type ViewParkNameAndAddress struct {
	ParkName string
	Address  string
}

//ViewVisitorNameAndParkName view representation
type ViewVisitorNameAndParkName struct {
	ParkName    string
	VisitorName string
}

//QueryResultRow for select queries
type QueryResultRow struct {
	VisitorName    string
	ParkName       string
	FirstName      string
	LastName       string
	AttractionName string
}

/**********     DB functions     **********/

//opens the database
func openDatabase() *sql.DB {
	//initialize the database
	db, err := sql.Open("sqlite3", "parks_miniworld.db")

	checkErr(err)
	return db
}

//maps db query results to global struct
func populateMaps(db *sql.DB) {

	//Amusement Parks
	rows, err := db.Query("select Name, ID from Amusement_Park")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r AmusementPark

		err = rows.Scan(&r.Name, &r.ID)
		checkErr(err)

		Miniworld.ParkCollection[r.ID] = r
	}

	//Attractions
	rows, err = db.Query("select Name, Mgr_ID, Park_ID from Attraction")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r Attraction

		err = rows.Scan(&r.Name, &r.MgrID, &r.ParkID)
		checkErr(err)

		Miniworld.AttractionCollection[r.Name] = r
	}

	//Locations
	rows, err = db.Query("select Address, Park_ID from Location")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r Location

		err = rows.Scan(&r.Address, &r.ParkID)
		checkErr(err)

		Miniworld.LocationCollection[r.Address] = r
	}

	//Visitors
	rows, err = db.Query("select Name, Fast_Pass, Ticket_ID from Visitor")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r Visitor

		err = rows.Scan(&r.Name, &r.FastPass, &r.TicketID)
		checkErr(err)

		Miniworld.VisitorCollection[r.TicketID] = r
	}

	//Employees
	rows, err = db.Query(`select Fname, Minit, Lname, ID, Birth_Date, Salary,
        Hire_Date, Address, Supervisor_ID, Attraction_Name,
        Attraction_Date_Assigned from Employee`)

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r Employee

		err = rows.Scan(&r.Fname, &r.Minit, &r.Lname, &r.ID, &r.BirthDate,
			&r.Salary, &r.HireDate, &r.Address, &r.SupervisorID,
			&r.AttractionName, &r.AttractionDateAssigned)
		checkErr(err)

		Miniworld.EmployeeCollection[r.ID] = r
	}

	//GoesTo
	rows, err = db.Query("select Park_ID, Visitor_Ticket_ID from Goes_To")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r GoesTo

		err = rows.Scan(&r.ParkID, &r.VisitorTicketID)
		checkErr(err)

		k := string(r.ParkID + r.VisitorTicketID)
		Miniworld.GoesToCollection[k] = r
	}

	//InteractsWith
	rows, err = db.Query("select Visitor_Ticket_ID, Attraction_Name from Interacts_With")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r InteractsWith

		err = rows.Scan(&r.VisitorTicketID, &r.AttractionName)
		checkErr(err)

		k := string(r.VisitorTicketID) + r.AttractionName
		Miniworld.InteractsWithCollection[k] = r
	}

	//View: Park Name & Address
	rows, err = db.Query("select Name, Address from Parks_And_Locations")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r ViewParkNameAndAddress

		err = rows.Scan(&r.ParkName, &r.Address)
		checkErr(err)

		Miniworld.ParkNameAndAddressView[r.ParkName] = r
	}

	//View: Visitor Name & Park Name
	rows, err = db.Query("select Visitor_Name, Park_Name from Visitors_At_Parks")

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var r ViewVisitorNameAndParkName

		err = rows.Scan(&r.VisitorName, &r.ParkName)
		checkErr(err)

		k := r.VisitorName + r.ParkName
		Miniworld.VisitorNameAndParkNameView[k] = r
	}
}

//runs the insert sql command
func insertIntoDB(tableName string, rowData string) {
	_, err := db.Exec("INSERT INTO " + tableName + " VALUES (" + rowData + ");")
	checkErr(err)

	//repopulate structs to reflect new data
	populateMaps(db)
}

//executes queries
func executeChooseSelectQuery(sc string, fc string, wc string, e string, hostVar string) *sql.Rows {
	queryStmt := "SELECT " + sc + " FROM " + fc + " WHERE " + wc + e + ";"
	// fmt.Println(queryStmt)
	var rows *sql.Rows
	var err error

	if hostVar == "" {
		rows, err = db.Query(queryStmt)
	} else {
		rows, err = db.Query(queryStmt, hostVar)
	}

	checkErr(err)
	return rows
}

//maps query results to struct
func chooseSelectQuery(input string, hostVar string) struct{ Rows map[string]QueryResultRow } {
	QueryResult := struct{ Rows map[string]QueryResultRow }{Rows: make(map[string]QueryResultRow)}

	switch input {
	case "VisitorsAtParks":
		//get the rows
		rows := executeChooseSelectQuery("v.Name, p.Name",
			"Amusement_Park p inner join Visitor v inner join Goes_To g",
			"g.Park_Id = p.Id and g.Visitor_Ticket_Id = v.Ticket_Id and p.Name = ?",
			"",
			hostVar)

		//save the rows to a struct
		for rows.Next() {
			var r QueryResultRow
			err := rows.Scan(&r.VisitorName, &r.ParkName)

			checkErr(err)
			QueryResult.Rows[r.VisitorName] = r
		}
	case "EmployeesWorkingAttractions":
		//get the rows
		rows := executeChooseSelectQuery("fname, lname, attraction_name, p.Name",
			"employee e inner join attraction a inner join amusement_park p",
			"a.Park_Id = p.Id and a.name = e.attraction_name and e.attraction_name = ?",
			"",
			hostVar)

		//save the rows to a struct
		for rows.Next() {
			var r QueryResultRow
			err := rows.Scan(&r.FirstName, &r.LastName, &r.AttractionName, &r.ParkName)

			checkErr(err)
			k := r.FirstName + r.LastName + r.AttractionName
			QueryResult.Rows[k] = r
		}
	}
	return QueryResult
}

/**********     Web functions     **********/
func mainHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")

	t.Execute(w, Miniworld)
	checkErr(err)

}

//handler for displaying queries
func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/query.html")
		t.Execute(w, nil)
	} else {
		//grab the query name
		m := validPath.FindStringSubmatch(r.URL.Path)
		input := m[2]

		//grab the query results
		QueryResult := chooseSelectQuery(input, r.FormValue("hostVariable"))

		t, err := template.ParseFiles("templates/" + input + ".html")
		checkErr(err)

		t.Execute(w, QueryResult)
		http.Redirect(w, r, "/query/"+m[2], http.StatusFound)
	}
}

//handler for displaying tables
func displayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/display.html")
		t.Execute(w, nil)
	} else {

		input := r.FormValue("tableName")

		t, err := template.ParseFiles("templates/" + input + ".html")
		checkErr(err)

		t.Execute(w, Miniworld)
	}
}

//handler for inserting values
func insertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/insert.html")
		t.Execute(w, nil)
	} else {

		//struct to send to template
		insertStruct := struct {
			TableName string
			RowValues string
		}{
			TableName: r.FormValue("dbName"),
			RowValues: r.FormValue("insertValue"),
		}

		insertIntoDB(insertStruct.TableName, insertStruct.RowValues)

		t, err := template.ParseFiles("templates/insert.html")
		checkErr(err)

		t.Execute(w, insertStruct)
	}
}

//function for clean error checking
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/**********     main function    **********/

func main() {
	//remove the db if it exists, user preference
	// os.Remove("parks_miniworld.db")

	/**********     set up the database    **********/
	fmt.Printf("\nStarting database...\n")

	defer db.Close()

	//populate the structs with info from the db
	populateMaps(db)

	/**********     set up the web page     **********/

	fmt.Printf("\nStarting web server...\n")

	//make the /resources directory visible
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	//set the handlers
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/display/", displayHandler)
	http.HandleFunc("/query/", queryHandler)
	http.HandleFunc("/insert/", insertHandler)

	//start the webpage
	err := http.ListenAndServe(":8080", nil) // set listen port

	checkErr(err)
}
