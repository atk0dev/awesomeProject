package model

type Book struct {
	Id     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   string `json:year`
}
