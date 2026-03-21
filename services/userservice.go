package services

import (
	"Payback_BE/models"
	"Payback_BE/repo"
	"context"
	"database/sql"
)

type UserService interface {
	SaveUser(ctx context.Context, number string, name string, id int) (*models.User, error)
	FindUserFromNum(ctx context.Context, number string) (*models.User, error)
	FindUser(ctx context.Context, id int) (*models.User, error)
}

type UserTaskService struct {
	UserRepository repo.UserRepo
}

func NewUserTaskService(db *sql.DB) *UserTaskService {
	return &UserTaskService{
		UserRepository: repo.NewPgUserRepo(db),
	}
}

func (us *UserTaskService) FindUserFromNum(ctx context.Context, number string) (*models.User, error) {
	user, err := us.UserRepository.FindUserFromNum(ctx, number)
	if err != nil {
		return nil, err
	}
	// we found a user
	return user, nil
}

func (us *UserTaskService) FindUser(ctx context.Context, id int) (*models.User, error) {
	user, err := us.UserRepository.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserTaskService) SaveUser(ctx context.Context, number string, userName string, id int) (*models.User, error) {
	user := &models.User{
		Number: number,
		Name:   userName,
		ID:     id,
	}
	user, err := us.UserRepository.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/*
func (us *UserService) LookupNumber(ctx context.Context, phoneNumber int, pointStep int) (*models.User, error) {
	user, err := UserService.UserRepository.FindUser()
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
*/
