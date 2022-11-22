package service

import (
	"github.com/vbua/go_setter_getter/internal/entity"
)

type UserGradeRepo interface {
	Set(entity.UserGrade)
	Get(userId string) (*entity.UserGrade, error)
}

type UserGradeNatsRepo interface {
	Publish(entity.UserGrade) error
}

type UserGradeService struct {
	userGradeRepo     UserGradeRepo
	userGradeNatsRepo UserGradeNatsRepo
}

func NewUserGradeService(userGradeRepo UserGradeRepo, userGradeNatsRepo UserGradeNatsRepo) UserGradeService {
	return UserGradeService{userGradeRepo, userGradeNatsRepo}
}

func (u *UserGradeService) Set(userGrade entity.UserGrade) {
	u.userGradeRepo.Set(userGrade)
	u.userGradeNatsRepo.Publish(userGrade)
}

func (u *UserGradeService) Get(userId string) (*entity.UserGrade, error) {
	userGrade, err := u.userGradeRepo.Get(userId)
	if err != nil {
		return nil, err
	}
	return userGrade, nil
}
