package services

import (
	"context"

	"github.com/aakritigkmit/payment-gateway/internal/dto"
	"github.com/aakritigkmit/payment-gateway/internal/helpers"
	"github.com/aakritigkmit/payment-gateway/internal/model"
	"github.com/aakritigkmit/payment-gateway/internal/repository"
)

type DBSService struct {
	DBSRepo *repository.DBSRepo
}

func NewDBSService(dbsRepo *repository.DBSRepo) *DBSService {
	return &DBSService{DBSRepo: dbsRepo}
}

func (s *DBSService) ProcessBankStatement(req dto.CAMT053Request) error {

	data := helpers.MapCAMT053DTOToModel(&req)
	// Save the raw incoming request as-is
	if err := s.DBSRepo.SaveBankStatement(context.Background(), data); err != nil {
		return err

	}
	return nil
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
