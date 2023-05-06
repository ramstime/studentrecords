package models

type Student struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
	//Subjects string `json:"subjects"`
	Address string `json:"address"`
}
