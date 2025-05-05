package model

import "time"

type DBSHeader struct {
	MsgID     string    `json:"msgId" bson:"msgId"`
	OrgID     string    `json:"orgId" bson:"orgId"`
	TimeStamp time.Time `json:"timeStamp" bson:"timeStamp"`
	Country   string    `json:"ctry" bson:"ctry"`
}

type DBSParty struct {
	Name         string `json:"name" bson:"name"`
	AccountNo    string `json:"accountNo,omitempty" bson:"accountNo,omitempty"`
	SenderBankID string `json:"senderBankId,omitempty" bson:"senderBankId,omitempty"`
}

type DBSAmountDetails struct {
	TxnCurrency string  `json:"txnCcy" bson:"txnCcy"`
	TxnAmount   float64 `json:"txnAmt" bson:"txnAmt"`
}

type DBSRmtInfo struct {
	PaymentDetails string `json:"paymentDetails,omitempty" bson:"paymentDetails,omitempty"`
	PurposeCode    string `json:"purposeCode,omitempty" bson:"purposeCode,omitempty"`
}

type DBSTxnInfo struct {
	TxnType           string           `json:"txnType" bson:"txnType"`
	CustomerReference string           `json:"customerReference" bson:"customerReference"`
	TxnRefID          string           `json:"txnRefId" bson:"txnRefId"`
	TxnDate           string           `json:"txnDate" bson:"txnDate"`
	ValueDate         string           `json:"valueDt" bson:"valueDt"`
	ReceivingParty    DBSParty         `json:"receivingParty" bson:"receivingParty"`
	AmountDetails     DBSAmountDetails `json:"amtDtls" bson:"amtDtls"`
	SenderParty       DBSParty         `json:"senderParty" bson:"senderParty"`
	RmtInfo           *DBSRmtInfo      `json:"rmtInf,omitempty" bson:"rmtInf,omitempty"`
}

type DBSTransaction struct {
	Header  DBSHeader  `json:"header" bson:"header"`
	TxnInfo DBSTxnInfo `json:"txnInfo" bson:"txnInfo"`
}
