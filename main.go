package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID      int     `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Postion string  `json:"position,omitempty"`
	Salary  float64 `json:"salary,omitempty"`
}

var employees []Employee

var idLock sync.Mutex
var id int

func GenerateID() int {
	idLock.Lock()
	defer idLock.Unlock()
	id++
	return id
}

// Handlers
// List Employee based on page number and per page records
func ListEmployee(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Param("page"))

	if err != nil {
		ctx.JSON(http.StatusForbidden, "error while parsing page number: "+err.Error())
		return
	}

	perpage, err := strconv.Atoi(ctx.Param("perpage"))

	if err != nil {
		ctx.JSON(http.StatusForbidden, "error while parsing per page records: "+err.Error())
		return
	}

	start := perpage * (page - 1)

	// For first page allowing any numbers of employee per page
	if page == 1 {
		start = 0
	}

	if len(employees)-1 <= start {
		ctx.JSON(http.StatusForbidden, "requested record does not found. start>=records")
		return
	}

	end := start + perpage

	var response []Employee
	if end > len(employees)-1 {
		response = employees[start:]
	} else {
		response = employees[start:end]
	}

	ctx.JSON(http.StatusOK, response)
	log.Println("Page number: ", page, "per page Records: ", perpage)
}

// Get EmployeeByID

func GetEmployeeById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Error: Unable to parse ID from request")
		return
	}

	for _, employee := range employees {
		if employee.ID == id {
			ctx.JSON(http.StatusOK, employee)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, fmt.Sprintf("Error: we don't have any employee associated with given Id: %d", id))

}

// CreateEmployee

func CreateEmployee(ctx *gin.Context) {
	var employee Employee
	if err := ctx.Bind(&employee); err != nil {
		ctx.JSON(http.StatusForbidden, "Error: unable to parse request body from request")
		return
	}

	employee.ID = GenerateID()

	employees = append(employees, employee)

	ctx.JSON(http.StatusOK, employee)
}

// Update Employee

func UpdateEmployee(ctx *gin.Context) {

	var employee Employee
	if err := ctx.Bind(&employee); err != nil {
		ctx.JSON(http.StatusForbidden, "Error: unable to parse request body from request")
		return
	}

	for i, emp := range employees {
		if emp.ID == employee.ID {
			// Updated employee
			employees[i] = employee
			ctx.JSON(http.StatusOK, employee)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, fmt.Sprintf("Error: we don't have any employee associated with given Id: %d to update", employee.ID))
}

// DeleteEmployee

func DeleteEmployee(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Error: Unable to parse ID from request for deletion")
		return
	}

	for i, employee := range employees {
		if employee.ID == id {
			employees = append(employees[0:i], employees[i+1:]...)
			ctx.JSON(http.StatusOK, employee)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, fmt.Sprintf("Error: we don't have any employee associated with given Id: %d to delete", id))
}

func StartApp(server *http.Server) error {
	app := gin.Default()
	app.GET("/employee/get/:id", GetEmployeeById)
	app.GET("employee/list/:page/:perpage", ListEmployee)
	app.POST("/employee/create", CreateEmployee)
	app.PUT("/employee/update", UpdateEmployee)
	app.DELETE("/employee/delete/:id", DeleteEmployee)
	server.Handler = app
	server.Addr = ":8080"

	return server.ListenAndServe()
}

func main() {
	server := &http.Server{}
	fmt.Println(StartApp(server))
}
