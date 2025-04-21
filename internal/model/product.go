package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Amount struct {
	Base              interface{} `bson:"base" json:"base"`
	PromotionBonus    interface{} `bson:"promotion_bonus" json:"promotion_bonus"`
	TotalExcludingTax interface{} `bson:"total_excluding_tax" json:"total_excluding_tax"`
	TotalIncludingTax interface{} `bson:"total_including_tax" json:"total_including_tax"`
}

type Prices struct {
	Retail    interface{} `bson:"retail" json:"retail"`
	Wholesale Wholesale   `bson:"wholesale" json:"wholesale"`
}

type Wholesale struct {
	Amount   interface{} `bson:"amount" json:"amount"`
	Fee      interface{} `bson:"fee" json:"fee"`
	Unit     string      `bson:"unit" json:"unit"`
	UnitType string      `bson:"unit_type" json:"unit_type"`
}

type Rates struct {
	Base      interface{} `bson:"base" json:"base"`
	Retail    interface{} `bson:"retail" json:"retail"`
	Wholesale interface{} `bson:"wholesale" json:"wholesale"`
}

type Operator struct {
	ID      int         `bson:"id" json:"id"`
	Name    string      `bson:"name" json:"name"`
	Country Country     `bson:"country" json:"country"`
	Regions interface{} `bson:"regions" json:"regions"`
}

type Country struct {
	ISOCode string   `bson:"iso_code" json:"iso_code"`
	Name    string   `bson:"name" json:"name"`
	Regions []Region `bson:"regions" json:"regions"`
}

type Region struct {
	Code string `bson:"code" json:"code"`
	Name string `bson:"name" json:"name"`
}

type Service struct {
	ID         int        `bson:"id" json:"id"`
	Name       string     `bson:"name" json:"name"`
	SubService SubService `bson:"subservice" json:"subservice"`
}

type SubService struct {
	ID   int    `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Status struct {
	Class   StatusClass `bson:"class" json:"class"`
	ID      int         `bson:"id" json:"id"`
	Message string      `bson:"message" json:"message"`
}

type StatusClass struct {
	ID      int    `bson:"id" json:"id"`
	Message string `bson:"message" json:"message"`
}

// Product Model
type Product struct {
	AvailabilityZones                   []string           `bson:"availability_zones" json:"availability_zones"`
	Benefits                            []Benefit          `bson:"benefits" json:"benefits"`
	Description                         string             `bson:"description" json:"description"`
	Destination                         Amount             `bson:"destination" json:"destination"`
	ID                                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UniqueId                            int                `bson:"unique_id,omitempty" json:"id"`
	Name                                string             `bson:"name" json:"name"`
	Operator                            Operator           `bson:"operator" json:"operator"`
	Prices                              Prices             `bson:"prices" json:"prices"`
	Promotions                          interface{}        `bson:"promotions" json:"promotions"`
	Rates                               Rates              `bson:"rates" json:"rates"`
	Regions                             interface{}        `bson:"regions" json:"regions"`
	RequiredAdditionalIdentifierFields  interface{}        `bson:"required_additional_identifier_fields" json:"required_additional_identifier_fields"`
	RequiredBeneficiaryFields           interface{}        `bson:"required_beneficiary_fields" json:"required_beneficiary_fields"`
	RequiredCreditPartyIdentifierFields [][]string         `bson:"required_credit_party_identifier_fields" json:"required_credit_party_identifier_fields"`
	RequiredDebitPartyIdentifierFields  interface{}        `bson:"required_debit_party_identifier_fields" json:"required_debit_party_identifier_fields"`
	RequiredSenderFields                interface{}        `bson:"required_sender_fields" json:"required_sender_fields"`
	RequiredStatementIdentifierFields   interface{}        `bson:"required_statement_identifier_fields" json:"required_statement_identifier_fields"`
	Service                             Service            `bson:"service" json:"service"`
	Source                              Amount             `bson:"source" json:"source"`
	Tags                                interface{}        `bson:"tags" json:"tags"`
	Type                                string             `bson:"type" json:"type"`
	Validity                            Validity           `bson:"validity" json:"validity"`
}

// Benefit Model
type Benefit struct {
	AdditionalInformation interface{} `bson:"additional_information" json:"additional_information"`
	Amount                Amount      `bson:"amount" json:"amount"`
	Type                  string      `bson:"type" json:"type"`
	Unit                  string      `bson:"unit" json:"unit"`
	UnitType              string      `bson:"unit_type" json:"unit_type"`
}

// ProductTransaction Model

type ProductTransaction struct {
	Benefits                   []Benefit             `bson:"benefits" json:"benefits"`
	ConfirmationDate           time.Time             `bson:"confirmation_date" json:"confirmation_date"`
	ConfirmationExpirationDate time.Time             `bson:"confirmation_expiration_date" json:"confirmation_expiration_date"`
	CreationDate               time.Time             `bson:"creation_date" json:"creation_date"`
	CreditPartyIdentifier      CreditPartyIdentifier `bson:"credit_party_identifier" json:"credit_party_identifier"`
	ExternalID                 string                `bson:"external_id" json:"external_id"`
	ID                         int64                 `bson:"id" json:"id"`
	OperatorReference          string                `bson:"operator_reference" json:"operator_reference"`
	Pin                        Pin                   `bson:"pin" json:"pin"`
	Prices                     Prices                `bson:"prices" json:"prices"`
	Product                    Product               `bson:"product" json:"product"`
	Promotions                 interface{}           `bson:"promotions" json:"promotions"`
	Rates                      Rates                 `bson:"rates" json:"rates"`
	Status                     Status                `bson:"status" json:"status"`
	CreatedAt                  time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt                  time.Time             `bson:"updated_at" json:"updated_at"`
	UpdationTime               string                `bson:"updation_time" json:"updation_time"`
	DeletedAt                  time.Time             `bson:"deleted_at" json:"deleted_at"`
}

type ProductPinItem struct {
	ExternalID string `bson:"external_id"`
	ProductID  int    `bson:"productId"`
	Pin        struct {
		Code   string `bson:"code"`
		Serial string `bson:"serial"`
	} `bson:"pin"`
}

type ProductPin struct {
	OrderID     string           `bson:"orderID"`
	ProductPins []ProductPinItem `bson:"productPins"`
	CreatedAt   time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time        `bson:"updated_at" json:"updated_at"`
	DeletedAt   time.Time        `bson:"deleted_at" json:"deleted_at"`
}

// CreditPartyIdentifier Model
type CreditPartyIdentifier struct {
	MobileNumber string `bson:"mobile_number" json:"mobile_number"`
}

// Validity Model
type Validity struct {
	Quantity int    `bson:"quantity" json:"quantity"`
	Unit     string `bson:"unit" json:"unit"`
}

type Pin struct {
	Code   string `bson:"code" json:"code"`
	Serial string `bson:"serial" json:"serial"`
}
