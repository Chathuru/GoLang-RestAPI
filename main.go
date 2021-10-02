package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type Employee struct {
	Emp_no     int
	First_name string
	Last_name  string
	Gender     string
	Birth_date string
	Hire_date  string
}

var message = map[string]string{"message": "hi..."}

func hello(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, message)
}

func api() {
	r := gin.Default()
	r.GET("/", hello)
	r.Run()
}

func CreateTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS employees(
		Emp_no int NOT NULL PRIMARY KEY,
		First_name TEXT,
		Last_name TEXT,
		Gender TEXT,
		Birth_date TEXT,
		Hire_date TEXT
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("DB Table creating fail")
	}
}

func ListAllEmployes(db *sql.DB) []Employee {
	var employees = []Employee{}
	result, _ := db.Query("SELECT * FROM employees")

	for result.Next() {
		var emp_no int
		var first_name, last_name, gender, birth_date, hire_date string
		result.Scan(&emp_no, &first_name, &last_name, &gender, &birth_date, &hire_date)
		employees = append(employees, Employee{emp_no, first_name, last_name, gender, birth_date, hire_date})
	}

	return employees
}

func SelectById(db *sql.DB, empNo int) (Employee, error) {
	var emploee Employee
	query_result := db.QueryRow("SELECT * FROM employees WHERE emp_no = ?", empNo)
	err := query_result.Scan(&emploee.Emp_no, &emploee.First_name, &emploee.Last_name, &emploee.Gender, &emploee.Birth_date, &emploee.Hire_date)

	if err != nil {
		return emploee, err
	}
	return emploee, nil
}

func AddEmployee(db *sql.DB, employee Employee) {
	query, _ := db.Prepare("INSERT INTO employees(Emp_no, First_name, Last_name, Gender, Birth_date, Hire_date) VALUES (?, ?, ?, ?, ?, ?)")
	_, err := query.Exec(employee.Emp_no, employee.First_name, employee.Last_name, employee.Gender, employee.Birth_date, employee.Hire_date)
	if err != nil {
		log.Fatal("Insert fail")
	}
}

func DeleteById(db *sql.DB, empNo int) {
	query, _ := db.Prepare("DELETE FROM employees WHERE emp_no = ?")
	query.Exec(empNo)
}

func main() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal("Error, DB not created")
	}
	CreateTable(db)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", hello)
	r.GET("/ListAllEmployes", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, ListAllEmployes(db))
	})

	r.GET("/GetById/:empNo", func(context *gin.Context) {
		empNo := context.Params.ByName("empNo")
		id, _ := strconv.Atoi(empNo)
		output, err := SelectById(db, id)

		if err != nil {
			context.JSON(http.StatusOK, gin.H{"message": err.Error()})
		} else {
			context.IndentedJSON(http.StatusOK, output)
		}
	})

	r.POST("/AddEmployee", func(context *gin.Context) {
		var employee Employee
		err := context.ShouldBind(&employee)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			AddEmployee(db, employee)
			context.String(http.StatusOK, "Success")
		}
	})

	r.GET("/DeleteById/:empNo", func(context *gin.Context) {
		empNo := context.Params.ByName("empNo")
		id, _ := strconv.Atoi(empNo)
		DeleteById(db, id)
		context.JSON(http.StatusOK, gin.H{"message": "Employee " + empNo + " deleted successfully"})
	})

	r.GET("/bad", func(context *gin.Context) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Bad"})
	})

	r.Run("0.0.0.0:8080")
}
