package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Address1 string `bson:"address1,omitempty" json:"address1,omitempty"`
	Address2 string `bson:"address2,omitempty" json:"address2,omitempty"`
	Address3 string `bson:"address3,omitempty" json:"address3,omitempty"`
	Pincode  string `bson:"pinCode,omitempty" json:"pinCode,omitempty"`
	City     string `bson:"city,omitempty" json:"city,omitempty"`
	State    string `bson:"state,omitempty" json:"state,omitempty"`
	Country  string `bson:"country,omitempty" json:"country,omitempty"`
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email          string             `bson:"email,omitempty" json:"emailId,omitempty"`
	FirstName      string             `bson:"firstName,omitempty" json:"firstName,omitempty"`
	LastName       string             `bson:"lastName,omitempty" json:"lastName,omitempty"`
	Password       string             `bson:"password,omitempty" json:"password,omitempty"`
	MobileNumber   string             `bson:"mobileNumber,omitempty" json:"mobileNumber,omitempty"`
	BillingAddress Address            `bson:"billingAddress,omitempty" json:"billingAddress,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
