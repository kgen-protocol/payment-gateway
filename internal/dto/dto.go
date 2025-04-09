package dto

import (
	"github.com/go-playground/validator/v10"
)

type UserRequest struct {
	FirstName      string `bson:"name,omitempty" validate:"required,min=2,max=100"`
	LastName       string `bson:"name,omitempty" validate:"required,min=2,max=100"`
	EmailID        string `bson:"email,omitempty" validate:"required,email"`
	Password       string `bson:"password,omitempty" json:"password" validate:"required,min=6"`
	MobileNumber   string `json:"mobile_number,omitempty"`
	BillingAddress string `json:"billing_address,omitempty"`
}

func (u *UserRequest) ValidateUser() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
