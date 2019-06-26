package controller

import (
	"awesomeProject/model"
	"html/template"
	"log"
	"net/http"
)

type Index struct{}

func (p Index) GetIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		person := model.IndexModel{ Version: "0.1", Name: "Book list web api service"}
		parsedTemplate, _ := template.ParseFiles("templates/index.html")
		err := parsedTemplate.Execute(w, person)
		if err != nil   {
			log.Printf("Error occurred while executing the template    or writing its output : ", err)
			return
		}
	}
}
