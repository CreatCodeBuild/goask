package fakeadapter

import (
	"goask/core/adapter"
	"goask/core/entity"

	"github.com/pkg/errors"
)

// UserDAO
type UserDAO struct {
	data *Data
}

func NewUserDAO(data *Data) *UserDAO {
	return &UserDAO{data: data}
}

func (d *UserDAO) UserByID(ID entity.ID) (entity.User, error) {
	for _, user := range d.data.users {
		if user.ID == ID {
			return user, nil
		}
	}
	return entity.User{}, errors.WithStack(&adapter.ErrUserNotFound{ID: ID})
}

func (d *UserDAO) Users() ([]entity.User, error) {
	return d.data.users, nil
}

func (d *UserDAO) CreateUser(name string) (entity.User, error) {
	user := entity.User{ID: entity.NewIDInt(len(d.data.users) + 1), Name: name}
	d.data.users = append(d.data.users, user)
	return user, d.data.serialize()
}

func (d *UserDAO) QuestionsOfUser(ID entity.ID) ([]entity.Question, error) {
	var ret []entity.Question
	for _, q := range d.data.questions {
		if q.AuthorID == ID {
			ret = append(ret, q)
		}
	}
	return ret, nil
}
