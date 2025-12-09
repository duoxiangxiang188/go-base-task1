package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int8
}
type Employee struct {
	EmployeeID int
	Person
}

func (e *Employee) PrintInfo() {
	fmt.Printf("EmployeeID = %d ，员工姓名 : %s , 年龄: %d \n", e.EmployeeID, e.Name, e.Age)
}

func main() {
	ems := []Employee{
		{
			EmployeeID: 1001,
			Person: Person{
				Name: "王",
				Age:  18,
			},
		},
		{
			EmployeeID: 1002,
			Person: Person{
				Name: "陈",
				Age:  19,
			},
		},
	}
	for _, e := range ems {
		e.PrintInfo()
	}
}
