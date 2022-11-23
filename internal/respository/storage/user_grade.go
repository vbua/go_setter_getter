package storage

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/vbua/go_setter_getter/internal/entity"
	"sync"
)

type UserGradeRepo struct {
	userGrades map[string]entity.UserGrade
	*sync.RWMutex
}

func NewUserGradeRepo() UserGradeRepo {
	events := make(map[string]entity.UserGrade)
	return UserGradeRepo{events, &sync.RWMutex{}}
}

func (u *UserGradeRepo) Set(userGrade entity.UserGrade) {
	u.RLock()
	exUserGrade, ok := u.userGrades[userGrade.UserId]
	u.RUnlock()
	if ok {
		err := mergo.Merge(&exUserGrade, userGrade, mergo.WithOverride)
		if err != nil {
			fmt.Println(err)
			return
		}

		u.Lock()
		u.userGrades[userGrade.UserId] = exUserGrade
		u.Unlock()
		fmt.Println(u.userGrades)
		return
	}

	u.Lock()
	u.userGrades[userGrade.UserId] = userGrade
	u.Unlock()
}

func (u *UserGradeRepo) Get(userId string) (*entity.UserGrade, error) {
	u.RLock()
	userGrade, ok := u.userGrades[userId]
	u.RUnlock()
	if !ok {
		return nil, fmt.Errorf("user grade wasn't found")
	}
	return &userGrade, nil
}

func (u *UserGradeRepo) GetAll() map[string]entity.UserGrade {
	return u.userGrades
}

func (u *UserGradeRepo) SetMany(userGrades map[string]entity.UserGrade) {
	u.Lock()
	u.userGrades = userGrades
	u.Unlock()
}
