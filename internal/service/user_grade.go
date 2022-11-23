package service

import (
	"github.com/vbua/go_setter_getter/internal/entity"
)

type UserGradeRepo interface {
	Set(entity.UserGrade)
	Get(userId string) (*entity.UserGrade, error)
	GetAll() map[string]entity.UserGrade
	SetMany(map[string]entity.UserGrade)
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

func (u *UserGradeService) Set(userGrade entity.UserGrade, needPublish bool) {
	u.userGradeRepo.Set(userGrade)
	if needPublish {
		u.userGradeNatsRepo.Publish(userGrade)
	}
}

func (u *UserGradeService) Get(userId string) (*entity.UserGrade, error) {
	userGrade, err := u.userGradeRepo.Get(userId)
	if err != nil {
		return nil, err
	}
	return userGrade, nil
}

func (u *UserGradeService) Backup() map[string]entity.UserGrade {
	userGrades := u.userGradeRepo.GetAll()
	return userGrades
}

func (u *UserGradeService) SetMany(userGrades map[string]entity.UserGrade) {
	u.userGradeRepo.SetMany(userGrades)
}
