package employees

import (
	"context"
	"emp/config"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee struct {
	// ID     bson.ObjectId // `bson:"_id"`
	EmpID     string  `bson:"empId"`
	Firstname string  `bson:"firstname"`
	Lastname  string  `bson:"lastname"`
	Position  string  `bson:"position"`
	Salary    float32 `bson:"salary"`
	Bonus     float32 `bson:"bonus,omitempty"`
	Total     float32 `bson:"total,omitempty"`
}

func AllEmps() ([]Employee, error) {
	emps := []Employee{}
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"empId", 1}})
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
func OneEmp(r *http.Request) (Employee, error) {
	emp := Employee{}
	empId := r.FormValue("empid")
	if empId == "" {
		return emp, errors.New("400. Bad Request.")
	}
	// Find one
	filter := bson.D{{"empId", empId}}
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
		fmt.Println("No Bonus & total")
	}

	return emp, nil
}

func PutEmp(r *http.Request) (Employee, error) {
	// get form values
	emp := Employee{}
	emp.EmpID = r.FormValue("empid")
	emp.Firstname = r.FormValue("firstname")
	emp.Lastname = r.FormValue("lastname")
	emp.Position = r.FormValue("position")
	s := r.FormValue("salary")

	// validate form values
	if emp.EmpID == "" || emp.Firstname == "" || emp.Lastname == "" || emp.Position == "" || s == "" {
		return emp, errors.New("400. Bad request. All fields must be complete.")
	}

	// convert form values
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return emp, errors.New("406. Not Acceptable. Price must be a number.")
	}
	emp.Salary = float32(f64)

	// insert the document
	res, err := config.Coll.InsertOne(context.TODO(), emp)
	if err != nil {
		return emp, errors.New("500. Internal Server Error." + err.Error())
	}
	fmt.Printf("inserted document with ID %v\n", res.InsertedID)
	return emp, nil
}

func UpdateEmp(r *http.Request) (Employee, error) {
	// get form values
	emp := Employee{}
	emp.EmpID = r.FormValue("empid")
	emp.Firstname = r.FormValue("firstname")
	emp.Lastname = r.FormValue("lastname")
	emp.Position = r.FormValue("position")
	s := r.FormValue("salary")

	// validate form values
	if emp.EmpID == "" || emp.Firstname == "" || emp.Lastname == "" || emp.Position == "" || s == "" {
		return emp, errors.New("400. Bad request. All fields must be complete.")
	}

	// convert form values
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return emp, errors.New("406. Not Acceptable. Price must be a number.")
	}
	emp.Salary = float32(f64)

	// Update one
	filter := bson.D{{"empId", emp.EmpID}}
	update := bson.D{{"$set",
		bson.D{
			{"empId", emp.EmpID},
			{"firstname", emp.Firstname},
			{"lastname", emp.Lastname},
			{"position", emp.Position},
			{"salary", emp.Salary},
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

func DeleteEmp(r *http.Request) error {
	empId := r.FormValue("empid")
	if empId == "" {
		return errors.New("400. Bad Request.")
	}
	// Delete one
	res, err := config.Coll.DeleteOne(context.TODO(), bson.D{{"empId", empId}})
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	return nil
}
