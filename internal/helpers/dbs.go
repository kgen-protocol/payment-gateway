package helpers

import (
	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/model"
)

func mapPartyID(id *dto.PartyID) *model.PartyID {
	if id == nil {
		return nil
	}
	othr := make([]model.IDValue, 0)
	for _, o := range id.OrgID.Othr {
		othr = append(othr, model.IDValue{
			ID: o.ID,
		})
	}
	return &model.PartyID{
		OrgID: model.OrgIDDetails{
			Othr: othr,
		},
	}
}

func MapCAMT053DTOToModel(dto *dto.CAMT053Request) model.CAMT053Request {
	statements := make([]model.StatementWrapper, 0)
	for _, sw := range dto.TxnEnqResponse.Statement {
		stmts := make([]model.Statement, 0)
		for _, stmt := range sw.BkToCstmrStmt.Stmt {
			// Map balances
			balances := make([]model.Balance, 0)
			for _, b := range stmt.Bal {
				balances = append(balances, model.Balance{
					Tp: model.BalanceType{
						CdOrPrtry: model.CodeOrProprietary{
							Cd: b.Tp.CdOrPrtry.Cd,
						},
					},
					Amt: model.DbsAmount{
						Value: b.Amt.Value,
						Ccy:   b.Amt.Ccy,
					},
					CdtDbtInd: b.CdtDbtInd,
					Dt: model.DateObj{
						Dt: b.Dt.Dt,
					},
				})
			}

			// Map entries
			entries := make([]model.Entry, 0)
			for _, e := range stmt.Ntry {
				// Entry details
				entryDetails := make([]model.EntryDetail, 0)
				for _, ed := range e.NtryDtls {
					txDetails := make([]model.TransactionDetail, 0)
					for _, td := range ed.TxDtls {
						var ccyXchg *model.CurrencyExchange
						if td.AmtDtls.InstdAmt.CcyXchg != nil {
							ccyXchg = &model.CurrencyExchange{
								SrcCcy:   td.AmtDtls.InstdAmt.CcyXchg.SrcCcy,
								TrgtCcy:  td.AmtDtls.InstdAmt.CcyXchg.TrgtCcy,
								XchgRate: td.AmtDtls.InstdAmt.CcyXchg.XchgRate,
								CtrctID:  td.AmtDtls.InstdAmt.CcyXchg.CtrctID,
							}
						}

						txDetails = append(txDetails, model.TransactionDetail{
							Refs: model.ReferenceDetails{
								EndToEndID: td.Refs.EndToEndID,
							},
							AmtDtls: model.AmountDetails{
								InstdAmt: model.InstructedAmount{
									Amt: model.DbsAmount{
										Value: td.AmtDtls.InstdAmt.Amt.Value,
										Ccy:   td.AmtDtls.InstdAmt.Amt.Ccy,
									},
									CcyXchg: ccyXchg,
								},
							},
							RltdPties: model.RelatedParties{
								Dbtr: model.Party{
									Nm: td.RltdPties.Dbtr.Nm,
									ID: mapPartyID(td.RltdPties.Dbtr.ID),
								},
								Cdtr: model.Party{
									Nm: td.RltdPties.Cdtr.Nm,
									ID: mapPartyID(td.RltdPties.Cdtr.ID),
								},
								CdtrAcct: model.AccountID{
									Othr: model.IDValue{
										ID: td.RltdPties.CdtrAcct.Othr.ID,
									},
								},
							},
						})
					}

					entryDetails = append(entryDetails, model.EntryDetail{
						TxDtls: txDetails,
					})
				}

				entries = append(entries, model.Entry{
					NtryRef:     e.NtryRef,
					Amt:         model.DbsAmount{Value: e.Amt.Value, Ccy: e.Amt.Ccy},
					CdtDbtInd:   e.CdtDbtInd,
					Sts:         e.Sts,
					BookgDt:     model.DateTimeObj{DtTm: e.BookgDt.DtTm},
					ValDt:       model.DateObj{Dt: e.ValDt.Dt},
					AcctSvcrRef: e.AcctSvcrRef,
					BkTxCd: model.BankTxCode{
						Prtry: model.ProprietaryCode{
							Cd: e.BkTxCd.Prtry.Cd,
						},
					},
					NtryDtls:     entryDetails,
					AddtlNtryInf: e.AddtlNtryInf,
				})
			}

			stmts = append(stmts, model.Statement{
				ID:      stmt.ID,
				CreDtTm: stmt.CreDtTm,
				Acct: model.Account{
					ID: model.AccountID{
						Othr: model.IDValue{
							ID: stmt.Acct.ID.Othr.ID,
						},
					},
					Ccy: stmt.Acct.Ccy,
					Nm:  stmt.Acct.Nm,
					Svcr: model.AccountSvcr{
						FinInstnID: model.FinancialInstitutionID{
							BIC: stmt.Acct.Svcr.FinInstnID.BIC,
						},
					},
				},
				Bal: balances,
				TxsSumm: model.TxnSummary{
					TtlNtries: model.TotalEntries{
						NbOfNtries:    stmt.TxsSumm.TtlNtries.NbOfNtries,
						Sum:           stmt.TxsSumm.TtlNtries.Sum,
						TtlNetNtryAmt: stmt.TxsSumm.TtlNtries.TtlNetNtryAmt,
						CdtDbtInd:     stmt.TxsSumm.TtlNtries.CdtDbtInd,
					},
				},
				Ntry: entries,
			})
		}

		statements = append(statements, model.StatementWrapper{
			BkToCstmrStmt: model.BankToCustomerStatement{
				GrpHdr: model.GroupHeader{
					MsgID:   sw.BkToCstmrStmt.GrpHdr.MsgID,
					CreDtTm: sw.BkToCstmrStmt.GrpHdr.CreDtTm,
				},
				Stmt: stmts,
			},
		})
	}

	return model.CAMT053Request{
		Header: model.Header{
			MsgID:     dto.Header.MsgID,
			OrgID:     dto.Header.OrgID,
			TimeStamp: dto.Header.TimeStamp,
			Country:   dto.Header.Country,
		},
		TxnEnqResponse: model.TxnEnqResponse{
			EnqStatus: dto.TxnEnqResponse.EnqStatus,
			AcctInfo: model.AcctInfo{
				AccountNo:  dto.TxnEnqResponse.AcctInfo.AccountNo,
				AccountCcy: dto.TxnEnqResponse.AcctInfo.AccountCcy,
			},
			BizDate:     dto.TxnEnqResponse.BizDate,
			MessageType: dto.TxnEnqResponse.MessageType,
			Statement:   statements,
		},
	}
}

