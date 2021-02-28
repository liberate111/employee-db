package main

import (
	"emp/employees"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/emps", employees.Index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/emps/show", employees.Show)
	http.HandleFunc("/emps/create", employees.Create)
	http.HandleFunc("/emps/create/process", employees.CreateProcess)
	http.HandleFunc("/emps/update", employees.Update)
	http.HandleFunc("/emps/update/process", employees.UpdateProcess)
	http.HandleFunc("/emps/delete/process", employees.DeleteProcess)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/emps", http.StatusSeeOther)
}
