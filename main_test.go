package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server = new(http.Server)

// Create server

func setup() {
	go StartApp(server)
}

func shutdown() {
	if err := server.Shutdown(context.Background()); err != nil {
		log.Println(err.Error())
	}
}

func reset() {
	employees = []Employee{}
}

//Create employee Helper function

func createEmp(emp Employee) (Employee, error) {
	payload, _ := json.Marshal(&emp)
	resp, err := http.Post("http://localhost:8080/employee/create", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return Employee{}, err
	}
	defer resp.Body.Close()
	var respEmp Employee
	err = json.NewDecoder(resp.Body).Decode(&respEmp)
	if err != nil {
		return Employee{}, nil
	}
	return respEmp, nil
}

// GetEmployeeByID Helper function
func getEmployeebyId(ID int) (Employee, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/employee/get/%d", ID))
	if err != nil {
		return Employee{}, nil
	}
	defer resp.Body.Close()

	var respEmp Employee
	err = json.NewDecoder(resp.Body).Decode(&respEmp)
	if err != nil {
		return Employee{}, nil
	}
	return respEmp, nil
}

// UpdateEmployee Helper function
func updateEmplyee(emp Employee) (Employee, error) {
	payload, err := json.Marshal(&emp)
	if err != nil {
		return Employee{}, nil
	}
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/employee/update", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return Employee{}, err
	}
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return Employee{}, err
	}
	defer resp.Body.Close()

	var respEmp Employee
	err = json.NewDecoder(resp.Body).Decode(&respEmp)
	if err != nil {
		return Employee{}, err
	}
	return respEmp, nil
}

// DeleteEmployee helper function
func deleteEmployee(ID int) (Employee, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/employee/delete/%d", ID), nil)

	if err != nil {
		return Employee{}, err
	}
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return Employee{}, err
	}
	defer resp.Body.Close()

	var respEmp Employee
	err = json.NewDecoder(resp.Body).Decode(&respEmp)
	if err != nil {
		return Employee{}, err
	}
	return respEmp, nil

}

func listEmployee(page, perpage int) ([]Employee, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/employee/list/%d/%d", page, perpage), nil)

	if err != nil {
		return nil, err
	}
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respEmp []Employee
	if resp.StatusCode != http.StatusOK {
		msg := ""
		json.NewDecoder(resp.Body).Decode(&msg)
		return nil, errors.New(msg)
	}
	err = json.NewDecoder(resp.Body).Decode(&respEmp)
	if err != nil {
		return nil, err
	}
	return respEmp, nil
}

func TestCreateEmployee(t *testing.T) {

	emp := Employee{
		Name:    "Person1 Surname1",
		Postion: "SWE",
		Salary:  80.4,
	}

	respEmp, err := createEmp(emp)

	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, emp.Name, respEmp.Name)
	assert.Equal(t, emp.Postion, respEmp.Postion)
	assert.Equal(t, emp.Salary, respEmp.Salary)
	assert.Equal(t, len(employees) > 0, true)
}

func TestGetEmployeeById(t *testing.T) {
	emp := Employee{
		Name:    "Person2 Surname2",
		Postion: "SWE",
		Salary:  80.4,
	}
	createdEmp, err := createEmp(emp)
	if err != nil {
		t.Error(err)
	}
	emp2, err := getEmployeebyId(createdEmp.ID)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, emp2.ID, createdEmp.ID)
}

func TestUpdateEmployee(t *testing.T) {
	emp := Employee{
		Name:    "Person3 Surname3",
		Postion: "SWE",
		Salary:  80.4,
	}
	updateEmp, err := createEmp(emp)
	if err != nil {
		t.Error(err)
	}

	updateEmp.Name = "UPdateName3"
	updateEmp.Salary = 100.5

	updated, err := updateEmplyee(updateEmp)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, updated.ID, updateEmp.ID)
	assert.Equal(t, updated.Name, updateEmp.Name)
	assert.Equal(t, updated.Postion, updateEmp.Postion)
	assert.Equal(t, updated.Salary, updateEmp.Salary)
}

func TestDeleteEmployee(t *testing.T) {
	emp := Employee{
		Name:    "Person3 Surname3",
		Postion: "SWE",
		Salary:  80.4,
	}
	emp2, err := createEmp(emp)
	if err != nil {
		t.Error(err)
	}

	// Length of employees slice before delete
	lenbeforeDelete := len(employees)
	deletedEmp, err := deleteEmployee(emp2.ID)
	if err != nil {
		t.Error(err)
	}

	lenAfterDelete := len(employees)
	assert.Equal(t, emp2.ID, deletedEmp.ID)
	assert.Equal(t, lenAfterDelete < lenbeforeDelete, true)
}

func TestListEmployeePagination(t *testing.T) {
	// Reseting the employees in memory Data
	reset()
	// Create more than 100 Employee
	for i := 0; i < 106; i++ {
		emp := Employee{
			Name:    fmt.Sprintf("Name%d Surname%d", i, i),
			Postion: fmt.Sprintf("SWE-%d", i),
			Salary:  float64(10 * i),
		}
		_, err := createEmp(emp)
		if err != nil {
			t.Error(err)
		}
	}
	page_number := 5
	per_page := 10
	data, err := listEmployee(page_number, per_page)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, len(data), per_page)

	page_number = 11

	// 11th page only contains 6 employee
	data, err = listEmployee(page_number, per_page)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, len(data), 6)

	// Request the page which have no data

	page_number = 12
	_, err = listEmployee(page_number, per_page)

	assert.Errorf(t, err, err.Error(), "requested record does not found. start>=records")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)

}
