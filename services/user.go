package services

import (
	"Payback_BE/models"
	"Payback_BE/repo"
	"context"
	"database/sql"
)

type UserService struct {
	repo repo.UserRepo
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		repo: repo.NewPgUserRepo(db),
	}
}

func (us *UserService) LookupNumber(ctx context.Context, phoneNumber int, pointStep int) (*models.User, error) {
	user, err := us.repo.LookUp(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	// we found a user
	visited := user.AddVisit(pointStep)

	// extra db call here, todo come back to this for telemetry purposes
	savedUser, savedErr := us.repo.SaveUser(ctx, visited)
	if savedErr != nil {
		return nil, savedErr
	}
	return savedUser, nil
}

func (us *UserService) RegisterNumber(ctx context.Context, phoneNumber int, pointStep int) (*models.User, error) {
	user, err := us.repo.LookUp(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	visited := user.AddVisit(pointStep)

	// extra db call here, todo come back to this for telemetry purposes
	savedUser, savedErr := us.repo.SaveUser(ctx, visited)
	if savedErr != nil {
		return nil, savedErr
	}
	return savedUser, nil
}

// this should have an automatic add visit call too
func (us *UserService) RedeemPoints(ctx context.Context, phoneNumber int, redeemAmount int) (*models.User, error) {
	user, err := us.repo.LookUp(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}
	redeemed := user.RedeemPoints(redeemAmount)

	savedUser, savedErr := us.repo.SaveUser(ctx, redeemed)
	if savedErr != nil {
		return nil, savedErr
	}
	return savedUser, nil
}
