package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Struct

type Knowledge struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Detail    string             `json:"detail,omitempty" bson:"detail,omitempty"`
	Viewcount int                `json:"viewcount,omitempty" bson:"viewcount,omitempty"`
}
