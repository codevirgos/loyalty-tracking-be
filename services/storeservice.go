package services

import (
	"Payback_BE/models"
	"Payback_BE/repo"
	"context"
	"database/sql"
	"log"
)

type StoreService interface {
	GetOverview(ctx context.Context, userId int) (*models.Overview, error)
	StoreScan(ctx context.Context, scanned models.Scanned) (*models.Overview, error)
	RedeemPoints(ctx context.Context, redeem models.Redeem) (*models.Overview, error)
}

type StoreTaskService struct {
	StoreRepository repo.StoreRepo
}

func NewStoreTaskService(db *sql.DB) *StoreTaskService {
	return &StoreTaskService{
		StoreRepository: repo.NewPgStoreRepo(db),
	}
}

func (st *StoreTaskService) GetOverview(ctx context.Context, id int) (*models.Overview, error) {
	overview, err := st.StoreRepository.GetOverview(ctx, id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return overview, nil
}

func (st *StoreTaskService) StoreScan(ctx context.Context, scanned models.Scanned) (*models.Overview, error) {
	overview, err := st.StoreRepository.StoreScan(ctx, scanned)
	// add the activity table call too
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return overview, nil
}

func (st *StoreTaskService) RedeemPoints(ctx context.Context, redeem models.Redeem) (*models.Overview, error) {
	overview, err := st.StoreRepository.RedeemPoints(ctx, redeem)
	// add the activity table call
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return overview, nil
}
