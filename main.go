package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
        Cassandra "./cassandra" 
	"github.com/gorilla/mux"
)

type Emp struct {
	id        string `json:"empid"`
	firstName string `json:"first_name"`
	lastName  string `json:"last_name"`
	age       string `json:"age"`
}

type studentDetails struct {
	EmpId	  string
	FirstName string
}

type updEmp struct {
	EmpId	  string
	Age       string
}

type delEmpDetails struct {
    EmplId	string
}

type getEmp struct {
	Id        string `json:"empid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}
type AllEmpsResponse struct {
	Emps   []getEmp `json:"emps"`
	Emplen int
        Msg string
}

type DelEmpsResponse struct {
	Emps []getEmp
	ListLen		int
	SuccessMessage	string	
	FailedMessage	string
	InputMessage	string
}

type InputResponse struct {
	InputMessage		string
}

//Get All Employess Details
func getEmps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/dataDisplay.html"))
	var empList []getEmp
	m := map[string]interface{}{}
	fmt.Println("Getting all Employees")

	iter := Cassandra.Session.Query("SELECT empid,first_name,last_name,age FROM emps").Iter()
	for iter.MapScan(m) {
		empList = append(empList, getEmp{
			Id:        m["empid"].(string),
			FirstName: m["first_name"].(string),
			LastName:  m["last_name"].(string),
			Age:       m["age"].(int),
		})
		m = map[string]interface{}{}

	}
	tmpl.Execute(w, AllEmpsResponse{Emps: empList})
}


//Insert New employee
func insertData(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/forms.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := Emp{
		id:        r.FormValue("id"),
		firstName: r.FormValue("firstName"),
		lastName:  r.FormValue("lastName"),
		age:       r.FormValue("age"),
	}
	i, err := strconv.ParseInt(details.age, 10, 64)
	if err != nil {
		panic(err)
	}
	fmt.Println(i)
	// do something with details
	_ = details
	if err := Cassandra.Session.Query("INSERT INTO emps(empid, first_name, last_name, age) VALUES(?, ?, ?, ?)",
		details.id, details.firstName, details.lastName, details.age).Exec(); err != nil {
		fmt.Println("Error while inserting Emp")
		fmt.Println(err)
	}

	tmpl.Execute(w, struct{ Success bool }{true})
}


//Update Employee Information
func updateEmp(w http.ResponseWriter, r *http.Request) {
	var err error

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
        t := template.Must(template.ParseFiles("templates/update.html"))

	 if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}

	details := updEmp{
		EmpId: r.FormValue("emplid"),
		Age:   r.FormValue("age"),
	}
	     
	fmt.Println(" **** Updating Employee Age ****")
	if err := Cassandra.Session.Query("update emps set age = ? where empid = ?",
		details.Age, details.EmpId).Exec(); err != nil {
		fmt.Println("Error while inserting Emp")
		fmt.Println(err)
	}

t.Execute(w, struct{ Success bool }{true})
}

//Default Page
func goHome(w http.ResponseWriter, r *http.Request) {
	var err error

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    t := template.Must(template.ParseFiles("templates/home.html"))

	 if err != nil {
		fmt.Fprintf(w, "Unable to load template")
    }
t.Execute(w, struct{ Success bool }{true})
}


// Perform Employee Search with either ID or FirstName
func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/searchIntString.html"))
        var msg string
		
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}


	details := studentDetails{
		EmpId:     r.FormValue("empid"),
		FirstName: r.FormValue("firstname"), 		
	}
        
	var empList []getEmp
	m := map[string]interface{}{}

	if (details.EmpId != "") {
		iter := Cassandra.Session.Query("SELECT empid, first_name, last_name, age FROM emps WHERE empid = ? ALLOW FILTERING", details.EmpId).Iter()
		for iter.MapScan(m) {
			empList = append(empList, getEmp{
				Id:        m["empid"].(string),
				FirstName: m["first_name"].(string),
				LastName:  m["last_name"].(string),
				Age:       m["age"].(int),
			})
			m = map[string]interface{}{}
                        fmt.Println(empList)
		}
	} else {
		iter := Cassandra.Session.Query("SELECT empid, first_name, last_name, age FROM emps WHERE first_name = ? ALLOW FILTERING", details.FirstName).Iter()
		for iter.MapScan(m) {
			empList = append(empList, getEmp{
				Id:        m["empid"].(string),
				FirstName: m["first_name"].(string),
				LastName:  m["last_name"].(string),
				Age:       m["age"].(int),
			})
			m = map[string]interface{}{}
		}
	}
	

	Emplen := len(empList);
        if Emplen <= 0{
           msg = "No Results Found"
        }

	tmpl.Execute(w, AllEmpsResponse{Emps: empList, Emplen: Emplen, Msg:msg})	
}

//Delete Employee Record
func deleteData(w http.ResponseWriter, r *http.Request) {
	//dialog.MsgDlg("Press a button!");
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl := template.Must(template.ParseFiles("templates/delete.html"))
		
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}


	details := delEmpDetails{
		EmplId:   r.FormValue("delete"),
	}

	var empList []getEmp
	var availableEmpList []getEmp
	m := map[string]interface{}{}
	var successMessage string
	var failedMessage string
	var inputMessage string

	iter := Cassandra.Session.Query("SELECT * FROM emps WHERE empid = ?", details.EmplId).Iter() 
	for iter.MapScan(m) {
		empList = append(empList, getEmp{
			Id:        m["empid"].(string),
			FirstName: m["first_name"].(string),
			LastName:  m["last_name"].(string),
			Age:       m["age"].(int),
		})
		m = map[string]interface{}{}
	}

	listLen := len(empList);

	if(details.EmplId == "") {
		inputMessage = "Please enter the Employee Id to delete"
		tmpl.Execute(w, InputResponse{InputMessage: inputMessage})	
		return
	}

	if(listLen > 0) {
		if err := Cassandra.Session.Query("DELETE FROM emps WHERE empid = ?", details.EmplId).Exec(); 
		err != nil {
			fmt.Println("Error while deleting Emp")
			fmt.Println(err)
		} else {
			successMessage = "Employee Deleted Successfully"
			iter := Cassandra.Session.Query("SELECT * FROM emps").Iter() 
			for iter.MapScan(m) {
				availableEmpList = append(availableEmpList, getEmp{
					Id:        m["empid"].(string),
					FirstName: m["first_name"].(string),
					LastName:  m["last_name"].(string),
					Age:       m["age"].(int),
				})
				m = map[string]interface{}{}
			}
		}
	} else {
		failedMessage = "There is no Employee with that Employee Id"
	}
	tmpl.Execute(w, DelEmpsResponse{Emps: availableEmpList, ListLen: listLen, SuccessMessage: successMessage, FailedMessage: failedMessage})	
}

//Welome Page
func welcome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/welcomepage.html")
	if err != nil {
		fmt.Println(err)
	}
	items := struct {
		Homepage string
	}{
		Homepage: "Employee Home Page",
	}
	t.Execute(w, items)
}
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", welcome)
	router.HandleFunc("/insert", insertData)
        router.HandleFunc("/update", updateEmp)
	router.HandleFunc("/getEmps", getEmps)
	router.HandleFunc("/search", search)
        router.HandleFunc("/delete", deleteData)
        router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",http.FileServer(http.Dir("view/css/"))))
	http.ListenAndServe(":8080", router)
}
