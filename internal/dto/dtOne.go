package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	AvailabilityZones                   []string           `json:"availability_zones"`
	Benefits                            []Benefit          `json:"benefits"`
	Description                         string             `json:"description"`
	Destination                         Amount             `json:"destination"`
	ID                                  primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ProductId                           int                `bson:"product_id,omitempty" json:"id"`
	Name                                string             `json:"name"`
	Operator                            Operator           `json:"operator"`
	Prices                              Prices             `json:"prices"`
	Promotions                          interface{}        `json:"promotions"` // or a specific type if known
	Rates                               Rates              `json:"rates"`
	Regions                             interface{}        `json:"regions"` // can be []Region if needed
	RequiredAdditionalIdentifierFields  interface{}        `json:"required_additional_identifier_fields"`
	RequiredBeneficiaryFields           interface{}        `json:"required_beneficiary_fields"`
	RequiredCreditPartyIdentifierFields [][]string         `json:"required_credit_party_identifier_fields"`
	RequiredDebitPartyIdentifierFields  interface{}        `json:"required_debit_party_identifier_fields"`
	RequiredSenderFields                interface{}        `json:"required_sender_fields"`
	RequiredStatementIdentifierFields   interface{}        `json:"required_statement_identifier_fields"`
	Service                             Service            `json:"service"`
	Source                              Amount             `json:"source"`
	Tags                                interface{}        `json:"tags"` // or []string if tags are string-based
	Type                                string             `json:"type"`
	Validity                            Validity           `json:"validity"`
}

type Benefit struct {
	AdditionalInformation interface{}     `json:"additional_information"`
	Amount                AmountBreakdown `json:"amount"`
	Type                  string          `json:"type"`
	Unit                  string          `json:"unit"`
	UnitType              string          `json:"unit_type"`
}

type AmountBreakdown struct {
	Base              interface{} `json:"base"`
	PromotionBonus    interface{} `json:"promotion_bonus"`
	TotalExcludingTax interface{} `json:"total_excluding_tax"`
	TotalIncludingTax interface{} `json:"total_including_tax"`
}

type AmountWithFee struct {
	Amount   interface{} `json:"amount"`
	Fee      interface{} `json:"fee"`
	Unit     string      `json:"unit"`
	UnitType string      `json:"unit_type"`
}

type Amount struct {
	Amount   interface{} `json:"amount"`
	Unit     string      `json:"unit"`
	UnitType string      `json:"unit_type"`
}

type Operator struct {
	Country Country     `json:"country"`
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Regions interface{} `json:"regions"` // or []Region if regions are consistent
}

type Country struct {
	ISOCode string   `json:"iso_code"`
	Name    string   `json:"name"`
	Regions []Region `json:"regions"`
}

type Region struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Prices struct {
	Retail    interface{} `json:"retail"` // nullable or specific price type
	Wholesale Amount      `json:"wholesale"`
}

type Rates struct {
	Base      float64     `json:"base"`
	Retail    interface{} `json:"retail"`
	Wholesale float64     `json:"wholesale"`
}

type Service struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Subservice Subservice `json:"subservice"`
}

type Subservice struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Validity struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}
