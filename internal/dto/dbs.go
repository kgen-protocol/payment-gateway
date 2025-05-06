// // dto/request.go
package dto

// type Camt053Header struct {
// 	MsgId     string `json:"msgId"`
// 	OrgId     string `json:"orgId"`
// 	TimeStamp string `json:"timeStamp"`
// 	Ctry      string `json:"ctry"`
// }

// type Camt053TxnInfo struct {
// 	AccountNo   string `json:"accountNo"`
// 	AccountCcy  string `json:"accountCcy"`
// 	BizDate     string `json:"bizDate"`
// 	MessageType string `json:"messageType"`
// }

// type Camt053Request struct {
// 	Header  Camt053Header  `json:"header"`
// 	TxnInfo Camt053TxnInfo `json:"txnInfo"`
// }

// type Camt053Response struct {
// 	Header         Camt053Header  `json:"header"`
// 	TxnEnqResponse TxnEnqResponse `json:"txnEnqResponse"`
// }

// type TxnEnqResponse struct {
// 	EnqStatus            string        `json:"enqStatus"`
// 	EnqRejectCode        string        `json:"enqRejectCode,omitempty"`
// 	EnqStatusDescription string        `json:"enqStatusDescription,omitempty"`
// 	AcctInfo             *AcctInfo     `json:"acctInfo,omitempty"`
// 	BizDate              string        `json:"bizDate,omitempty"`
// 	MessageType          string        `json:"messageType,omitempty"`
// 	Statement            []interface{} `json:"statement,omitempty"`
// }

// type AcctInfo struct {
// 	AccountNo  string `json:"accountNo"`
// 	AccountCcy string `json:"accountCcy"`
// }

type Camt053Request struct {
	Header         Camt053Header  `json:"header"`
	TxnEnqResponse TxnEnqResponse `json:"txnEnqResponse"`
}

type Camt053Header struct {
	MsgId     string `json:"msgId"`
	OrgId     string `json:"orgId"`
	TimeStamp string `json:"timeStamp"`
	Country   string `json:"ctry"`
}

type TxnEnqResponse struct {
	EnqStatus   string      `json:"enqStatus"`
	AcctInfo    AcctInfo    `json:"acctInfo"`
	BizDate     string      `json:"bizDate"`
	MessageType string      `json:"messageType"`
	Statement   []Statement `json:"statement"`
}

type AcctInfo struct {
	AccountNo  string `json:"accountNo"`
	AccountCcy string `json:"accountCcy"`
}

type Statement struct {
	BkToCstmrStmt BkToCstmrStmt `json:"bkToCstmrStmt"`
}

type BkToCstmrStmt struct {
	GrpHdr GrpHdr   `json:"grpHdr"`
	Stmt   []StmtEl `json:"stmt"`
}

type GrpHdr struct {
	MsgId   string `json:"msgId"`
	CreDtTm string `json:"creDtTm"`
}

type StmtEl struct {
	Id        string    `json:"id"`
	CreDtTm   string    `json:"creDtTm"`
	Acct      Account   `json:"acct"`
	Bal       []Bal     `json:"bal"`
	TxsSummry TxsSummry `json:"txsSummry"`
	Ntry      []Entry   `json:"ntry"`
}

type Account struct {
	Id struct {
		Othr struct {
			Id string `json:"id"`
		} `json:"othr"`
	} `json:"id"`
	Ccy  string `json:"ccy"`
	Nm   string `json:"nm"`
	Svcr struct {
		FinInstnId struct {
			Bic string `json:"bic"`
		} `json:"finInstnId"`
	} `json:"svcr"`
}

type Bal struct {
	Tp struct {
		CdOrPrtry struct {
			Cd string `json:"cd"`
		} `json:"cdOrPrtry"`
	} `json:"tp"`
	Amt struct {
		Value float64 `json:"value"`
		Ccy   string  `json:"ccy"`
	} `json:"amt"`
	CdtDbtInd string `json:"cdtDbtInd"`
	Dt        struct {
		Dt string `json:"dt"`
	} `json:"dt"`
}

type TxsSummry struct {
	TtlNtries struct {
		NbOfNtries    string  `json:"nbOfNtries"`
		Sum           float64 `json:"sum"`
		TtlNetNtryAmt float64 `json:"ttlNetNtryAmt"`
		CdtDbtInd     string  `json:"cdtDbtInd"`
	} `json:"ttlNtries"`
}

type Entry struct {
	NtryRef string `json:"ntryRef"`
	Amt     struct {
		Value float64 `json:"value"`
		Ccy   string  `json:"ccy"`
	} `json:"amt"`
	CdtDbtInd string `json:"cdtDbtInd"`
	Sts       string `json:"sts"`
	BookgDt   struct {
		DtTm string `json:"dtTm"`
	} `json:"bookgDt"`
	ValDt struct {
		Dt string `json:"dt"`
	} `json:"valDt"`
	AcctSvcrRef string `json:"acctSvcrRef"`
	BkTxCd      struct {
		Prtry struct {
			Cd string `json:"cd"`
		} `json:"prtry"`
	} `json:"bkTxCd"`
	NtryDtls []struct {
		TxDtls []struct {
			Refs struct {
				EndToEndId string `json:"endToEndId"`
			} `json:"refs"`
			AmtDtls struct {
				InstdAmt struct {
					Amt struct {
						Value float64 `json:"value"`
						Ccy   string  `json:"ccy"`
					} `json:"amt"`
					CcyXchg struct {
						SrcCcy   string  `json:"srcCcy"`
						TrgtCcy  string  `json:"trgtCcy"`
						XchgRate float64 `json:"xchgRate"`
						CtrctId  string  `json:"ctrctId"`
					} `json:"ccyXchg"`
				} `json:"instdAmt"`
			} `json:"amtDtls"`
			RltdPties struct {
				Dbtr struct {
					Nm string `json:"nm"`
					Id struct {
						OrgId struct {
							Othr []struct {
								Id string `json:"id"`
							} `json:"othr"`
						} `json:"orgId"`
					} `json:"id"`
				} `json:"dbtr"`
				Cdtr struct {
					Nm string `json:"nm"`
				} `json:"cdtr"`
				CdtrAcct struct {
					Id struct {
						Othr struct {
							Id string `json:"id"`
						} `json:"othr"`
					} `json:"id"`
				} `json:"cdtrAcct"`
			} `json:"rltdPties"`
		} `json:"txDtls"`
	} `json:"ntryDtls"`
	AddtlNtryInf string `json:"addtlNtryInf"`
}
