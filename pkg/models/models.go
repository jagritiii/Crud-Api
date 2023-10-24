package models

type User struct {
	Id         int    `json:"Id" bson:"Id"`
	Company    string `json:"Company" bson:"Company" validate:"required,min=1,max=20"`
	Profile    string `json:"Profile" bson:"Profile" validate:"required"`
	Age        int    `json:"Age" bson:"Age" validate:"required"`
	Experience int    `json:"Exp" bson:"Exp" validate:"required"`
}
