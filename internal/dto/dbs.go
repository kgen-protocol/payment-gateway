package dto

type CAMT053Request struct {
	Header         Header         `json:"header"`
	TxnEnqResponse TxnEnqResponse `json:"txnEnqResponse"`
}

type Header struct {
	MsgID     string `json:"msgId"`
	OrgID     string `json:"orgId"`
	TimeStamp string `json:"timeStamp"`
	Country   string `json:"ctry"`
}

type TxnEnqResponse struct {
	EnqStatus   string             `json:"enqStatus"`
	AcctInfo    AcctInfo           `json:"acctInfo"`
	BizDate     string             `json:"bizDate"`
	MessageType string             `json:"messageType"`
	Statement   []StatementWrapper `json:"statement"`
}

type AcctInfo struct {
	AccountNo  string `json:"accountNo"`
	AccountCcy string `json:"accountCcy"`
}

type StatementWrapper struct {
	BkToCstmrStmt BankToCustomerStatement `json:"bkToCstmrStmt"`
}

type BankToCustomerStatement struct {
	GrpHdr GroupHeader `json:"grpHdr"`
	Stmt   []Statement `json:"stmt"`
}

type GroupHeader struct {
	MsgID   string `json:"msgId"`
	CreDtTm string `json:"creDtTm"`
}

type Statement struct {
	ID      string     `json:"id"`
	CreDtTm string     `json:"creDtTm"`
	Acct    Account    `json:"acct"`
	Bal     []Balance  `json:"bal"`
	TxsSumm TxnSummary `json:"txsSummry"`
	Ntry    []Entry    `json:"ntry"`
}

type Account struct {
	ID   AccountID   `json:"id"`
	Ccy  string      `json:"ccy"`
	Nm   string      `json:"nm"`
	Svcr AccountSvcr `json:"svcr"`
}

type AccountID struct {
	Othr IDValue `json:"othr"`
}

type IDValue struct {
	ID string `json:"id"`
}

type AccountSvcr struct {
	FinInstnID FinancialInstitutionID `json:"finInstnId"`
}

type FinancialInstitutionID struct {
	BIC string `json:"bic"`
}

type Balance struct {
	Tp        BalanceType `json:"tp"`
	Amt       DbsAmount   `json:"amt"`
	CdtDbtInd string      `json:"cdtDbtInd"`
	Dt        DateObj     `json:"dt"`
}

type BalanceType struct {
	CdOrPrtry CodeOrProprietary `json:"cdOrPrtry"`
}

type CodeOrProprietary struct {
	Cd string `json:"cd"`
}

type DbsAmount struct {
	Value float64 `json:"value"`
	Ccy   string  `json:"ccy"`
}

type DateObj struct {
	Dt string `json:"dt"`
}

type TxnSummary struct {
	TtlNtries TotalEntries `json:"ttlNtries"`
}

type TotalEntries struct {
	NbOfNtries    string  `json:"nbOfNtries"`
	Sum           float64 `json:"sum"`
	TtlNetNtryAmt float64 `json:"ttlNetNtryAmt"`
	CdtDbtInd     string  `json:"cdtDbtInd"`
}

type Entry struct {
	NtryRef      string        `json:"ntryRef"`
	Amt          DbsAmount     `json:"amt"`
	CdtDbtInd    string        `json:"cdtDbtInd"`
	Sts          string        `json:"sts"`
	BookgDt      DateTimeObj   `json:"bookgDt"`
	ValDt        DateObj       `json:"valDt"`
	AcctSvcrRef  string        `json:"acctSvcrRef"`
	BkTxCd       BankTxCode    `json:"bkTxCd"`
	NtryDtls     []EntryDetail `json:"ntryDtls"`
	AddtlNtryInf string        `json:"addtlNtryInf"`
}

type DateTimeObj struct {
	DtTm string `json:"dtTm"`
}

type BankTxCode struct {
	Prtry ProprietaryCode `json:"prtry"`
}

type ProprietaryCode struct {
	Cd string `json:"cd"`
}

type EntryDetail struct {
	TxDtls []TransactionDetail `json:"txDtls"`
}

type TransactionDetail struct {
	Refs      ReferenceDetails `json:"refs"`
	AmtDtls   AmountDetails    `json:"amtDtls"`
	RltdPties RelatedParties   `json:"rltdPties"`
}

type ReferenceDetails struct {
	EndToEndID string `json:"endToEndId"`
}

type AmountDetails struct {
	InstdAmt InstructedAmount `json:"instdAmt"`
}

type InstructedAmount struct {
	Amt     DbsAmount         `json:"amt"`
	CcyXchg *CurrencyExchange `json:"ccyXchg,omitempty"`
}

type CurrencyExchange struct {
	SrcCcy   string  `json:"srcCcy"`
	TrgtCcy  string  `json:"trgtCcy"`
	XchgRate float64 `json:"xchgRate"`
	CtrctID  string  `json:"ctrctId"`
}

type RelatedParties struct {
	Dbtr     Party     `json:"dbtr"`
	Cdtr     Party     `json:"cdtr"`
	CdtrAcct AccountID `json:"cdtrAcct"`
}

type Party struct {
	Nm string   `json:"nm"`
	ID *PartyID `json:"id,omitempty"`
}

type PartyID struct {
	OrgID OrgIDDetails `json:"orgId"`
}

type OrgIDDetails struct {
	Othr []IDValue `json:"othr"`
}
