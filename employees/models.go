package employees

import (
	"context"
	"emp/config"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee struct {
	// ID        bson.ObjectId `bson:"_id"`
	EmpID     string  `bson:"empId"`
	Firstname string  `bson:"firstname"`
	Lastname  string  `bson:"lastname"`
	Position  string  `bson:"position"`
	Salary    float64 `bson:"salary"`
	Bonus     float64 `bson:"bonus,omitempty"`
	Total     float64 `bson:"total,omitempty"`
	// SalaryStr string  `bson:"SalaryStr,omitempty"`
}

func AllEmps() ([]Employee, error) {
	emps := []Employee{}
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{Key: "empId", Value: 1}})
	cursor, err := config.Coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	// get a list of all returned documents
	if err = cursor.All(context.TODO(), &emps); err != nil {
		return nil, err
	}
	return emps, nil
}
func OneEmp(c *fiber.Ctx) (Employee, error) {
	emp := Employee{}
	empId := c.FormValue("empid")
	if empId == "" {
		return emp, errors.New("400. Bad Request")
	}
	// Find one
	filter := bson.D{{Key: "empId", Value: empId}}
	err := config.Coll.FindOne(context.TODO(), filter).Decode(&emp)
	if err != nil {
		return emp, err
	}
	fmt.Printf("found document %v\n", emp)

	switch {
	case emp.Salary > 0 && emp.Salary <= 10000:
		emp.Bonus = 0.7 * emp.Salary
		emp.Total = emp.Salary + emp.Bonus
	case emp.Salary > 10000 && emp.Salary <= 20000:
		emp.Bonus = 0.6 * emp.Salary
		emp.Total = emp.Salary + emp.Bonus
	case emp.Salary > 20000:
		emp.Bonus = 0.5 * emp.Salary
		emp.Total = emp.Salary + emp.Bonus
	default:
		fmt.Println("default case")
	}

	return emp, nil
}

func PutEmp(c *fiber.Ctx) (Employee, error) {
	// get form values
	emp := Employee{}
	emp.EmpID = c.FormValue("empid")
	emp.Firstname = c.FormValue("firstname")
	emp.Lastname = c.FormValue("lastname")
	emp.Position = c.FormValue("position")
	s := c.FormValue("salary")

	// validate form values
	if emp.EmpID == "" || emp.Firstname == "" || emp.Lastname == "" || emp.Position == "" || s == "" {
		return emp, errors.New("400. Bad request. All fields must be complete")
	}

	// convert form values
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return emp, errors.New("406. Not Acceptable. Salary must be a number")
	}
	emp.Salary = f64

	// insert the document
	res, err := config.Coll.InsertOne(context.TODO(), emp)
	if err != nil {
		return emp, errors.New("500. Internal Server Error." + err.Error())
	}
	fmt.Printf("inserted document with ID %v\n", res.InsertedID)
	return emp, nil
}

func UpdateEmp(c *fiber.Ctx) (Employee, error) {
	// get form values
	emp := Employee{}
	emp.EmpID = c.FormValue("empid")
	emp.Firstname = c.FormValue("firstname")
	emp.Lastname = c.FormValue("lastname")
	emp.Position = c.FormValue("position")
	s := c.FormValue("salary")

	// validate form values
	if emp.EmpID == "" || emp.Firstname == "" || emp.Lastname == "" || emp.Position == "" || s == "" {
		return emp, errors.New("400. Bad request. All fields must be complete")
	}

	// convert form values
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return emp, errors.New("406. Not Acceptable. Salary must be a number")
	}
	emp.Salary = f64

	// Update one
	filter := bson.D{{Key: "empId", Value: emp.EmpID}}
	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "empId", Value: emp.EmpID},
			{Key: "firstname", Value: emp.Firstname},
			{Key: "lastname", Value: emp.Lastname},
			{Key: "position", Value: emp.Position},
			{Key: "salary", Value: emp.Salary},
		}}}
	res, err := config.Coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return emp, err
	}
	if res.MatchedCount != 0 {
		fmt.Println("matched an existing document")
	}
	return emp, nil
}

func DeleteEmp(c *fiber.Ctx) error {
	empId := c.FormValue("empid")
	if empId == "" {
		return errors.New("400. Bad Request")
	}
	// Delete one
	res, err := config.Coll.DeleteOne(context.TODO(), bson.D{{Key: "empId", Value: empId}})
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	// the employee might not exist
	if res.DeletedCount < 1 {
		return c.SendStatus(404)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	return nil
}
