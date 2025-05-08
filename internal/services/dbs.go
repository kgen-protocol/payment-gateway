package services

import (
	"context"
	"time"

	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
)

type DBSService struct {
	DBSRepo *repository.DBSRepo
}

func NewDBSService(dbsRepo *repository.DBSRepo) *DBSService {
	return &DBSService{DBSRepo: dbsRepo}
}

func (s *DBSService) ProcessBankStatement(req model.Camt053Request) (*model.Camt053Response, error) {
	// Save the raw incoming request as-is
	if err := s.DBSRepo.SaveBankStatement(context.Background(), req); err != nil {
		return nil, err
	}

	// Simulate or generate response (as needed)
	resp := &model.Camt053Response{
		Header: model.Camt053Header{
			MsgId:     req.Header.MsgId,
			OrgId:     req.Header.OrgId,
			TimeStamp: time.Now().Format(time.RFC3339),
			Country:   req.Header.Country,
		},
		TxnEnqResponse: model.TxnEnqResponse{
			EnqStatus: "ACSP",
			AcctInfo: &model.AcctInfo{
				AccountNo:  req.TxnInfo.AccountNo,
				AccountCcy: req.TxnInfo.AccountCcy,
			},
			BizDate:     req.TxnInfo.BizDate,
			MessageType: req.TxnInfo.MessageType,
		},
	}

	return resp, nil
}

func (s *DBSService) ProcessIntradayNotification(payload model.IntradayNotificationPayload) error {
	// Save to DB
	if err := s.DBSRepo.SaveIntradayNotification(context.Background(), payload); err != nil {
		return err
	}

	return nil
}

func (s *DBSService) ProcessIncomingNotification(payload model.IncomingNotificationPayload) error {
	// Save to DB
	if err := s.DBSRepo.SaveIncomingNotification(context.Background(), payload); err != nil {
		return err
	}

	return nil
}
