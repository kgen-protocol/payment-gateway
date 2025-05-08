package model

type Camt053Header struct {
	MsgId     string `json:"msgId"`
	OrgId     string `json:"orgId"`
	TimeStamp string `json:"timeStamp"`
	Country   string `json:"ctry"`
}

type Camt053TxnInfo struct {
	AccountNo   string `json:"accountNo"`
	AccountCcy  string `json:"accountCcy"`
	BizDate     string `json:"bizDate"`
	MessageType string `json:"messageType"`
}

type Camt053Request struct {
	Header  Camt053Header  `json:"header"`
	TxnInfo Camt053TxnInfo `json:"txnInfo"`
}

type Camt053Response struct {
	Header         Camt053Header  `json:"header"`
	TxnEnqResponse TxnEnqResponse `json:"txnEnqResponse"`
}

type TxnEnqResponse struct {
	EnqStatus            string        `json:"enqStatus"`
	EnqRejectCode        string        `json:"enqRejectCode,omitempty"`
	EnqStatusDescription string        `json:"enqStatusDescription,omitempty"`
	AcctInfo             *AcctInfo     `json:"acctInfo,omitempty"`
	BizDate              string        `json:"bizDate,omitempty"`
	MessageType          string        `json:"messageType,omitempty"`
	Statement            []interface{} `json:"statement,omitempty"`
}

type AcctInfo struct {
	AccountNo  string `json:"accountNo"`
	AccountCcy string `json:"accountCcy"`
}

type NotificationHeader struct {
	MsgID     string `json:"msgId"`
	OrgID     string `json:"orgId"`
	TimeStamp string `json:"timeStamp"`
	Country   string `json:"ctry"`
}

type ReceivingParty struct {
	Name             string `json:"name"`
	AccountNo        string `json:"accountNo"`
	VirtualAccountNo string `json:"virtualAccountNo,omitempty"`
}

type AmountDetails struct {
	TxnCurrency string  `json:"txnCcy"`
	TxnAmount   float64 `json:"txnAmt"`
}

type SenderParty struct {
	Name         string `json:"name"`
	AccountNo    string `json:"accountNo"`
	SenderBankID string `json:"senderBankId"`
}

type TxnInfo struct {
	TxnType           string         `json:"txnType"`
	CustomerReference string         `json:"customerReference"`
	TxnRefID          string         `json:"txnRefId"`
	TxnDate           string         `json:"txnDate"`
	ValueDate         string         `json:"valueDt"`
	ReceivingParty    ReceivingParty `json:"receivingParty"`
	AmountDetails     AmountDetails  `json:"amtDtls"`
	SenderParty       SenderParty    `json:"senderParty"`
	PaymentDetails    string         `json:"paymentDetails"`
}

type IntradayNotificationPayload struct {
	Header  NotificationHeader `json:"header"`
	TxnInfo TxnInfo            `json:"txnInfo"`
}

type IncomingNotificationPayload struct {
	Header  NotificationHeader `json:"header"`
	TxnInfo TxnInfo            `json:"txnInfo"`
}
