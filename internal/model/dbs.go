package model

type Camt053Header struct {
	MsgId     string `json:"msgId"`
	OrgId     string `json:"orgId"`
	TimeStamp string `json:"timeStamp"`
	Country      string `json:"ctry"`
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
