package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	Address3 string `json:"address3,omitempty"`
	Pincode  string `json:"pincode,omitempty"`
	City     string `json:"city,omitempty"`
	State    string `json:"state,omitempty"`
	Country  string `json:"country,omitempty"`
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmailID        string             `json:"email_id,omitempty"`
	FirstName      string             `json:"first_name,omitempty"`
	LastName       string             `json:"last_name,omitempty"`
	Password       string             `json:"password,omitempty"`
	MobileNumber   string             `json:"mobile_number,omitempty"`
	BillingAddress Address            `json:"billing_address,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
