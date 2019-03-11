package models

import (
	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	Name       string `json : "name"`
	PhoneNumer string `json : "phone"`
	UserId     uint
}

// func validate(){

// }
// func getContact() {

// }
// func getContacts() {

// }
// func createContact() {

// }