// MapNotificationPayloadToModel maps the NotificationPayload to the corresponding model

func MapNotificationPayload(dto *dto.NotificationPayload) model.NotificationPayload {
	return model.NotificationPayload{
		Header: model.Header{
			MsgID:     dto.Header.MsgID,
			OrgID:     dto.Header.OrgID,
			TimeStamp: dto.Header.TimeStamp,
			Country:   dto.Header.Country,
		},
		TxnInfo: model.TxnInfo{
			TxnType:           dto.TxnInfo.TxnType,
			CustomerReference: dto.TxnInfo.CustomerReference,
			TxnRefID:          dto.TxnInfo.TxnRefID,
			TxnDate:           dto.TxnInfo.TxnDate,
			ValueDate:         dto.TxnInfo.ValueDate,
			ReceivingParty: model.ReceivingParty{
				Name:             dto.TxnInfo.ReceivingParty.Name,
				AccountNo:        dto.TxnInfo.ReceivingParty.AccountNo,
				VirtualAccountNo: dto.TxnInfo.ReceivingParty.VirtualAccountNo,
			},
			AmountDetails: model.NotificationAmountDetails{
				TxnCurrency: dto.TxnInfo.AmountDetails.TxnCurrency,
				TxnAmount:   dto.TxnInfo.AmountDetails.TxnAmount,
			},
			SenderParty: model.SenderParty{
				Name:         dto.TxnInfo.SenderParty.Name,
				AccountNo:    dto.TxnInfo.SenderParty.AccountNo,
				SenderBankID: dto.TxnInfo.SenderParty.SenderBankID,
			},
		
			PaymentDetails: dto.TxnInfo.PaymentDetails,
		},
	}
}
